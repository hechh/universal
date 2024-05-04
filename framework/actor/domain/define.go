package domain

import (
	"universal/common/pb"
	"universal/framework/fbasic"
)

// 默认实现 (base/actor)
type IActor interface {
	Start()
	Stop()
	Send(*pb.Packet)
	SetObject(string, fbasic.IData) error
	UUID() string
	GetUpdateTime() int64
	SetUpdateTime(int64)
}

type ActorHandle func(*fbasic.Context, []byte) func()
