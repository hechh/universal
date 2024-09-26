package player

import (
	"context"
	"corps/base"
	"corps/framework/actor"
	"corps/server/game/playerMgr/player"
	"reflect"
)

type PlayerMgr struct {
	actor.Actor
	actor.VirtualActor
}

func (p *PlayerMgr) Init() {
	p.Actor.Init()
	p.InitActor(p, reflect.TypeOf(player.Player{}))
	actor.MGR.RegisterActor(p)
	p.Actor.Start()
}

// 玩家登录
func (p *PlayerMgr) LoginPlayerRequest(ctx context.Context) {
	head := p.GetRpcHead(ctx)
	accountId := head.Id

	//加载角色
	if p.GetIActor(accountId) != nil {
		actor.MGR.SendMsgTo(base.NewActorRpcHead(head, "Player", "ReLogin"))
		return
	}

	//创建actor
	pPlayer := &player.Player{}
	pPlayer.AccountId = accountId
	pPlayer.SetId(accountId)
	pPlayer.Init()
	p.AddActor(pPlayer)

	//登录
	actor.MGR.SendMsgTo(base.NewActorRpcHead(head, "Player", "Login"))
}
