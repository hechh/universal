package domain

import (
	"universal/common/pb"
	"universal/framework/fbasic"
)

// 默认实现 (base/actor)
type IActor interface {
	Start()
	Stop()
	Send(*pb.PacketHead, []byte)
	UUID() string
	GetUpdateTime() int64
	SetUpdateTime(int64)
	SetObject(string, fbasic.IData) error
}

type ActorHandle func(*fbasic.Context, []byte) func()
