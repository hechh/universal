package packet

import (
	"universal/framework/define"
)

type Packet struct {
	header define.IHeader
	body   []byte
}

func NewPacket(header define.IHeader, body []byte) define.IPacket {
	return &Packet{
		header: header.(*Header),
		body:   body,
	}
}

func ParsePacket(data []byte) define.IPacket {
	pack := &Packet{
		header: &Header{
			Table: &define.RouteInfo{},
		},
	}
	pack.header.Parse(data)
	llen := pack.header.GetSize()
	copy(pack.body, data[llen:])
	return pack
}

func (p *Packet) GetHeader() define.IHeader {
	return p.header
}

func (p *Packet) GetBody() []byte {
	return p.body
}

func (p *Packet) ToBytes() (rets []byte) {
	llen := p.header.GetSize()
	rets = make([]byte, len(p.body)+llen)
	rets = p.header.ToBytes(rets)
	copy(rets[llen:], p.body)
	return
}
