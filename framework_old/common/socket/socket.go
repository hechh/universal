package socket

import (
	"io"
	"net"
	"time"
	"universal/framework/common/uerror"
)

const (
	SocketFrameMaxSize = 1024 * 1024 // 包头大小
)

type IFrame interface {
	Check([]byte, []byte) bool   // 数据包校验
	Build([]byte, []byte) []byte // 转成完整帧
	GetHeadSize() int            // 获取包头大小
	GetBodySize([]byte) int      // 获取包体
}

type Socket struct {
	IFrame
	conn       net.Conn      // 通信
	readExpire time.Duration // 读超时
	sendExpire time.Duration // 写超时
	recvBuff   []byte        // 接受缓存
	sendBuff   []byte        // 接受缓存
}

func NewSocket(fr IFrame, conn net.Conn) *Socket {
	if fr == nil {
		fr = &Frame{}
	}
	return &Socket{
		IFrame:   fr,
		conn:     conn,
		recvBuff: make([]byte, 512*1024),
		sendBuff: make([]byte, 512*1024),
	}
}

func (d *Socket) getSendBytes(size int) []byte {
	if cap(d.sendBuff) >= size {
		return d.sendBuff[:size]
	}
	d.sendBuff = make([]byte, size)
	return d.sendBuff
}

func (d *Socket) SendBytes(buf []byte) error {
	// 设置发送超时时间，避免阻塞
	if d.sendExpire > 0 {
		d.conn.SetWriteDeadline(time.Now().Add(d.sendExpire))
	}
	// 限制检测
	size := len(buf) + d.GetHeadSize()
	if size > SocketFrameMaxSize {
		return uerror.NewUErrorf(1, -1, "max network size(%d) < %d", SocketFrameMaxSize, size)
	}
	// 发送数据包
	if _, err := d.conn.Write(d.Build(d.getSendBytes(size), buf)); err != nil {
		return uerror.NewUErrorf(1, -1, "%v", err)
	}
	return nil
}

func (d *Socket) getReadBytes(size int) []byte {
	if cap(d.recvBuff) >= size {
		return d.recvBuff[:size]
	}
	d.recvBuff = make([]byte, size)
	return d.recvBuff
}

func (d *Socket) ReadBytes() ([]byte, error) {
	// 设置接受超时时间，避免阻塞
	if d.readExpire > 0 {
		d.conn.SetReadDeadline(time.Now().Add(d.readExpire))
	}
	// 读取包头
	headSize := d.GetHeadSize()
	buf := d.getReadBytes(headSize)
	n, err := io.ReadFull(d.conn, buf)
	if err != nil {
		return nil, uerror.NewUErrorf(1, -1, "%v", err)
	}
	if n != headSize {
		return nil, uerror.NewUErrorf(1, -1, "receive size(%d) != headSize(%d)", n, headSize)
	}
	// 获取包体长度
	bodySize := d.GetBodySize(buf)
	if bodySize <= 0 || bodySize > SocketFrameMaxSize {
		return nil, uerror.NewUErrorf(1, -1, "body size(%d) > frameMaxSize(%d)", bodySize, SocketFrameMaxSize)
	}
	// 接受包体
	headStr := string(buf)
	buf = d.getReadBytes(bodySize)
	n, err = io.ReadFull(d.conn, buf)
	if err != nil {
		return nil, uerror.NewUErrorf(1, -1, "%v", err)
	}
	if n != int(bodySize) {
		return nil, uerror.NewUErrorf(1, -1, "receive size(%d) != body size(%d)", n, bodySize)
	}
	// 校验数据包
	if !d.Check([]byte(headStr), buf) {
		return nil, uerror.NewUError(1, -1, "crc32 check failed")
	}
	return buf, nil
}
