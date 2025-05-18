package network

import (
	"encoding/binary"
	"io"
	"net"
	"poker_server/common/pb"
	"poker_server/framework/domain"
	"poker_server/framework/library/uerror"
	"time"

	"github.com/golang/protobuf/proto"
)

// len(四个字节表示包长) | cmd（4byte） | uid(8byte) | seq(4byte) |version(4byte) | extra(4byte) | data(数据包体)

type Socket struct {
	frame   domain.IFrame // 帧协议
	maxSize int           // 最大包长限制
	conn    net.Conn      // 通信
	rbytes  []byte        // 接受缓存
	sbytes  []byte        // 接受缓存
}

func NewSocket(conn net.Conn, maxSize int) *Socket {
	return &Socket{
		maxSize: maxSize,
		conn:    conn,
		rbytes:  make([]byte, maxSize/2),
		sbytes:  make([]byte, maxSize/2),
	}
}

func (d *Socket) Register(frame domain.IFrame) {
	d.frame = frame
}

func (d *Socket) Close() error {
	return d.conn.Close()
}

func (d *Socket) newRead(size int) (ret []byte) {
	if cap(d.sbytes) < size {
		d.sbytes = make([]byte, size)
	}
	ret = d.sbytes[:size]
	return
}

func (d *Socket) newSend(size int) (ret []byte) {
	if cap(d.rbytes) < size {
		d.rbytes = make([]byte, size)
	}
	ret = d.rbytes[:size]
	return
}

// 设置接受超时时间，避免阻塞
func (d *Socket) SetReadExpire(expire int64) {
	if expire > 0 {
		d.conn.SetReadDeadline(time.Now().Add(time.Duration(expire) * time.Second))
	} else {
		d.conn.SetReadDeadline(time.Time{})
	}
}

// 设置发送超时时间，避免阻塞
func (d *Socket) SetWriteExpire(expire int64) {
	if expire > 0 {
		d.conn.SetWriteDeadline(time.Now().Add(time.Duration(expire) * time.Second))
	} else {
		d.conn.SetWriteDeadline(time.Time{})
	}
}

func (d *Socket) WriteMsg(head *pb.Head, msg proto.Message) error {
	buf, err := proto.Marshal(msg)
	if err != nil {
		return uerror.New(1, -1, "序列化数据包失败: %v", err)
	}
	return d.Write(&pb.Packet{Head: head, Body: buf})
}

func (d *Socket) Write(pack *pb.Packet) error {
	// 获取数据帧长度
	ll := d.frame.GetSize(pack)
	if ll > d.maxSize {
		return uerror.New(1, -1, "超过最大包长限制: %d", d.maxSize)
	}
	// 组包
	buf := d.newSend(ll + 4)
	binary.BigEndian.PutUint32(buf, uint32(ll))
	if err := d.frame.Encode(pack, buf[4:]); err != nil {
		return err
	}
	// 发送数据包
	if n, err := d.conn.Write(buf); err != nil {
		return uerror.New(1, -1, "发送数据包失败: %v", err)
	} else if n != ll+4 {
		return uerror.New(1, -1, "发送数据包不完整, expect: %d, receive: %d", ll+4, n)
	}
	return nil
}

func (d *Socket) Read(recv *pb.Packet) error {
	// 读取包头
	head := d.newRead(4)
	n, err := io.ReadFull(d.conn, head)
	if err != nil {
		return uerror.New(1, -1, "读取包长失败: %v", err)
	}
	if n != 4 {
		return uerror.New(1, -1, "包头不完整, expect: 4, receive: %d", n)
	}
	// 读取包体
	bodySize := int(binary.BigEndian.Uint32(head))
	body := d.newRead(bodySize)
	nbody, err := io.ReadFull(d.conn, body)
	if err != nil {
		return uerror.New(1, -1, "读取数据帧失败: %v", err)
	}
	if nbody != bodySize {
		return uerror.New(1, -1, "包体不完整, expect: %d, receive: %d", bodySize, nbody)
	}
	if nbody > d.maxSize {
		return uerror.New(1, -1, "超过最大包长限制: %d", d.maxSize)
	}
	return d.frame.Decode(body, recv)
}
