package repository

import (
	"universal/common/pb"
	"universal/framework/basic"
	"universal/framework/packet/domain"
)

type index struct {
	ActorName string
	FuncName  string
}

type ActorPacket struct {
	apis map[index]*domain.IPacket
}

func (d *ActorPacket) Call(ctx *basic.Context, pac *pb.Packet) (*pb.Packet, error) {
	return nil, nil
}
