package texas

import (
	"poker_server/common/dao/repository/redis/texas_room_player"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/library/async"
	"poker_server/framework/library/mlog"
	"poker_server/server/room/texas/internal/base"
	"poker_server/server/room/texas/internal/texas_room"
	"reflect"
	"time"
)

type TexasRoomMgr struct {
	framework.Actor
	mgr *framework.ActorMgr
}

func NewTexasRoomMgr() *TexasRoomMgr {
	// 预先注册
	mgr := new(framework.ActorMgr)
	room := &texas_room.TexasRoom{}
	mgr.Register(room)
	mgr.ParseFunc(reflect.TypeOf(room))
	framework.RegisterActor(mgr)

	// 创建房间管理器
	ret := &TexasRoomMgr{mgr: mgr}
	ret.Actor.Register(ret)
	ret.Actor.ParseFunc(reflect.TypeOf(ret))
	ret.Actor.Start()
	framework.RegisterActor(ret)

	ret.loop()
	return ret
}

func (d *TexasRoomMgr) loop() {
	async.SafeGo(mlog.Fatalf, func() {
		tt := time.NewTicker(100 * time.Millisecond)
		defer tt.Stop()

		for {
			<-tt.C
			now := time.Now().UnixMilli()
			head := &pb.Head{ActorName: "TexasRoom", FuncName: "OnTick", SendType: pb.SendType_BROADCAST}
			d.mgr.SendMsg(head, now)
		}
	})
}

// 删除房间
func (d *TexasRoomMgr) RemoveRoom(roomId uint64) {
	if act := d.mgr.GetActor(roomId); act != nil {
		act.Stop()
		d.mgr.DelActor(roomId)
	}
}

// 加入房间
func (d *TexasRoomMgr) JoinRoomReq(uid uint64, req *pb.TexasJoinRoomReq, rsp *pb.TexasJoinRoomRsp) {
	// 是否已经加入房间
	roomId, err := texas_room_player.HGet(uid)
	if err != nil {
		rsp.Head = base.ToRet(pb.ErrorCode_EC_TEXAS_REDIS_LOAD_ROOM, err.Error())
		texas_room.SendPackageByPlayerIDRummy(uid, pb.CMD_TEXAS_JOIN_ROOM_RSP, rsp)
		mlog.Errorf("redis加载房间失败, uid:%d, roomId:%d, err:%s", uid, req.RoomId, err.Error())
		return
	}
	if roomId > 0 {
		req.RoomId = roomId
	}

	if act := d.mgr.GetActor(req.RoomId); act != nil {
		head := &pb.Head{ActorName: "TexasRoom", FuncName: "JoinRoomReq", Id: req.RoomId}
		act.SendMsg(head, uid, req, rsp)
		return
	}

	// 加载房间
	newRoom := texas_room.NewTexasRoom()
	newRoom.Start()
	d.mgr.AddActor(req.RoomId, newRoom)

	// 加入房间
	head := &pb.Head{ActorName: "TexasRoom", FuncName: "JoinRoomReq", Id: req.RoomId}
	framework.SendMsg(head, uid, req, rsp)
}
