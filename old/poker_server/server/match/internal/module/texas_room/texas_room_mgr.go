package texas_room

import (
	"poker_server/common/config/repository/texas_config"
	"poker_server/common/pb"
	"poker_server/common/room_util"
	"poker_server/framework"
	"poker_server/framework/actor"
	"poker_server/framework/cluster"
	"poker_server/library/mlog"
	"reflect"
)

type MatchTexasRoomMgr struct {
	actor.Actor
	mgr *actor.ActorMgr
}

func NewMatchTexasRoomMgr() *MatchTexasRoomMgr {
	mgr := new(actor.ActorMgr)
	rr := &MatchTexasRoom{}
	mgr.Register(rr)
	mgr.ParseFunc(reflect.TypeOf(rr))
	actor.Register(mgr)

	ret := &MatchTexasRoomMgr{mgr: mgr}
	ret.Actor.Register(ret)
	ret.Actor.ParseFunc(reflect.TypeOf(ret))
	ret.Actor.SetId(uint64(pb.DataType_DataTypeTexasRoom))
	ret.Actor.Start()
	actor.Register(ret)
	return ret
}

func (m *MatchTexasRoomMgr) Close() {
	m.mgr.Stop()
	m.Actor.Stop()
	mlog.Infof("MatchTexasRoomMgr关闭成功")
}

func (m *MatchTexasRoomMgr) Load() error {
	return cluster.Send(&pb.Head{
		Dst: framework.NewDbRouter(uint64(pb.DataType_DataTypeTexasRoom), "DbTexasRoomMgr", "Query"),
		Src: framework.NewSrcRouter(m.GetId(), "MatchTexasRoomMgr", "LoadComplete"),
	}, &pb.GetTexasRoomDataReq{DataType: pb.DataType_DataTypeTexasRoom})
}

func (m *MatchTexasRoomMgr) LoadComplete(head *pb.Head, rsp *pb.GetTexasRoomDataRsp) {
	tmps := map[uint64][]*pb.TexasRoomData{}
	for i, item := range rsp.List {
		mlog.Infof("加载房间数据%d, %v", i, item)
		cfg := texas_config.MGetID(item.GameId)
		id := room_util.ToMatchGameId(cfg.MatchType, cfg.GameType, cfg.CoinType)
		tmps[id] = append(tmps[id], item)

		newHead := &pb.Head{
			Src: framework.NewSrcRouter(uint64(pb.DataType_DataTypeTexasRoom), "MatchTexasRoomMgr"),
			Dst: framework.NewRoomRouter(item.RoomId, "TexasGameMgr", "TexasRoomRestart"),
		}
		err := cluster.Send(newHead, &pb.TexasQueryRsp{Data: item})
		mlog.Infof("主动重启房间: %v, error:%v", item, err)
	}

	for id, list := range tmps {
		m.mgr.AddActor(NewMatchTexasRoom(id, list...))
	}

	// 补全数据
	texas_config.Walk(func(cfg *pb.TexasConfig) bool {
		id := room_util.ToMatchGameId(cfg.MatchType, cfg.GameType, cfg.CoinType)
		if _, ok := tmps[id]; ok {
			return true
		}
		m.mgr.AddActor(NewMatchTexasRoom(id))
		return true
	})
}
