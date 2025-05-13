package socket

import (
	"fmt"
	"io"
	"net"
	"time"
	"universal/framework/library/util"
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

func (d *Socket) Close() error {
	return d.conn.Close()
}

func (d *Socket) getSendBytes(size int) (ret []byte) {
	if cap(d.sendBytes) < size {
		d.sendBytes = make([]byte, size)
	}
	ret = d.sendBytes[:size]
	return
}

func (d *Socket) Send(buf []byte) (int, error) {
	// 设置发送超时时间，避免阻塞
	if d.sendExpire > 0 {
		d.conn.SetWriteDeadline(time.Now().Add(d.sendExpire))
	}
	// 组包
	size := len(buf) + d.GetHeadSize()
	pack := d.Build(d.getSendBytes(size), buf)
	// 发送数据包
	return d.conn.Write(pack)
}

func (d *Socket) getReadBytes(size int) (ret []byte) {
	if cap(d.readBytes) < size {
		d.readBytes = make([]byte, size)
	}
	ret = d.readBytes[:size]
	return
}

func (d *Socket) Read() (recv []byte, err error) {
	// 设置接受超时时间，避免阻塞
	if d.readExpire > 0 {
		d.conn.SetReadDeadline(time.Now().Add(d.readExpire))
	}
	// 读取包头
	headSize := d.GetHeadSize()
	head := d.getReadBytes(headSize)
	if n, err := io.ReadFull(d.conn, head); err != nil {
		return nil, err
	} else if n != headSize {
		return nil, fmt.Errorf("packet header is incomplete, size: %d, receive: %d", headSize, n)
	}
	headStr := string(head)
	// 读取包体
	bodySize := d.GetBodySize(head)
	body := d.getReadBytes(bodySize)
	if n, err := io.ReadFull(d.conn, body); err != nil {
		return nil, err
	} else if n != int(bodySize) {
		return nil, fmt.Errorf("packet body is incomplete, size: %d, receive: %d", bodySize, n)
	}
	// 校验数据包
	if !d.Check(util.StringToBytes(headStr), body) {
		return nil, fmt.Errorf("packet check failed")
	}
	recv = append(recv, body...)
	return
}
