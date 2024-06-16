package socket

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"net"
	"time"
	"universal/framework_new/common/base"
)

const (
	MAX_PACKET_BODY_SIZE = 1024 * 1024 // 包头大小
)

type IFrame interface {
	GetHeadSize() int            // 获取包头大小
	GetBodySize([]byte) int      // 获取包体
	Check([]byte, []byte) bool   // 数据包校验
	Build([]byte, []byte) []byte // 转成完整帧
}

type Socket struct {
	IFrame
	conn       net.Conn      // 通信
	readExpire time.Duration // 读超时
	readBytes  []byte        // 接受缓存
	sendExpire time.Duration // 写超时
	sendBytes  []byte        // 接受缓存
}

func NewSocket(fr IFrame, conn net.Conn) *Socket {
	return &Socket{
		IFrame:    fr,
		conn:      conn,
		readBytes: make([]byte, 512*1024),
		sendBytes: make([]byte, 512*1024),
	}
}

func (d *Socket) getBytes(size int, isread bool) (ret []byte) {
	if isread {
		if cap(d.readBytes) < size {
			d.readBytes = make([]byte, size)
		}
		ret = d.readBytes[:size]
	} else {
		if cap(d.sendBytes) < size {
			d.sendBytes = make([]byte, size)
		}
		ret = d.sendBytes[:size]
	}
	return
}

func (d *Socket) Send(buf []byte) error {
	// 设置发送超时时间，避免阻塞
	if d.sendExpire > 0 {
		d.conn.SetWriteDeadline(time.Now().Add(d.sendExpire))
	}
	// 组包
	size := len(buf) + d.GetHeadSize()
	pack := d.Build(d.getBytes(size, false), buf)
	// 发送数据包
	_, err := d.conn.Write(pack)
	return err
}

func (d *Socket) Read() ([]byte, error) {
	// 设置接受超时时间，避免阻塞
	if d.readExpire > 0 {
		d.conn.SetReadDeadline(time.Now().Add(d.readExpire))
	}
	// 读取包头
	headSize := d.GetHeadSize()
	buf := d.getBytes(headSize, true)
	n, err := io.ReadFull(d.conn, buf)
	if err != nil {
		return nil, err
	}
	if n != headSize {
		return nil, fmt.Errorf("packet header is incomplete, size: %d, receive: %d", headSize, n)
	}
	// 获取包体长度
	bodySize := d.GetBodySize(buf)
	if bodySize <= 0 {
		return nil, fmt.Errorf("packet body is empty")
	}
	if bodySize > MAX_PACKET_BODY_SIZE {
		return nil, fmt.Errorf("packet body exceed maximum limit(%d)", MAX_PACKET_BODY_SIZE)
	}
	// 接受包体
	headStr := string(buf)
	buf = d.getBytes(bodySize, true)
	n, err = io.ReadFull(d.conn, buf)
	if err != nil {
		return nil, err
	}
	if n != int(bodySize) {
		return nil, fmt.Errorf("packet body is incomplete, size: %d, receive: %d", bodySize, n)
	}
	// 校验数据包
	if !d.Check(base.StringToBytes(headStr), buf) {
		return nil, fmt.Errorf("packet check failed")
	}
	return buf, nil
}

// websocket包结构
// bodySize(4B) | md5(4B) | body
type Frame struct{}

func (d *Frame) GetHeadSize() int {
	return 10
}

func (d *Frame) GetBodySize(head []byte) int {
	return int(binary.LittleEndian.Uint32(head))
}

func (d *Frame) Check(head []byte, body []byte) bool {
	oldCrc := binary.LittleEndian.Uint32(head[4:])
	crc := crc32.ChecksumIEEE(body)
	return crc == oldCrc
}

func (d *Frame) Build(frame []byte, body []byte) []byte {
	// 设置包头
	binary.LittleEndian.PutUint32(frame, uint32(len(body)))
	binary.LittleEndian.PutUint32(frame[4:], crc32.ChecksumIEEE(body))
	// 拷贝
	headSize := d.GetHeadSize()
	copy(frame[headSize:], body)
	return frame
}
