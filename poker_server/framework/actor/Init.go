package actor

import (
	"poker_server/common/pb"
	"poker_server/framework/domain"
	"poker_server/library/timer"
	"poker_server/library/uerror"
	"reflect"
	"strings"
	"time"
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
	return uerror.New(1, pb.ErrorCode_ACTOR_NOT_SUPPORTED, "Actor(%s)不存在", head.ActorName)
}

func Send(head *pb.Head, body []byte) error {
	if act, ok := actors[head.ActorName]; ok {
		return act.Send(head, body)
	}
	return uerror.New(1, pb.ErrorCode_ACTOR_NOT_SUPPORTED, "Actor%s不存在", head.ActorName)
}

func RegisterTimer(id *uint64, f func(), ttl time.Duration, times int32) error {
	return t.Register(id, f, ttl, times)
}

func parseName(rr reflect.Type) string {
	name := rr.String()
	if index := strings.Index(name, "."); index > -1 {
		name = name[index+1:]
	}
	return name
}
