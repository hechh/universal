package manager

import (
	"poker_server/common/pb"
	"poker_server/common/room_util"
	"poker_server/framework"
	"poker_server/framework/actor"
	"poker_server/framework/cluster"
	"poker_server/library/mlog"
	"poker_server/library/safe"
	"poker_server/library/uerror"
	"poker_server/server/room/internal/internal/sng"
	"poker_server/server/room/internal/internal/texas"
	"reflect"
)

type TexasGameMgr struct {
	actor.Actor
	mgr    *actor.ActorMgr
	sngMgr *actor.ActorMgr
}

func NewTexasGameMgr() *TexasGameMgr {
	sngMgr := new(actor.ActorMgr)
	sngRoom := &sng.SngTexasGame{}
	sngMgr.Register(sngRoom)
	sngMgr.ParseFunc(reflect.TypeOf(sngRoom))
	actor.Register(sngMgr)

	mgr := new(actor.ActorMgr)
	game := &texas.TexasGame{}
	mgr.Register(game)
	mgr.ParseFunc(reflect.TypeOf(game))
	actor.Register(mgr)

	ret := &TexasGameMgr{mgr: mgr, sngMgr: sngMgr}
	ret.Actor.Register(ret)
	ret.Actor.ParseFunc(reflect.TypeOf(ret))
	ret.Actor.SetId(uint64(pb.DataType_DataTypeTexasRoom))
	ret.Actor.Start()
	actor.Register(ret)
	return ret
}

func (d *TexasGameMgr) Stop() {
	d.sngMgr.Stop()
	d.mgr.Stop()
	d.Actor.Stop()
	mlog.Infof("TexasGameMgr关闭成功")
}

func (d *TexasGameMgr) Remove(id uint64) {
	if act := d.mgr.GetActor(id); act != nil {
		d.mgr.DelActor(id)
		d.Add(1)
		safe.Go(func() {
			act.Stop()
			d.Done()
		})
	}
}

func (d *TexasGameMgr) SngRemove(id uint64) {
	if act := d.sngMgr.GetActor(id); act != nil {
		d.sngMgr.DelActor(id)
		d.Add(1)
		safe.Go(func() {
			act.Stop()
			d.Done()
		})
	}
}

// 服务重启需要重新拉起房间
func (d *TexasGameMgr) TexasRoomRestart(head *pb.Head, rsp *pb.TexasQueryRsp) {
	if act := d.mgr.GetActor(rsp.Data.RoomId); act != nil {
		return
	}
	rr := texas.NewTexasGame(rsp.Data)
	d.mgr.AddActor(rr)
	rr.Init()
}

// 加入房间请求
func (d *TexasGameMgr) JoinRoomReq(head *pb.Head, req *pb.TexasJoinRoomReq, rsp *pb.TexasJoinRoomRsp) error {
	matchType, gameType, coinType := room_util.TexasRoomIdTo(req.RoomId)
	switch matchType {
	case pb.MatchType_MatchTypeNone:
		if act := d.mgr.GetActor(req.RoomId); act != nil {
			return act.SendMsg(head, req, rsp)
		}

		id := room_util.ToMatchGameId(matchType, gameType, coinType)
		newHead := &pb.Head{
			Uid: head.Uid,
			Src: framework.NewSrcRouter(req.RoomId, d.GetActorName()),
			Dst: framework.NewMatchRouter(id, "MatchTexasRoom", "Query"),
		}
		newreq := &pb.TexasQueryReq{RoomId: req.RoomId}
		newrsp := &pb.TexasQueryRsp{}
		if err := cluster.Request(newHead, newreq, newrsp); err != nil {
			return err
		}
		if newrsp.Head != nil {
			return uerror.ToError(newrsp.Head)
		}
		// 创建房间
		rr := texas.NewTexasGame(newrsp.Data)
		d.mgr.AddActor(rr)
		rr.Init()
		return rr.SendMsg(head, req, rsp)
	case pb.MatchType_MatchTypeSNG:
		if act := d.sngMgr.GetActor(req.RoomId); act != nil {
			return act.SendMsg(head, req, rsp)
		}
		// 请求房间数据
		newHead := &pb.Head{
			Src: framework.NewSrcRouter(d.GetId(), d.GetActorName()),
			Dst: framework.NewMatchRouter(uint64(pb.DataType_DataTypeSngRoom), "SngRoomMgr", "Query"),
		}
		data := &pb.TexasRoomData{}
		if err := cluster.Request(newHead, req.RoomId, data); err != nil {
			return err
		}
		// 创建房间
		rr := sng.NewSngTexasGame(data)
		d.sngMgr.AddActor(rr)
		return rr.SendMsg(head, req, rsp)
	}
	return uerror.New(1, pb.ErrorCode_ACTOR_ID_NOT_FOUND, "actor不存在")
}

func (d *TexasGameMgr) QuitRoomReq(head *pb.Head, req *pb.TexasQuitRoomReq, rsp *pb.TexasQuitRoomRsp) error {
	matchType := pb.MatchType((req.RoomId >> 40) & 0xFF)
	switch matchType {
	case pb.MatchType_MatchTypeNone:
		if act := d.mgr.GetActor(req.RoomId); act != nil {
			return act.SendMsg(head, req, rsp)
		}
	case pb.MatchType_MatchTypeSNG:
		if act := d.sngMgr.GetActor(req.RoomId); act != nil {
			return act.SendMsg(head, req, rsp)
		}
	}
	d.kick(head)
	return uerror.New(1, pb.ErrorCode_ACTOR_ID_NOT_FOUND, "actor不存在")
}

func (d *TexasGameMgr) kick(head *pb.Head) {
	cluster.SendToClient(&pb.Head{
		Uid: head.Uid,
		Src: framework.NewSrcRouter(d.GetId(), d.GetActorName()),
		Dst: framework.NewGateRouter(head.Uid, "Player", "SendToClient"),
	}, &pb.TexasQuitRoomNotify{})
}

func (d *TexasGameMgr) SitDownReq(head *pb.Head, req *pb.TexasSitDownReq, rsp *pb.TexasSitDownRsp) error {
	matchType := pb.MatchType((req.RoomId >> 40) & 0xFF)
	switch matchType {
	case pb.MatchType_MatchTypeNone:
		if act := d.mgr.GetActor(req.RoomId); act != nil {
			return act.SendMsg(head, req, rsp)
		}
	case pb.MatchType_MatchTypeSNG:
		if act := d.sngMgr.GetActor(req.RoomId); act != nil {
			return act.SendMsg(head, req, rsp)
		}
	}
	d.kick(head)
	return uerror.New(1, pb.ErrorCode_ACTOR_ID_NOT_FOUND, "actor不存在")
}

func (d *TexasGameMgr) StandUpReq(head *pb.Head, req *pb.TexasStandUpReq, rsp *pb.TexasStandUpRsp) error {
	matchType := pb.MatchType((req.RoomId >> 40) & 0xFF)
	switch matchType {
	case pb.MatchType_MatchTypeNone:
		if act := d.mgr.GetActor(req.RoomId); act != nil {
			return act.SendMsg(head, req, rsp)
		}
	case pb.MatchType_MatchTypeSNG:
		if act := d.sngMgr.GetActor(req.RoomId); act != nil {
			return act.SendMsg(head, req, rsp)
		}
	}
	d.kick(head)
	return uerror.New(1, pb.ErrorCode_ACTOR_ID_NOT_FOUND, "actor不存在")
}

func (d *TexasGameMgr) BuyInReq(head *pb.Head, req *pb.TexasBuyInReq, rsp *pb.TexasBuyInRsp) error {
	matchType := pb.MatchType((req.RoomId >> 40) & 0xFF)
	switch matchType {
	case pb.MatchType_MatchTypeNone:
		if act := d.mgr.GetActor(req.RoomId); act != nil {
			return act.SendMsg(head, req, rsp)
		}
	case pb.MatchType_MatchTypeSNG:
		if act := d.sngMgr.GetActor(req.RoomId); act != nil {
			return act.SendMsg(head, req, rsp)
		}
	}
	d.kick(head)
	return uerror.New(1, pb.ErrorCode_ACTOR_ID_NOT_FOUND, "actor不存在")
}

func (d *TexasGameMgr) DoBetReq(head *pb.Head, req *pb.TexasDoBetReq, rsp *pb.TexasDoBetRsp) error {
	matchType := pb.MatchType((req.RoomId >> 40) & 0xFF)
	switch matchType {
	case pb.MatchType_MatchTypeNone:
		if act := d.mgr.GetActor(req.RoomId); act != nil {
			return act.SendMsg(head, req, rsp)
		}
	case pb.MatchType_MatchTypeSNG:
		if act := d.sngMgr.GetActor(req.RoomId); act != nil {
			return act.SendMsg(head, req, rsp)
		}
	}
	d.kick(head)
	return uerror.New(1, pb.ErrorCode_ACTOR_ID_NOT_FOUND, "actor不存在")
}

func (d *TexasGameMgr) StatisticsReq(head *pb.Head, req *pb.TexasStatisticsReq, rsp *pb.TexasStatisticsRsp) error {
	matchType := pb.MatchType((req.RoomId >> 40) & 0xFF)
	switch matchType {
	case pb.MatchType_MatchTypeNone:
		if act := d.mgr.GetActor(req.RoomId); act != nil {
			return act.SendMsg(head, req, rsp)
		}
	case pb.MatchType_MatchTypeSNG:
		/*
			if act := d.sngMgr.GetActor(req.RoomId); act != nil {
				return act.SendMsg(head, req, rsp)
			}
		*/
	}
	d.kick(head)
	return uerror.New(1, pb.ErrorCode_ACTOR_ID_NOT_FOUND, "actor不存在")
}
