package actor

import (
	"universal/common/pb"
	"universal/framework/actor/domain"
	"universal/framework/actor/internal/manager"
)

func SetActorHandle(h domain.ActorHandle) {
	manager.SetActorHandle(h)
}

func Send(key string, pac *pb.Packet) {
	manager.Send(key, pac)
}

func GetIActor(key string) domain.IActor {
	return manager.GetIActor(key)
}
