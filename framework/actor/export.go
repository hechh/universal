package actor

import (
	"universal/common/pb"
	"universal/framework/actor/domain"
	"universal/framework/actor/internal/base"
	"universal/framework/actor/internal/manager"
)

func SetActorHandle(h domain.ActorHandle) {
	manager.SetActorHandle(h)
}

func Send(key string, pac *pb.Packet) {
	manager.GetIActor(key).Send(pac)
}

func NewActor(uuid string, h domain.ActorHandle) *base.Actor {
	return base.NewActor(uuid, h)
}

func LoadActor(uuid string) *base.Actor {
	return manager.LoadActor(uuid)
}

func StoreActor(aa *base.Actor) {
	manager.StoreActor(aa)
}
