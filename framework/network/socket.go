package network

import (
	"encoding/binary"
	"time"
	"universal/common/pb"
	"universal/framework/define"
	"universal/library/uerror"

	"github.com/gorilla/websocket"
)

type Socket struct {
	max    int
	frame  define.IFrame
	conn   *websocket.Conn
	wbytes []byte
}

func NewSocket(conn *websocket.Conn, max int) *Socket {
	return &Socket{
		max:    max,
		conn:   conn,
		wbytes: make([]byte, max/2),
	}
}

func (d *Socket) Register(frame define.IFrame) {
	d.frame = frame
}

func (d *Socket) Close() error {
	return d.conn.Close()
}

func (d *Socket) newWrite(size int) (ret []byte) {
	if cap(d.wbytes) < size {
		d.wbytes = make([]byte, size)
	}
	ret = d.wbytes[:size]
	return
}

func (d *Socket) SetReadExpire(expire int64) {
	if expire > 0 {
		d.conn.SetReadDeadline(time.Now().Add(time.Duration(expire) * time.Second))
	} else {
		d.conn.SetReadDeadline(time.Time{})
	}
}

func (d *Socket) SetWriteExpire(expire int64) {
	if expire > 0 {
		d.conn.SetWriteDeadline(time.Now().Add(time.Duration(expire) * time.Second))
	} else {
		d.conn.SetWriteDeadline(time.Time{})
	}
}

func (d *Socket) Write(pack *pb.Packet) error {
	ll := d.frame.GetSize(pack)
	if ll > d.max {
		return uerror.New(1, int32(pb.ErrorCode_SOCKET_MAX_LIMIT), "超过最大包长限制: %d", d.max)
	}

	buf := d.newWrite(ll + 4)
	binary.BigEndian.PutUint32(buf, uint32(ll))
	if err := d.frame.Encode(pack, buf[4:]); err != nil {
		return err
	}
	return d.conn.WriteMessage(websocket.BinaryMessage, buf)
}

func (d *Socket) Read(recv *pb.Packet) error {
	_, buf, err := d.conn.ReadMessage()
	if err != nil {
		return err
	}
	return d.frame.Decode(buf[4:], recv)
}
