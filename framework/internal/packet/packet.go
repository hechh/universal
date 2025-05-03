package packet

import (
	"encoding/binary"
	"universal/framework/define"
)

type Packet struct {
	header *Header
	body   []byte
}

func NewPacket(header define.IHeader, body []byte) define.IPacket {
	return &Packet{
		header: header.(*Header),
		body:   body,
	}
}

func (p *Packet) GetHeader() define.IHeader {
	return p.header
}

func (p *Packet) GetBody() []byte {
	return p.body
}

func (p *Packet) ToBytes() (rets []byte) {
	rets = make([]byte, 44+len(p.body)+len(p.header.ActorName)+len(p.header.FuncName))
	pos := 0
	binary.BigEndian.PutUint32(rets[pos:], uint32(p.header.SrcNodeType))
	pos += 4
	binary.BigEndian.PutUint32(rets[pos:], uint32(p.header.SrcNodeId))
	pos += 4
	binary.BigEndian.PutUint32(rets[pos:], uint32(p.header.DstNodeType))
	pos += 4
	binary.BigEndian.PutUint32(rets[pos:], uint32(p.header.DstNodeId))
	pos += 4
	binary.BigEndian.PutUint64(rets[pos:], p.header.Uid)
	pos += 8
	binary.BigEndian.PutUint64(rets[pos:], p.header.RouteId)
	pos += 8
	binary.BigEndian.PutUint32(rets[pos:], p.header.Cmd)
	pos += 4
	lactor := len(p.header.ActorName)
	binary.BigEndian.PutUint32(rets[pos:], uint32(lactor))
	pos += 4
	copy(rets[pos:], []byte(p.header.ActorName))
	pos += lactor
	lfunc := len(p.header.FuncName)
	binary.BigEndian.PutUint32(rets[pos:], uint32(lfunc))
	pos += 4
	copy(rets[pos:], []byte(p.header.FuncName))
	pos += lfunc
	copy(rets[pos:], p.body)
	return
}

func ParsePacket(data []byte) (define.IPacket, error) {
	pack := &Packet{header: &Header{}}
	pos := 0
	pack.header.SrcNodeType = binary.BigEndian.Uint32(data[pos:])
	pos += 4
	pack.header.SrcNodeId = binary.BigEndian.Uint32(data[pos:])
	pos += 4
	pack.header.DstNodeType = binary.BigEndian.Uint32(data[pos:])
	pos += 4
	pack.header.DstNodeId = binary.BigEndian.Uint32(data[pos:])
	pos += 4
	pack.header.Uid = binary.BigEndian.Uint64(data[pos:])
	pos += 8
	pack.header.RouteId = binary.BigEndian.Uint64(data[pos:])
	pos += 8
	pack.header.Cmd = binary.BigEndian.Uint32(data[pos:])
	pos += 4
	lmethod := int(binary.BigEndian.Uint32(data[pos:]))
	pos += 4
	pack.header.Methond = string(data[pos : pos+lmethod])
	pos += lmethod
	pack.body = make([]byte, len(data)-pos)
	copy(pack.body, data[pos:])
	return pack, nil
}
