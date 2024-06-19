package socket

import (
	"fmt"
	"io"
	"net"
	"time"
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
		readBytes: make([]byte, 512),
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

func (d *Socket) Send(buf []byte) (int, error) {
	// 设置发送超时时间，避免阻塞
	if d.sendExpire > 0 {
		d.conn.SetWriteDeadline(time.Now().Add(d.sendExpire))
	}
	// 组包
	size := len(buf) + d.GetHeadSize()
	pack := d.Build(d.getBytes(size, false), buf)
	// 发送数据包
	return d.conn.Write(pack)
}

func (d *Socket) Read(recv []byte) (int, error) {
	// 设置接受超时时间，避免阻塞
	if d.readExpire > 0 {
		d.conn.SetReadDeadline(time.Now().Add(d.readExpire))
	}
	// 读取包头
	headSize := d.GetHeadSize()
	head := d.getBytes(headSize, true)
	n, err := io.ReadFull(d.conn, head)
	if err != nil {
		return 0, err
	}
	if n != headSize {
		return 0, fmt.Errorf("packet header is incomplete, size: %d, receive: %d", headSize, n)
	}
	// 获取包体长度
	bodySize := d.GetBodySize(head)
	if bodySize <= 0 {
		return 0, fmt.Errorf("packet body is empty")
	}
	if bodySize > len(recv) {
		return 0, fmt.Errorf("packet body exceed receive buffer(%d)", len(recv))
	}
	// 接受包体
	n, err = io.ReadFull(d.conn, recv)
	if err != nil {
		return 0, err
	}
	if n != int(bodySize) {
		return 0, fmt.Errorf("packet body is incomplete, size: %d, receive: %d", bodySize, n)
	}
	// 校验数据包
	if !d.Check(head, recv[:bodySize]) {
		return 0, fmt.Errorf("packet check failed")
	}
	return bodySize, nil
}
