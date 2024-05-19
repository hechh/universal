package domain

import (
	"universal/common/pb"
	"universal/framework/common/fbasic"
)

// 默认实现 (base/actor)
type IActor interface {
	Start()
	Stop()
	GetUID() string
	GetUpdateTime() int64
	SetUpdateTime(int64)
	Send(*pb.PacketHead, []byte)
}

type ActorHandle func(*fbasic.Context, []byte) func()
