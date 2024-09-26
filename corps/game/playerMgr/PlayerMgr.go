package playerMgr

import (
	"context"
	"corps/base"
	"corps/common"
	"corps/common/serverCommon"
	"corps/framework/actor"
	"corps/framework/cluster"
	"corps/framework/plog"
	"corps/pb"
	"corps/server/game/playerMgr/player"
	"reflect"
)

// ********************************************************
// 玩家管理
// ********************************************************
var (
	MGR PlayerMgr
)

type (
	PlayerMgr struct {
		actor.Actor
		actor.VirtualActor
	}

	IPlayerMgr interface {
		Update()
	}
)

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

func (p *PlayerMgr) OnStubRegister(ctx context.Context) {
	//这里可以是加载db数据
	plog.Info("Stub db register sucess id:%d", p.GetId())
}

func (p *PlayerMgr) OnStubUnRegister(ctx context.Context) {
	//lease一致性这里要清理缓存数据了
	plog.Info("Stub db unregister sucess id:%d", p.GetId())
}

// mail注销
func (p *PlayerMgr) OnMailUnRegister(ctx context.Context, id uint64) {
	head := p.GetRpcHead(ctx)
	head.Id = id

	newHead := base.DeepCopyRpcHead(head)
	newHead.Id = id
	newHead.ActorName = "Player"
	newHead.FuncName = "OnMailUnRegister"
	actor.MGR.SendMsgTo(newHead)
}

// 停服
func (this *PlayerMgr) ShutDown() {
	// 不等定时器触发，立即设施保存任务
	this.StopAll()
	plog.Info(" (this *PlayerMgr) ShutDown() finished")
}

// actor删除
func (this *PlayerMgr) OnStopSon(ctx context.Context) {
	head := this.GetRpcHead(ctx)
	this.DelActor(head.Id)
	plog.Info("OnStopSon id:%d", head.Id)
}

// actor删除
func (this *PlayerMgr) DelPlayer(ctx context.Context) {
	head := this.GetRpcHead(ctx)
	ac := this.GetIActor(head.Id)
	if ac == nil {
		return
	}
	ac.Stop()
	this.DelActor(head.Id)
	plog.Info("DelPlayer id:%d", head.Id)
}

// 发送奖励
func (this *PlayerMgr) RewardRequest(ctx context.Context, req *pb.RankRewardRequest) {
	head := this.GetRpcHead(ctx)
	if ac := this.GetIActor(head.Id); ac != nil {
		head.ActorName = "Player"
		head.FuncName = "RewardRequest"
		ac.SendMsg(head, req)
	} else {
		// 写入离线缓存中
		serverCommon.SendPlayerItem(head.Id, common.AddItemToAddItemData(req.Doing, req.Rewards...)...)
	}
	// 给dip回包
	if len(head.Reply) > 0 {
		cluster.ReplyMsgTo(head)
	}
}
