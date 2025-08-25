package actor

import (
	"reflect"
	"strings"
	"universal/common/pb"
	"universal/framework/define"
	"universal/library/timer"
	"universal/library/uerror"
)

var (
	self    *pb.Node
	actors  = make(map[string]define.IActor)
	t       = timer.NewTimer(4)
	sendrsp define.SendRspFunc
)

func Init(nn *pb.Node, srsp define.SendRspFunc) {
	self = nn
	sendrsp = srsp
}

func Register(ac define.IActor) {
	actors[ac.GetActorName()] = ac
}

func SendMsg(head *pb.Head, args ...interface{}) error {
	if act, ok := actors[head.ActorName]; ok {
		return act.SendMsg(head, args...)
	}
	return uerror.New(1, -1, "%v", head)
}

func Send(head *pb.Head, buf []byte) error {
	if act, ok := actors[head.ActorName]; ok {
		return act.Send(head, buf)
	}
	return uerror.New(1, -1, "%v", head)
}

func parseName(rr reflect.Type) string {
	name := rr.String()
	if index := strings.Index(name, "."); index > -1 {
		name = name[index+1:]
	}
	return name
}
