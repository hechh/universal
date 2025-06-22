package texas_room

import (
	"poker_server/common/config/repository/texas_config"
	"poker_server/common/dao/repository/redis/texas_room_data"
	"poker_server/common/pb"
	"poker_server/framework/actor"
	"poker_server/library/mlog"
	"reflect"
	"time"
)

type TexasRoomData struct {
	*pb.TexasRoomData
	isDelete bool
	isChange bool
}

type DbTexasRoomMgr struct {
	actor.Actor
	datas map[uint64]*TexasRoomData
}

func NewTexasRoomMgr() *DbTexasRoomMgr {
	ret := &DbTexasRoomMgr{}
	ret.Actor.Register(ret)
	ret.Actor.ParseFunc(reflect.TypeOf(ret))
	ret.Actor.SetId(uint64(pb.DataType_DataTypeTexasRoom))
	return ret
}

func (m *DbTexasRoomMgr) Init() error {
	rooms, err := texas_room_data.HGetAll()
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	m.datas = make(map[uint64]*TexasRoomData)
	for _, room := range rooms {
		cfg := texas_config.MGetID(room.GameId)
		if cfg == nil || room.CreateTime+cfg.RoomKeepLive*60 <= now {
			m.datas[room.RoomId] = &TexasRoomData{TexasRoomData: room, isDelete: true}
		} else {
			m.datas[room.RoomId] = &TexasRoomData{TexasRoomData: room}
		}
	}

	m.Start()
	actor.Register(m)
	m.RegisterTimer(&pb.Head{FuncName: "OnTick"}, 5*time.Second, -1)
	return nil
}

func (m *DbTexasRoomMgr) Close() {
	m.Save()
	m.Delete()
	m.Actor.Stop()
}

func (m *DbTexasRoomMgr) OnTick() {
	now := time.Now().Unix()
	isDelete := false
	isChange := false
	for _, item := range m.datas {
		cfg := texas_config.MGetID(item.GameId)
		if cfg == nil || item.CreateTime+cfg.RoomKeepLive*60 <= now {
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

func (m *DbTexasRoomMgr) Save() error {
	saves := map[string]*pb.TexasRoomData{}
	ids := []uint64{}
	for _, item := range m.datas {
		if item.isChange {
			ids = append(ids, item.RoomId)
			saves[texas_room_data.GetField(item.RoomId)] = item.TexasRoomData
		}
	}

	mlog.Infof("DbTexasRoomMgr Save rooms data")
	if err := texas_room_data.HMSet(saves); err != nil {
		return err
	}

	for _, id := range ids {
		m.datas[id].isChange = false
	}
	return nil
}

func (m *DbTexasRoomMgr) Delete() error {
	deletes := []string{}
	ids := []uint64{}
	for _, item := range m.datas {
		if item.isDelete {
			deletes = append(deletes, texas_room_data.GetField(item.RoomId))
			ids = append(ids, item.RoomId)
		}
	}
	if err := texas_room_data.HDel(deletes...); err != nil {
		return err
	}
	for _, id := range ids {
		delete(m.datas, id)
	}
	return nil
}

func (m *DbTexasRoomMgr) Update(head *pb.Head, req *pb.UpdateTexasRoomDataNotify) error {
	if len(req.List) <= 0 {
		return nil
	}
	for _, item := range req.List {
		m.datas[item.RoomId] = &TexasRoomData{
			TexasRoomData: item,
			isDelete:      false,
			isChange:      true,
		}
	}
	return nil
}

func (m *DbTexasRoomMgr) Query(head *pb.Head, req *pb.GetTexasRoomDataReq, rsp *pb.GetTexasRoomDataRsp) error {
	for _, item := range m.datas {
		rsp.List = append(rsp.List, item.TexasRoomData)
	}
	return nil
}
