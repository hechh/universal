package texas_room

import (
	"poker_server/common/config/repository/texas_config"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/actor"
	"reflect"
	"time"
)

type MatchTexasRoomMgr struct {
	actor.Actor
	mgr   *actor.ActorMgr
	datas map[uint64]*pb.TexasRoomData
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
}

func (m *MatchTexasRoomMgr) Load() error {
	dst := framework.NewDbRouter(uint64(pb.DataType_DataTypeTexasRoom), "DbTexasRoomMgr", "Query")
	head := framework.NewHead(dst, pb.RouterType_RouterTypeGameType, m.GetId(), "MatchTexasRoomMgr", "LoadComplete")
	return framework.Send(head, &pb.GetTexasRoomDataReq{DataType: pb.DataType_DataTypeTexasRoom})
}

func (m *MatchTexasRoomMgr) LoadComplete(head *pb.Head, rsp *pb.GetTexasRoomDataRsp) error {
	tmps := map[uint64][]*pb.TexasRoomData{}
	for _, item := range rsp.List {
		cfg := texas_config.MGetID(item.GameId)
		id := uint64(cfg.GameType&0xFFFF)<<16 | uint64(cfg.CoinType&0xFFFF)
		tmps[id] = append(tmps[id], item)
	}
	for id, list := range tmps {
		m.mgr.AddActor(NewMatchTexasRoom(id, list...))
	}

	// 补全数据
	texas_config.Walk(func(cfg *pb.TexasConfig) bool {
		id := uint64(cfg.MatchType&0xFFFF)<<32 | uint64(cfg.GameType&0xFFFF)<<16 | uint64(cfg.CoinType&0xFFFF)
		if _, ok := tmps[id]; ok {
			return true
		}
		m.mgr.AddActor(NewMatchTexasRoom(id))
		return true
	})
	m.datas = make(map[uint64]*pb.TexasRoomData)

	// 注册定时器
	return m.RegisterTimer(&pb.Head{
		SendType:  pb.SendType_POINT,
		ActorName: "MatchTexasRoomMgr",
		FuncName:  "OnTick",
	}, 5*time.Second, -1)
}

// 定时落地到db服务
func (m *MatchTexasRoomMgr) OnTick() {
	if len(m.datas) <= 0 {
		return
	}
	m.SendMsg(&pb.Head{FuncName: "Save"})
}

func (m *MatchTexasRoomMgr) Collect(notify *pb.UpdateTexasRoomDataNotify) {
	for _, item := range notify.List {
		m.datas[item.RoomId] = item
	}
}

// 保存数据
func (m *MatchTexasRoomMgr) Save() error {
	if len(m.datas) <= 0 {
		return nil
	}

	notify := &pb.UpdateTexasRoomDataNotify{DataType: pb.DataType_DataTypeTexasRoom}
	for _, item := range m.datas {
		notify.List = append(notify.List, item)
	}

	dst := framework.NewDbRouter(uint64(pb.DataType_DataTypeTexasRoom), "DbTexasRoomMgr", "Update")
	head := framework.NewHead(dst, pb.RouterType_RouterTypeGameType, uint64(pb.DataType_DataTypeTexasRoom))
	if err := framework.Send(head, notify); err != nil {
		return err
	}

	// 清空数据
	for key := range m.datas {
		delete(m.datas, key)
	}
	return nil
}
