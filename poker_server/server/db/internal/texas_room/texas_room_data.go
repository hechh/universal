package texas_room

import (
	"poker_server/common/config/repository/texas_config"
	"poker_server/common/pb"
	"poker_server/common/redis/repository/texas_room_data"
	"poker_server/framework/actor"
	"poker_server/library/mlog"
	"reflect"
	"time"
)

type DbTexasRoomMgr struct {
	actor.Actor
	datas map[uint64]*pb.TexasRoomData
	list  []uint64
	dels  []uint64
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
	m.datas = make(map[uint64]*pb.TexasRoomData)
	dels := []string{}
	for _, room := range rooms {
		cfg := texas_config.MGetID(room.GameId)
		if cfg == nil || room.CreateTime+cfg.RoomKeepLive*60 <= now {
			dels = append(dels, texas_room_data.GetField(room.RoomId))
		} else {
			m.datas[room.RoomId] = room
		}
	}

	if err := texas_room_data.HDel(dels...); err != nil {
		return err
	}

	m.Start()
	actor.Register(m)
	m.RegisterTimer(&pb.Head{FuncName: "OnTick"}, 5*time.Second, -1)
	return nil
}

func (m *DbTexasRoomMgr) Close() {
	m.Actor.Stop()
	mlog.Infof("DbTexasRoomMgr关闭成功")
}

func (m *DbTexasRoomMgr) OnTick() {
	// 删除
	tmps := []string{}
	for _, roomId := range m.dels {
		delete(m.datas, roomId)
		tmps = append(tmps, texas_room_data.GetField(roomId))
	}
	texas_room_data.HDel(tmps...)
	m.dels = m.dels[:0]

	// 保存
	saves := map[string]*pb.TexasRoomData{}
	for _, roomId := range m.list {
		item, ok := m.datas[roomId]
		if !ok {
			continue
		}
		saves[texas_room_data.GetField(item.RoomId)] = item
	}
	m.list = m.list[:0]
	texas_room_data.HMSet(saves)
}

func (m *DbTexasRoomMgr) Update(head *pb.Head, req *pb.UpdateTexasRoomDataNotify) error {
	if req.Data == nil {
		return nil
	}
	m.datas[req.Data.RoomId] = req.Data
	m.list = append(m.list, req.Data.RoomId)
	return nil
}

func (m *DbTexasRoomMgr) Delete(head *pb.Head, req *pb.DeleteTexasRoomDataNotify) error {
	m.dels = append(m.dels, req.RoomId)
	return nil
}

func (m *DbTexasRoomMgr) Query(head *pb.Head, req *pb.GetTexasRoomDataReq, rsp *pb.GetTexasRoomDataRsp) error {
	for _, item := range m.datas {
		rsp.List = append(rsp.List, item)
	}
	return nil
}
