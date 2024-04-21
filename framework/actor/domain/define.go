package domain

import (
	"universal/common/pb"
	"universal/framework/basic"
)

const (
	SessionExpireTime = 30 * 60 //单位：秒
)

type ISession interface {
	Send(*pb.Packet)
	SetActor(string, interface{}) error
}

type PacketHandle func(*basic.Context, *pb.Packet) func()
