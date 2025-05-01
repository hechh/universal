package packet

import (
	"encoding/binary"
	"universal/framework/define"
)

type Packet struct {
	header *Header
	body   []byte
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
	pack.header.Cmd = binary.BigEndian.Uint32(data[pos:])
	pos += 4
	pack.header.Uid = binary.BigEndian.Uint64(data[pos:])
	pos += 8
	pack.header.RouteId = binary.BigEndian.Uint64(data[pos:])
	pos += 8
	pack.body = make([]byte, len(data)-pos)
	copy(pack.body, data[pos:])
	return pack, nil
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
	rets = make([]byte, 36+len(p.body))
	pos := 0
	binary.BigEndian.PutUint32(rets[pos:], uint32(p.header.SrcNodeType))
	pos += 4
	binary.BigEndian.PutUint32(rets[pos:], uint32(p.header.SrcNodeId))
	pos += 4
	binary.BigEndian.PutUint32(rets[pos:], uint32(p.header.DstNodeType))
	pos += 4
	binary.BigEndian.PutUint32(rets[pos:], uint32(p.header.DstNodeId))
	pos += 4
	binary.BigEndian.PutUint32(rets[pos:], p.header.Cmd)
	pos += 4
	binary.BigEndian.PutUint64(rets[pos:], p.header.Uid)
	pos += 8
	binary.BigEndian.PutUint64(rets[pos:], p.header.RouteId)
	pos += 8
	copy(rets[pos:], p.body)
	return
}
