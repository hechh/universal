package actor

import (
	"poker_server/common/pb"
	"poker_server/framework/domain"
	"poker_server/framework/library/uerror"
)

var (
	busObj domain.IBus
	actors = make(map[string]domain.IActor)
)

func Init(bus domain.IBus) {
	busObj = bus
}

func Register(ac domain.IActor) {
	actors[ac.GetActorName()] = ac
}

func SendMsg(head *pb.Head, args ...interface{}) error {
	if act, ok := actors[head.ActorName]; ok {
		return act.SendMsg(head, args...)
	}
	return uerror.New(1, -1, "Actor%s不存在", head.ActorName)
}

func Send(head *pb.Head, body []byte) error {
	if act, ok := actors[head.ActorName]; ok {
		return act.Send(head, body)
	}
	return uerror.New(1, -1, "Actor%s不存在", head.ActorName)
}
