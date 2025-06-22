package rummy_room

import (
	"poker_server/common/config/repository/rummy_config"
	"poker_server/common/dao/repository/redis/rummy_room_data"
	"poker_server/common/pb"
	"poker_server/framework/actor"
	"poker_server/library/mlog"
	"reflect"
	"time"
)

type RummyRoomData struct {
	*pb.RummyRoomData
	isDelete bool
	isChange bool
}

type DbRummyRoomMgr struct {
	actor.Actor
	datas map[uint64]*RummyRoomData // 房间数据
}

func NewDbRummyRoomMgr() *DbRummyRoomMgr {
	ret := &DbRummyRoomMgr{}
	ret.Actor.Register(ret)
	ret.Actor.ParseFunc(reflect.TypeOf(ret))
	ret.Actor.SetId(uint64(pb.DataType_DataTypeRummyRoom))
	return ret
}

func (m *DbRummyRoomMgr) Init() error {
	rooms, err := rummy_room_data.HGetAll()
	if err != nil {
		return err
	}

	// 初始化所有数据
	m.datas = make(map[uint64]*RummyRoomData)
	for _, room := range rooms {
		cfg := rummy_config.MGetID(room.GameId)
		if cfg == nil {
			m.datas[room.RoomId] = &RummyRoomData{RummyRoomData: room, isDelete: true}
		} else {
			m.datas[room.RoomId] = &RummyRoomData{RummyRoomData: room}
		}
	}

	m.Start()
	actor.Register(m)

	// 注册定时器
	return m.RegisterTimer(&pb.Head{
		SendType:  pb.SendType_POINT,
		ActorName: "DbRummyRoomMgr",
		FuncName:  "OnTick",
	}, 5*time.Second, -1)
}

func (m *DbRummyRoomMgr) OnTick() {
	isDelete := false
	isChange := false
	for _, item := range m.datas {
		cfg := rummy_config.MGetID(item.GameId)
		if cfg == nil {
			item.isDelete = true
			isDelete = true
			continue
		}
		isDelete = isDelete || item.isDelete
		isChange = isChange || item.isChange
	}
	if isDelete {
		m.SendMsg(&pb.Head{FuncName: "Delete"})
	}
	if isChange {
		m.SendMsg(&pb.Head{FuncName: "Save"})
	}
}

// 保存数据
func (m *DbRummyRoomMgr) Save() error {
	saves := map[string]*pb.RummyRoomData{}
	ids := []uint64{}
	for _, item := range m.datas {
		if item.isChange {
			ids = append(ids, item.RoomId)
			saves[rummy_room_data.GetField(item.RoomId)] = item.RummyRoomData
		}
	}

	// 保存数据
	mlog.Infof("DbRummyRoomMgr Save rooms data")
	if err := rummy_room_data.HMSet(saves); err != nil {
		return err
	}

	// 清除修改标记
	for _, id := range ids {
		m.datas[id].isChange = false
	}
	return nil
}

// 删除数据
func (m *DbRummyRoomMgr) Delete() error {
	deletes := []string{}
	ids := []uint64{}
	for _, item := range m.datas {
		if item.isDelete {
			deletes = append(deletes, rummy_room_data.GetField(item.RoomId))
			ids = append(ids, item.RoomId)
		}
	}

	// 删除数据
	if err := rummy_room_data.HDel(deletes...); err != nil {
		return err
	}

	// 清除删除标记
	for _, id := range ids {
		delete(m.datas, id)
	}
	return nil
}

// 更新数据
func (m *DbRummyRoomMgr) Update(head *pb.Head, req *pb.UpdateRummyRoomDataNotify) error {
	if len(req.List) <= 0 {
		return nil
	}

	for _, item := range req.List {
		m.datas[item.RoomId] = &RummyRoomData{
			RummyRoomData: item,
			isDelete:      false,
			isChange:      true,
		}
	}
	return nil
}

// 加载数据请求(同步和异步请求支持)
func (m *DbRummyRoomMgr) Query(head *pb.Head, req *pb.GetRummyRoomDataReq, rsp *pb.GetRummyRoomDataRsp) error {
	for _, item := range m.datas {
		rsp.List = append(rsp.List, item.RummyRoomData)
	}
	return nil
}
