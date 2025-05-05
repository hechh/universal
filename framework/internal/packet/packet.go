package packet

import (
	"universal/framework/define"
)

type Packet struct {
	header define.IHeader
	body   []byte
}

func NewPacket() define.IPacket {
	return &Packet{}
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
	p.header.ToBytes(rets)
	copy(rets[llen:], p.body)
	return
}

func (p *Packet) Parse(data []byte) define.IPacket {
	p.header.Parse(data)
	llen := p.header.GetSize()
	copy(p.body, data[llen:])
	return p
}

func (p *Packet) SetHeader(h define.IHeader) define.IPacket {
	p.header = h
	return p
}

func (p *Packet) SetBody(buf []byte) define.IPacket {
	p.body = buf
	return p
}
