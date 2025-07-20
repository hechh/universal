package room_info

import (
	"poker_server/common/pb"
	"poker_server/common/redis/repository/room_info"
	"poker_server/common/room_util"
	"poker_server/framework/actor"
	"poker_server/library/mlog"
	"reflect"
	"time"
)

type RoomInfo struct {
	*pb.RoomInfo
	updateTime int64
}

type RoomInfoMgr struct {
	actor.Actor
	datas map[uint64]*RoomInfo
	uids  []uint64
}

func NewRoomInfoMgr() *RoomInfoMgr {
	ret := &RoomInfoMgr{datas: make(map[uint64]*RoomInfo)}
	ret.Actor.Register(ret)
	ret.Actor.ParseFunc(reflect.TypeOf(ret))
	ret.Actor.SetId(uint64(pb.DataType_DataTypeRoomInfo))
	ret.Actor.Start()
	actor.Register(ret)
	ret.RegisterTimer(&pb.Head{FuncName: "OnTick"}, 5*time.Second, -1)
	return ret
}

func (d *RoomInfoMgr) Init() error {
	return nil
}

func (d *RoomInfoMgr) Close() {
	d.Save()
	d.Stop()
	mlog.Infof("RoomInfoMgr关闭成功")
}

func (d *RoomInfoMgr) OnTick() {
	d.Save()
	// 清楚过期数据
	now := time.Now().Unix()
	for key, item := range d.datas {
		if now-item.updateTime < 30*60 {
			continue
		}
		delete(d.datas, key)
	}
}

func (d *RoomInfoMgr) Save() {
	if len(d.uids) <= 0 {
		return
	}
	tmps := map[string]*pb.RoomInfo{}
	for _, uid := range d.uids {
		tmps[room_info.GetField(uid)] = d.datas[uid].RoomInfo
	}
	if err := room_info.HMSet(tmps); err != nil {
		mlog.Errorf("RoomInfo数据保存失败%v", err)
	} else {
		d.uids = d.uids[:0]
	}
}

func (m *RoomInfoMgr) Update(head *pb.Head, req *pb.UpdateRoomInfoNotify) error {
	if req.Data == nil {
		return nil
	}
	if vals, ok := m.datas[req.Data.Uid]; ok {
		if vals.GameType == req.Data.GameType && vals.RoomId == req.Data.RoomId && vals.TableId == req.Data.TableId {
			return nil
		}
		vals.RoomInfo = req.Data
		vals.updateTime = time.Now().Unix()
	} else {
		m.datas[req.Data.Uid] = &RoomInfo{RoomInfo: req.Data, updateTime: time.Now().Unix()}
	}
	for _, uid := range m.uids {
		if uid == req.Data.Uid {
			return nil
		}
	}
	m.uids = append(m.uids, req.Data.Uid)
	return nil
}

func (m *RoomInfoMgr) Query(head *pb.Head, req *pb.QueryPlayerDataReq, rsp *pb.QueryPlayerDataRsp) error {
	item, ok := m.datas[head.Uid]
	if !ok || item.RoomId <= 0 {
		data, ok, err := room_info.HGet(head.Uid)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		m.datas[data.Uid] = &RoomInfo{updateTime: time.Now().Unix(), RoomInfo: data}
		item = m.datas[data.Uid]
	}

	matchType, gameType, coinType := room_util.TexasRoomIdTo(item.RoomId)
	rsp.Data = &pb.PlayerDataBase{
		RoomInfo: &pb.PlayerRoomInfo{
			GameType: item.GameType,
			RoomId:   item.RoomId,
			TableId:  item.TableId,
		},
	}
	rsp.MatchType = matchType
	rsp.GameType = gameType
	rsp.CoinType = coinType
	return nil
}
