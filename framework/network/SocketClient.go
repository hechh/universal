package network

import (
	"fmt"
	"io"
	"net"
	"time"
	"universal/common/pb"
	"universal/framework/fbasic"

	"google.golang.org/protobuf/proto"
)

type SocketClient struct {
	conn       net.Conn
	readExpire time.Duration // 读超时
	sendExpire time.Duration // 写超时
	recvBuff   []byte        // 接受缓存
	sendBuff   []byte        // 接受缓存
}

func (d *SocketClient) getSendBytes(size int) []byte {
	if cap(d.sendBuff) >= size {
		return d.sendBuff[:size]
	}
	d.sendBuff = make([]byte, size)
	return d.sendBuff
}

func (d *SocketClient) Send(pac *pb.Packet) error {
	// 设置发送超时时间，避免阻塞
	if d.sendExpire > 0 {
		d.conn.SetWriteDeadline(time.Now().Add(d.sendExpire))
	}
	// 序列化
	buf, err := proto.Marshal(pac)
	if err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_ProtoMarshal, err)
	}
	// 限制检测
	size := len(buf) + WebSocketHeaderSize
	if size > WebSocketMaxLimit {
		return fbasic.NewUError(1, pb.ErrorCode_SocketClientMaxLimit, fmt.Sprintf("send: %d, limit: %d", size, WebSocketMaxLimit))
	}
	// 获取发送缓存
	sendBuff := d.getSendBytes(size)
	head := WebSocketHeader(sendBuff[:WebSocketHeaderSize])
	head.SetSize(uint32(len(buf)))
	// 设置crc
	head.SetCrc32(fbasic.GetCrc32(buf))
	// 设置body
	copy(sendBuff[:WebSocketHeaderSize], buf)
	// 发送数据包
	if _, err := d.conn.Write(sendBuff); err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_SocketClientSend, err)
	}
	return nil
}

func (d *SocketClient) getRecvBytes(size int) []byte {
	if cap(d.recvBuff) >= size {
		return d.recvBuff[:size]
	}
	d.recvBuff = make([]byte, size)
	return d.recvBuff
}

func (d *SocketClient) Read() (*pb.Packet, error) {
	// 设置接受超时时间，避免阻塞
	if d.readExpire > 0 {
		d.conn.SetReadDeadline(time.Now().Add(d.readExpire))
	}
	// 读取包头
	buf := d.getRecvBytes(WebSocketHeaderSize)
	n, err := io.ReadFull(d.conn, buf)
	if err != nil {
		return nil, fbasic.NewUError(1, pb.ErrorCode_SocketClientRead, err)
	}
	if n != WebSocketHeaderSize {
		return nil, fbasic.NewUError(1, pb.ErrorCode_SocketClientHeaderSize, fmt.Sprint(n, WebSocketHeaderSize))
	}
	// 获取包体长度
	crc := WebSocketHeader(buf).GetCrc32()
	bodySize := WebSocketHeader(buf).GetSize()
	if bodySize <= 0 || bodySize > WebSocketMaxBodySize {
		return nil, fbasic.NewUError(1, pb.ErrorCode_SocketClientBodySizeLimit, fmt.Sprint(bodySize, WebSocketMaxBodySize))
	}
	// 接受包体
	buf = d.getRecvBytes(int(bodySize))
	n, err = io.ReadFull(d.conn, buf)
	if err != nil {
		return nil, fbasic.NewUError(1, pb.ErrorCode_SocketClientRead, err)
	}
	if n != int(bodySize) {
		return nil, fbasic.NewUError(1, pb.ErrorCode_SocketClientRead, fmt.Sprint(bodySize, n))
	}
	// 校验包头是否被篡改
	if vcrc := fbasic.GetCrc32(buf); vcrc != crc {
		return nil, fbasic.NewUError(1, pb.ErrorCode_SocketClientCheck, fmt.Sprint(crc, vcrc))
	}
	// 解析包
	ret := &pb.Packet{}
	if err := proto.Unmarshal(buf, ret); err != nil {
		return nil, fbasic.NewUError(1, pb.ErrorCode_ProtoUnmarshal, err)
	}
	return ret, nil
}
