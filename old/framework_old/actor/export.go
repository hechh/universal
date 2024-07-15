package actor

import (
	"universal/common/pb"
	"universal/framework/actor/domain"
	"universal/framework/actor/internal/manager"
)

func SetActorClearExpire(expire int64) {
	manager.SetClearExpire(expire)
}

func GetIActor(key string, ff domain.ActorHandle) domain.IActor {
	return manager.GetIActor(key, ff)
}

func Send(key string, ff domain.ActorHandle, pac *pb.Packet) {
	manager.Send(key, ff, pac)
}

func StopAll() {
	manager.StopAll()
}
