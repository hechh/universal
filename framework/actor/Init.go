package actor

import (
	"reflect"
	"strings"
	"universal/common/pb"
	"universal/framework/domain"
	"universal/library/timer"
	"universal/library/uerror"
	"universal/library/util"
)

var (
	names  = make(map[string]uint32)
	apis   = make(map[uint32]domain.IFuncs)
	actors = make(map[string]domain.IActor)
	t      = timer.NewTimer(4)
)

func GetCrc32(actorFunc string) uint32 {
	if _, ok := names[actorFunc]; !ok {
		names[actorFunc] = util.GetCrc32(actorFunc)
	}
	return names[actorFunc]
}

func Parse(head *pb.Head, ffs ...string) error {
	var ok bool
	var rr domain.IFuncs
	if head.Dst.ActorFunc > 0 {
		rr, ok = apis[head.Dst.ActorFunc]
	} else if len(ffs) > 0 {
		rr, ok = apis[GetCrc32(ffs[0])]
	}
	if !ok {
		return uerror.New(1, -1, "请求接口不存在%v", head.Dst)
	}
	return rr.Parse(head)
}

func Register(ac domain.IActor) {
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
