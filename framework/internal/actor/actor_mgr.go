package actor

import (
	"universal/common/pb"
	"universal/framework/domain"
	"universal/framework/library/uerror"
)

type ActorMgr struct {
	actors map[string]domain.IActor
}

func NewActorMgr() *ActorMgr {
	return &ActorMgr{
		actors: make(map[string]domain.IActor),
	}
}

func (mgr *ActorMgr) Register(ac domain.IActor) {
	mgr.actors[ac.GetActorName()] = ac
}

func (mgr *ActorMgr) Get(name string) domain.IActor {
	if ac, ok := mgr.actors[name]; ok {
		return ac
	}
	return nil
}

func (mgr *ActorMgr) Send(head *pb.Head, args ...interface{}) error {
	if ac, ok := mgr.actors[head.ActorName]; ok {
		return ac.Send(head, args...)
	}
	return uerror.New(1, -1, "%s.%s未实现", head.ActorName, head.FuncName)
}

func (mgr *ActorMgr) SendRpc(head *pb.Head, data []byte) error {
	if ac, ok := mgr.actors[head.ActorName]; ok {
		return ac.SendRpc(head, data)
	}
	return uerror.New(1, -1, "%s.%s未实现", head.ActorName, head.FuncName)
}
