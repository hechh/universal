package actor

import (
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

func (mgr *ActorMgr) Send(head domain.IHead, args ...interface{}) error {
	if ac, ok := mgr.actors[head.GetActorName()]; ok {
		return ac.Send(head, args...)
	}
	return uerror.New(1, -1, "%s.%s未实现", head.GetActorName(), head.GetFuncName())
}

func (mgr *ActorMgr) SendRpc(head domain.IHead, data []byte) error {
	if ac, ok := mgr.actors[head.GetActorName()]; ok {
		return ac.SendRpc(head, data)
	}
	return uerror.New(1, -1, "%s.%s未实现", head.GetActorName(), head.GetFuncName())
}
