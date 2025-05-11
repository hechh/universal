package packet

import "universal/framework/domain"

type Packet struct {
	head  domain.IHead  // 包头
	route domain.IRoute // 路由
	body  []byte        // 包体
}

func (d *Packet) GetHead() domain.IHead {
	return d.head
}

func (d *Packet) GetRoute() domain.IRoute {
	return d.route
}

func (d *Packet) GetBody() []byte {
	return d.body
}

func (d *Packet) SetHead(head domain.IHead) domain.IPacket {
	d.head = head
	return d
}

func (d *Packet) SetRoute(route domain.IRoute) domain.IPacket {
	d.route = route
	return d
}

func (d *Packet) SetBody(body []byte) domain.IPacket {
	d.body = body
	return d
}

func (d *Packet) GetSize() int {
	return d.head.GetSize() + d.route.GetSize() + len(d.body)
}

func (d *Packet) WriteTo(buf []byte) error {
	pos := 0
	if err := d.head.WriteTo(buf); err != nil {
		return err
	}
	pos += d.head.GetSize()
	if err := d.route.WriteTo(buf[pos:]); err != nil {
		return err
	}
	pos += d.route.GetSize()
	copy(buf[pos:], d.body)
	return nil
}

func (d *Packet) ReadFrom(buf []byte) error {
	pos := 0
	if err := d.head.ReadFrom(buf); err != nil {
		return err
	}
	pos += d.head.GetSize()
	if err := d.route.ReadFrom(buf[pos:]); err != nil {
		return err
	}
	pos += d.route.GetSize()
	d.body = buf[pos:]
	return nil
}
