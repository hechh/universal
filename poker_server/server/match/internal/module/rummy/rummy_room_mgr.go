package rummy

import (
	"poker_server/common/config/repository/rummy_config"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/actor"
	"poker_server/framework/cluster"
	"poker_server/library/mlog"
	"poker_server/library/uerror"
	"reflect"
	"time"
)

type MatchRummyRoomMgr struct {
	actor.Actor
	mgr   *actor.ActorMgr
	datas map[uint64]*pb.RummyRoomData
}

func NewMatchRummyRoomMgr() *MatchRummyRoomMgr {
	mgr := new(actor.ActorMgr)
	rr := &MatchRummyRoom{}
	mgr.Register(rr)
	mgr.ParseFunc(reflect.TypeOf(rr))
	actor.Register(mgr)

	ret := &MatchRummyRoomMgr{mgr: mgr}
	ret.Actor.Register(ret)
	ret.Actor.ParseFunc(reflect.TypeOf(ret))
	ret.Actor.SetId(uint64(pb.DataType_DataTypeRummyRoom))
	ret.Actor.Start()
	actor.Register(ret)
	return ret
}

func (m *MatchRummyRoomMgr) Close() {
	m.mgr.Stop()
	m.Actor.Stop()
	mlog.Infof("MatchRummyRoomMgr关闭成功")
}

// 加载数据
func (m *MatchRummyRoomMgr) Load() error {
	head := &pb.Head{
		Src: framework.NewSrcRouter(m.GetId(), m.GetActorName()),
		Dst: framework.NewDbRouter(uint64(pb.DataType_DataTypeRummyRoom), "DbRummyRoomMgr", "Query"),
	}
	req := &pb.GetRummyRoomDataReq{DataType: pb.DataType_DataTypeRummyRoom}
	rsp := &pb.GetRummyRoomDataRsp{}
	if err := cluster.Request(head, req, rsp); err != nil {
		return err
	}
	if rsp.Head != nil {
		return uerror.ToError(rsp.Head)
	}

	// 初始化
	tmps := map[uint64][]*pb.RummyRoomData{}
	for _, item := range rsp.List {
		cfg := rummy_config.MGetID(item.GameId)
		id := uint64(cfg.GameType)<<32 | uint64(cfg.CoinType)
		tmps[id] = append(tmps[id], item)
	}
	for id, list := range tmps {
		m.mgr.AddActor(NewMatchRummyRoom(id, list...))
	}

	// 补全数据
	rummy_config.Walk(func(cfg *pb.RummyConfig) bool {
		id := uint64(cfg.GameType)<<32 | uint64(cfg.CoinType)
		if _, ok := tmps[id]; ok {
			return true
		}
		m.mgr.AddActor(NewMatchRummyRoom(id))
		return true
	})
	m.datas = make(map[uint64]*pb.RummyRoomData)

	// 注册定时器
	return m.RegisterTimer(&pb.Head{
		SendType:  pb.SendType_POINT,
		ActorName: "MatchRummyRoomMgr",
		FuncName:  "OnTick",
	}, 5*time.Second, -1)
}

// OnTick 定时落地到db服务
func (m *MatchRummyRoomMgr) OnTick() {
	if len(m.datas) <= 0 {
		return
	}
	m.SendMsg(&pb.Head{FuncName: "Save"})
}

func (m *MatchRummyRoomMgr) Collect(notify *pb.UpdateRummyRoomDataNotify) {
	for _, item := range notify.List {
		m.datas[item.RoomId] = item
	}
}

// Save 保存数据
func (m *MatchRummyRoomMgr) Save() error {
	if len(m.datas) <= 0 {
		return nil
	}

	notify := &pb.UpdateRummyRoomDataNotify{DataType: pb.DataType_DataTypeRummyRoom}
	for _, item := range m.datas {
		notify.List = append(notify.List, item)
	}
	head := &pb.Head{
		Src: framework.NewSrcRouter(m.GetId(), m.GetActorName()),
		Dst: framework.NewDbRouter(uint64(pb.DataType_DataTypeRummyRoom), "DbRummyRoomMgr", "Update"),
	}
	if err := cluster.Send(head, notify); err != nil {
		return err
	}

	// 清空数据
	for key := range m.datas {
		delete(m.datas, key)
	}
	return nil
}

func (m *MatchRummyRoomMgr) Remove(head *pb.Head, req *pb.RummyRemoveRoomReq, rsp *pb.RummyRemoveRoomRsp) error {
	m.datas[req.RoomId] = nil
	err := m.SendMsg(&pb.Head{FuncName: "Save"})
	if err != nil {
		mlog.Infof("MatchRummyRoomMgr Remove err:%v ", err)
	}

	rsp.RoomId = req.RoomId
	return nil
}
