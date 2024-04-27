package actor

import (
	"universal/common/pb"
	"universal/framework/actor/domain"
	"universal/framework/actor/internal/base"
	"universal/framework/actor/internal/manager"
	"universal/framework/fbasic"
)

func SetActorHandle(h domain.ActorHandle) {
	manager.SetActorHandle(h)
}

func Send(key string, pac *pb.Packet) {
	manager.Send(key, pac)
}

func NewActor(uuid string, h domain.ActorHandle) domain.IActor {
	return base.NewActor(uuid, h)
}

func Load(uuid string) domain.IActor {
	return manager.Load(uuid)
}

func Store(aa interface{}) error {
	switch vv := aa.(type) {
	case *base.Actor:
		manager.Store(vv)
	case domain.ICustom:
		manager.Store(vv)
	default:
		return fbasic.NewUError(1, pb.ErrorCode_TypeNotSupported, "*Actor or IMgrActor is expected")
	}
	return nil
}
