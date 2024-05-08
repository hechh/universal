package base

import (
	"universal/common/pb"
	"universal/framework/notify/domain"
)

type Unicast struct {
	key string
	f   domain.NotifyHandle
}

func NewUnicast(key string, f domain.NotifyHandle) *Unicast {
	return &Unicast{
		key: key,
		f:   f,
	}
}

func (d *Unicast) GetKey() string {
	return d.key
}

func (d *Unicast) Handle(pac *pb.Packet) {
	d.f(pac)
}
