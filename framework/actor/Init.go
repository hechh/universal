package actor

import (
	"reflect"
	"strings"
	"universal/common/pb"
	"universal/framework/define"
	"universal/framework/handler"
	"universal/library/timer"
	"universal/library/uerror"
)

var (
	actors  = make(map[string]define.IActor)
	t       = timer.NewTimer(4)
	self    *pb.Node
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

func RpcCall(head *pb.Head, buf []byte) error {
	head.ActorName, head.FuncName = handler.GetActorFunc(head.Dst.ActorFunc)
	if head.Dst.ActorId > 0 {
		head.ActorId = head.Dst.ActorId
	} else {
		head.ActorId = handler.ParseRouterId(head.Dst.RouterId)
	}
	if act, ok := actors[head.ActorName]; ok {
		return act.Send(head, buf)
	}
	return uerror.New(1, -1, "接口%s.%s未注册", head.ActorName, head.FuncName)
}

func parseName(rr reflect.Type) string {
	name := rr.String()
	if index := strings.Index(name, "."); index > -1 {
		name = name[index+1:]
	}
	return name
}
