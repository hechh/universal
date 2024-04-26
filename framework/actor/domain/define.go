package domain

import (
	"universal/common/pb"
	"universal/framework/fbasic"
)

type IActor interface {
	Start()
	Stop()
	Send(*pb.Packet)
	SetObject(string, fbasic.IData) error
}

type ActorHandle func(*fbasic.Context, []byte) func()
