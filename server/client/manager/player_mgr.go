package manager

import (
	"reflect"
	"universal/common/pb"
	"universal/common/yaml"
	"universal/framework/actor"
	"universal/library/async"
	"universal/library/mlog"
	"universal/server/client/player"
)

type PlayerMgr struct {
	actor.Actor
	mgr  *actor.ActorMgr
	cfg  *yaml.ServerConfig
	node *pb.Node
}

func NewPlayerMgr(node *pb.Node, cfg *yaml.ServerConfig) *PlayerMgr {
	mgr := new(actor.ActorMgr)
	pl := &player.Player{}
	mgr.Register(pl)
	mgr.ParseFunc(reflect.TypeOf(pl))
	actor.Register(mgr)

	ret := new(PlayerMgr)
	ret.Actor.Register(ret)
	ret.Actor.ParseFunc(reflect.TypeOf(ret))
	ret.Start()
	actor.Register(ret)

	ret.mgr = mgr
	ret.cfg = cfg
	ret.node = node
	return ret
}

func (p *PlayerMgr) Remove(uid uint64) {
	if ac := p.mgr.GetActor(uid); ac != nil {
		p.mgr.DelActor(uid)
		async.SafeGo(mlog.Errorf, func() {
			ac.Stop()
		})
	}
}

func (p *PlayerMgr) Login(begin, end uint64) {
	for i := begin; i <= end; i++ {
		pl := player.NewPlayer(p.node, p.cfg, i)
		p.mgr.AddActor(pl)
		pl.SendMsg(&pb.Head{FuncName: "Login"})
	}
}
