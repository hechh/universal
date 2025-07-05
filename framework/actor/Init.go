package actor

import (
	"reflect"
	"strings"
	"universal/common/pb"
	"universal/framework/domain"
	"universal/library/timer"
	"universal/library/uerror"
)

var (
	actors = make(map[string]domain.IActor)
	t      = timer.NewTimer(4)
)

func Register(ac domain.IActor) {
	actors[ac.GetActorName()] = ac
}

func SendMsg(head *pb.Head, args ...interface{}) error {
	if act, ok := actors[head.ActorName]; ok {
		return act.SendMsg(head, args...)
	}
	return uerror.N(1, -1, "%v", head)
}

func Send(head *pb.Head, buf []byte) error {
	if act, ok := actors[head.ActorName]; ok {
		return act.Send(head, buf)
	}
	return uerror.N(1, -1, "%v", head)
}

func parseName(rr reflect.Type) string {
	name := rr.String()
	if index := strings.Index(name, "."); index > -1 {
		name = name[index+1:]
	}
	return name
}
