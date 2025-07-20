package internal

import (
	"poker_server/common/pb"
	"poker_server/common/room_util"
	"poker_server/framework"
	"poker_server/framework/actor"
	"poker_server/framework/cluster"
	"poker_server/library/mlog"
	"poker_server/library/safe"
	"poker_server/library/snowflake"
	"poker_server/library/uerror"
	"poker_server/library/util"
	"poker_server/server/game/internal/player"
	"poker_server/server/game/module/http"
	"reflect"
	"strconv"
)

var (
	playerMgr = NewPlayerMgr()
)

type PlayerMgr struct {
	actor.Actor
	mgr *actor.ActorMgr
}

func Init() {
	util.Must(cluster.SetBroadcastHandler(framework.DefaultHandler))
	util.Must(cluster.SetSendHandler(framework.DefaultHandler))
	util.Must(cluster.SetReplyHandler(framework.DefaultHandler))
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

func (p *PlayerMgr) Close() {
	p.mgr.Stop()
	p.Actor.Stop()
	mlog.Infof("PlayerMgr关闭成功")
}

func (p *PlayerMgr) Kick(uid uint64) {
	act := p.mgr.GetActor(uid)
	if act == nil {
		return
	}
	p.mgr.DelActor(uid)

	p.Add(1)
	safe.Go(func() {
		act.Stop()
		p.Done()
	})
}

// 登录请求
func (p *PlayerMgr) Login(head *pb.Head, req *pb.GateLoginRequest, rsp *pb.GateLoginResponse) error {
	mlog.Debugf("PlayerDataPool.Login: rsp:%v", req.PlayerData)
	if act := p.mgr.GetActor(head.Uid); act != nil {
		head.FuncName = "Relogin"
		return act.SendMsg(head, req, rsp)
	}
	if req.PlayerData == nil {
		return uerror.New(1, pb.ErrorCode_NIL_POINTER, "玩家数据为空: %v", req)
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
	return uerror.New(1, pb.ErrorCode_GAME_PLAYER_NOT_LOGIN, "玩家未登录")
}

func (p *PlayerMgr) TexasQuitRoomReq(head *pb.Head, req *pb.TexasQuitRoomReq, rsp *pb.TexasQuitRoomRsp) error {
	if act := p.mgr.GetActor(head.Uid); act != nil {
		return act.SendMsg(head, req, rsp)
	}
	return uerror.New(1, pb.ErrorCode_GAME_PLAYER_NOT_LOGIN, "玩家未登录")
}

func (p *PlayerMgr) TexasBuyInReq(head *pb.Head, req *pb.TexasBuyInReq, rsp *pb.TexasBuyInRsp) error {
	if act := p.mgr.GetActor(head.Uid); act != nil {
		return act.SendMsg(head, req, rsp)
	}
	return uerror.New(1, pb.ErrorCode_GAME_PLAYER_NOT_LOGIN, "玩家未登录")
}

func (p *PlayerMgr) TexasChangeReq(head *pb.Head, req *pb.TexasChangeRoomReq, rsp *pb.TexasChangeRoomRsp) error {
	if act := p.mgr.GetActor(head.Uid); act != nil {
		return act.SendMsg(head, req, rsp)
	}
	return uerror.New(1, pb.ErrorCode_GAME_PLAYER_NOT_LOGIN, "玩家未登录")
}

func (p *PlayerMgr) TexasFinishNotify(head *pb.Head, event *pb.TexasFinishNotify) error {
	if act := p.mgr.GetActor(head.Uid); act != nil {
		return act.SendMsg(head, event)
	}
	if event.Incr <= 0 {
		return nil
	}
	// 离线玩家，道具先存入redis保存
	return propMgr.SendMsg(head, event)
}

// RummyJoinRoomReq 加入rummy请求
func (p *PlayerMgr) RummyJoinRoomReq(head *pb.Head, req *pb.RummyJoinRoomReq, rsp *pb.RummyJoinRoomRsp) error {
	if act := p.mgr.GetActor(head.Uid); act != nil {
		return act.SendMsg(head, req, rsp)
	}
	return uerror.New(1, pb.ErrorCode_GAME_PLAYER_NOT_LOGIN, "玩家未登录")
}

func (p *PlayerMgr) RummyQuitRoomReq(head *pb.Head, req *pb.RummyQuitRoomReq, rsp *pb.RummyQuitRoomRsp) error {
	if act := p.mgr.GetActor(head.Uid); act != nil {
		return act.SendMsg(head, req, rsp)
	}
	return uerror.New(1, pb.ErrorCode_GAME_PLAYER_NOT_LOGIN, "玩家未登录")
}

// RummyChangeRoomReq 换桌
func (p *PlayerMgr) RummyChangeRoomReq(head *pb.Head, req *pb.RummyChangeRoomReq, rsp *pb.RummyChangeRoomRsp) error {
	if act := p.mgr.GetActor(head.Uid); act != nil {
		return act.SendMsg(head, req, rsp)
	}
	return uerror.New(1, pb.ErrorCode_GAME_PLAYER_NOT_LOGIN, "玩家未登录")
}

// QueryPlayerData 重连信息查询
func (p *PlayerMgr) QueryPlayerData(head *pb.Head, req *pb.QueryPlayerDataReq, rsp *pb.QueryPlayerDataRsp) error {
	if act := p.mgr.GetActor(head.GetActorId()); act != nil {
		return act.SendMsg(head, req, rsp)
	}
	return cluster.Send(framework.SwapToDb(head, uint64(pb.DataType_DataTypeRoomInfo), "RoomInfoMgr", "Query"), req)
}

// RummyMatchReq Rummy分支玩法开启匹配请求
func (p *PlayerMgr) RummyMatchReq(head *pb.Head, req *pb.RummyMatchReq, rsp *pb.RummyMatchRsp) error {
	if act := p.mgr.GetActor(head.GetActorId()); act != nil {
		return act.SendMsg(head, req, rsp)
	}
	return uerror.New(1, pb.ErrorCode_GAME_PLAYER_NOT_LOGIN, "玩家未登录")
}

func (p *PlayerMgr) RummyGiveUpReq(head *pb.Head, req *pb.RummyGiveUpReq, rsp *pb.RummyGiveUpRsp) error {
	if act := p.mgr.GetActor(head.GetActorId()); act != nil {
		return act.SendMsg(head, req, rsp)
	}
	return uerror.New(1, pb.ErrorCode_GAME_PLAYER_NOT_LOGIN, "玩家未登录")
}

func (p *PlayerMgr) RummyCancelMatchReq(head *pb.Head, req *pb.RummyCancelMatchReq, rsp *pb.RummyCancelMatchRsp) error {
	if act := p.mgr.GetActor(head.GetActorId()); act != nil {
		return act.SendMsg(head, req, rsp)
	}
	return uerror.New(1, pb.ErrorCode_GAME_PLAYER_NOT_LOGIN, "玩家未登录")
}

// RummyKickPlayerReq rummy用户超时情况回补
func (p *PlayerMgr) RummyKickPlayerReq(head *pb.Head, req *pb.RummyKickPlayerReq) {
	if act := p.mgr.GetActor(head.GetUid()); act != nil {
		err := act.SendMsg(head, req)
		mlog.Infof("RummyKickPlayerReq err:%v", err)
		return
	}

	// 如果玩家超时被清理
	_, gameType, coinType := room_util.TexasRoomIdTo(req.GetRoomId())
	uuid, _ := snowflake.GenUUID()
	param := &pb.TransParam{
		GameSn:   strconv.FormatUint(uuid, 10),
		GameType: gameType,
		CoinType: coinType,
		Incr:     req.GetCharge(),
		Uid:      head.Uid,
	}
	rsp := &pb.HttpTransferOutRsp{}
	if err := http.ChargeTransOutRequest(param, rsp); err != nil {
		mlog.Error(head, "ChargeTransOut http err: %v", err)
		return
	}
	mlog.Info(head, "PlayerMgr ChargeTransOut uid:%d, param:%v, rsp:%v", head.Uid, param, rsp)
}
