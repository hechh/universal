package internal

import (
	"poker_server/common/pb"
	"poker_server/framework/actor"
	"poker_server/library/async"
	"poker_server/library/mlog"
	"poker_server/library/uerror"
	"poker_server/server/game/internal/player"
	"reflect"
)

var (
	playerMgr = NewPlayerMgr()
)

type PlayerMgr struct {
	actor.Actor
	mgr *actor.ActorMgr
}

func Init() error {
	return nil
}

func GetPlayerMgr() *PlayerMgr {
	return playerMgr
}

func NewPlayerMgr() *PlayerMgr {
	mgr := new(actor.ActorMgr)
	pp := &player.Player{}
	mgr.Register(pp)
	mgr.ParseFunc(reflect.TypeOf(pp))
	actor.Register(mgr)

	ret := &PlayerMgr{mgr: mgr}
	ret.Actor.Register(ret)
	ret.Actor.ParseFunc(reflect.TypeOf(ret))
	ret.Start()
	actor.Register(ret)
	return ret
}

func (p *PlayerMgr) Stop() {
	p.mgr.Stop()
	p.Actor.Stop()
}

func (p *PlayerMgr) Kick(uid uint64) {
	act := p.mgr.GetActor(uid)
	if act == nil {
		return
	}
	p.mgr.DelActor(uid)

	p.Add(1)
	async.SafeGo(mlog.Errorf, func() {
		act.Stop()
		p.Done()
	})
}

// 登录请求
func (p *PlayerMgr) Login(head *pb.Head, req *pb.GateLoginRequest, rsp *pb.GateLoginResponse) error {
	if act := p.mgr.GetActor(head.Uid); act != nil {
		head.FuncName = "Relogin"
		return act.SendMsg(head, req, rsp)
	}
	if req.PlayerData == nil {
		return uerror.NEW(pb.ErrorCode_NIL_POINTER, head, "玩家数据为空: %v", req)
	}

	usr := player.NewPlayer(head.Uid, req.PlayerData)
	usr.Start()
	p.mgr.AddActor(usr)
	return usr.SendMsg(head, req, rsp)
}

// 加入德州请求
func (p *PlayerMgr) TexasJoinRoomReq(head *pb.Head, req *pb.TexasJoinRoomReq, rsp *pb.TexasJoinRoomRsp) error {
	if act := p.mgr.GetActor(head.Uid); act != nil {
		return act.SendMsg(head, req, rsp)
	}
	return uerror.NEW(pb.ErrorCode_GAME_PLAYER_NOT_LOGIN, head, "玩家未登录")
}

func (p *PlayerMgr) TexasQuitRoomReq(head *pb.Head, req *pb.TexasQuitRoomReq, rsp *pb.TexasQuitRoomRsp) error {
	if act := p.mgr.GetActor(head.Uid); act != nil {
		return act.SendMsg(head, req, rsp)
	}
	return uerror.NEW(pb.ErrorCode_GAME_PLAYER_NOT_LOGIN, head, "玩家未登录")
}

func (p *PlayerMgr) TexasBuyInReq(head *pb.Head, req *pb.TexasBuyInReq, rsp *pb.TexasBuyInRsp) error {
	if act := p.mgr.GetActor(head.Uid); act != nil {
		return act.SendMsg(head, req, rsp)
	}
	return uerror.NEW(pb.ErrorCode_GAME_PLAYER_NOT_LOGIN, head, "玩家未登录")
}

// RummyJoinRoomReq 加入rummy请求
func (p *PlayerMgr) RummyJoinRoomReq(head *pb.Head, req *pb.RummyJoinRoomReq, rsp *pb.RummyJoinRoomRsp) error {
	if act := p.mgr.GetActor(head.Uid); act != nil {
		return act.SendMsg(head, req, rsp)
	}
	return uerror.NEW(pb.ErrorCode_GAME_PLAYER_NOT_LOGIN, head, "玩家未登录")
}

func (p *PlayerMgr) RummyQuitRoomReq(head *pb.Head, req *pb.RummyQuitRoomReq, rsp *pb.RummyQuitRoomRsp) error {
	if act := p.mgr.GetActor(head.Uid); act != nil {
		return act.SendMsg(head, req, rsp)
	}
	return uerror.NEW(pb.ErrorCode_GAME_PLAYER_NOT_LOGIN, head, "玩家未登录")
}

// RummyChangeRoomReq 换桌
func (p *PlayerMgr) RummyChangeRoomReq(head *pb.Head, req *pb.RummyChangeRoomReq, rsp *pb.RummyChangeRoomRsp) error {
	if act := p.mgr.GetActor(head.Uid); act != nil {
		return act.SendMsg(head, req, rsp)
	}
	return uerror.NEW(pb.ErrorCode_GAME_PLAYER_NOT_LOGIN, head, "玩家未登录")
}

// QueryPlayerData 重连信息查询
func (p *PlayerMgr) QueryPlayerData(head *pb.Head, req *pb.QueryPlayerDataReq, rsp *pb.QueryPlayerDataRsp) error {
	if act := p.mgr.GetActor(head.GetActorId()); act != nil {
		return act.SendMsg(head, req, rsp)
	}
	mlog.Infof("==========================head:%v", head)
	return uerror.NEW(pb.ErrorCode_GAME_PLAYER_NOT_LOGIN, head, "玩家未登录")
}
