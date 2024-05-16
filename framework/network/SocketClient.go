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

const (
	SocketFrameMaxSize = 1024 * 1024 // 包头大小
)

type IFrame interface {
	Check([]byte, []byte) error  // 数据包校验
	Build([]byte, []byte) []byte // 转成完整帧
	GetHeadSize() int            // 获取包头大小
	GetBodySize([]byte) int      // 获取包体
}

type SocketClient struct {
	IFrame
	conn       net.Conn      // 通信
	readExpire time.Duration // 读超时
	sendExpire time.Duration // 写超时
	recvBuff   []byte        // 接受缓存
	sendBuff   []byte        // 接受缓存
}

func NewSocketClient(conn net.Conn) *SocketClient {
	return &SocketClient{
		IFrame:   &SocketFrame{},
		conn:     conn,
		recvBuff: make([]byte, 512*1024),
		sendBuff: make([]byte, 512*1024),
	}
}

func (d *SocketClient) SendRsp(head *pb.PacketHead, item proto.Message, params ...interface{}) error {
	pac, err := fbasic.RspToPacket(head, item, params...)
	if err != nil {
		return err
	}
	return d.Send(pac)
}

func (d *SocketClient) Send(pac *pb.Packet) error {
	buf, err := proto.Marshal(pac)
	if err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_ProtoMarshal, err)
	}
	return d.SendBytes(buf)
}

func (d *SocketClient) Read() (*pb.Packet, error) {
	buf, err := d.ReadBytes()
	if err != nil {
		return nil, err
	}
	ret := &pb.Packet{}
	if err := proto.Unmarshal(buf, ret); err != nil {
		return nil, fbasic.NewUError(1, pb.ErrorCode_ProtoUnmarshal, err)
	}
	return ret, nil
}

func (d *SocketClient) getSendBytes(size int) []byte {
	if cap(d.sendBuff) >= size {
		return d.sendBuff[:size]
	}
	d.sendBuff = make([]byte, size)
	return d.sendBuff
}

func (d *SocketClient) SendBytes(buf []byte) error {
	// 设置发送超时时间，避免阻塞
	if d.sendExpire > 0 {
		d.conn.SetWriteDeadline(time.Now().Add(d.sendExpire))
	}
	// 限制检测
	size := len(buf) + d.GetHeadSize()
	if size > SocketFrameMaxSize {
		return fbasic.NewUError(1, pb.ErrorCode_SocketFrameSizeMaxLimit, size, SocketFrameMaxSize)
	}
	// 发送数据包
	if _, err := d.conn.Write(d.Build(d.getSendBytes(size), buf)); err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_SocketClientSend, err)
	}
	return nil
}

func (d *SocketClient) getReadBytes(size int) []byte {
	if cap(d.recvBuff) >= size {
		return d.recvBuff[:size]
	}
	d.recvBuff = make([]byte, size)
	return d.recvBuff
}

func (d *SocketClient) ReadBytes() ([]byte, error) {
	// 设置接受超时时间，避免阻塞
	if d.readExpire > 0 {
		d.conn.SetReadDeadline(time.Now().Add(d.readExpire))
	}
	// 读取包头
	headSize := d.GetHeadSize()
	buf := d.getReadBytes(headSize)
	n, err := io.ReadFull(d.conn, buf)
	if err != nil {
		return nil, fbasic.NewUError(1, pb.ErrorCode_SocketClientRead, n, err)
	}
	if n != headSize {
		return nil, fbasic.NewUError(1, pb.ErrorCode_SocketFrameHeaderSize, n, headSize)
	}
	// 获取包体长度
	bodySize := d.GetBodySize(buf)
	if bodySize <= 0 || bodySize > SocketFrameMaxSize {
		return nil, fbasic.NewUError(1, pb.ErrorCode_SocketFrameBodySizeMaxLimit, bodySize, SocketFrameMaxSize)
	}
	// 接受包体
	headStr := string(buf)
	buf = d.getReadBytes(bodySize)
	n, err = io.ReadFull(d.conn, buf)
	if err != nil {
		return nil, fbasic.NewUError(1, pb.ErrorCode_SocketClientRead, err)
	}
	if n != int(bodySize) {
		return nil, fbasic.NewUError(1, pb.ErrorCode_SocketClientRead, fmt.Sprint(bodySize, n))
	}
	// 校验数据包
	if err := d.Check([]byte(headStr), buf); err != nil {
		return nil, fbasic.NewUError(1, pb.ErrorCode_SocketFrameCheck, err)
	}
	return buf, nil
}
