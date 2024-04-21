package repository

import (
	"fmt"
	"universal/common/pb"
	"universal/framework/basic"
)

type EmptyPacket struct {
	actorName string
	name      string
}

func NewEmptyPacket(a, b string) *EmptyPacket {
	return &EmptyPacket{a, b}
}

func (d *EmptyPacket) Call(ctx *basic.Context, pac *pb.Packet) *pb.Packet {
	return &pb.Packet{
		Head:   pac.Head,
		Code:   int32(pb.ErrorCode_NotSupported),
		ErrMsg: fmt.Sprintf("%s.%s() not supported", d.actorName, d.name),
	}
}
