package manager

import (
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/actor"
	"poker_server/library/uerror"
	"poker_server/server/room/internal/internal/rummy"
	"reflect"
)

type RummyGameMgr struct {
	actor.Actor
	mgr *actor.ActorMgr
}

func NewRummyGameMgr() *RummyGameMgr {
	// 预先注册
	mgr := new(actor.ActorMgr)
	room := &rummy.RummyGame{}
	mgr.Register(room)
	mgr.ParseFunc(reflect.TypeOf(room))
	actor.Register(mgr)

	// 创建房间actor 管理器
	ret := &RummyGameMgr{mgr: mgr}
	ret.Actor.Register(ret)
	ret.Actor.ParseFunc(reflect.TypeOf(ret))
	ret.Actor.SetId(uint64(pb.DataType_DataTypeRummyRoom))
	ret.Actor.Start()
	actor.Register(ret)
	return ret
}

func (d *RummyGameMgr) Stop() {
	// 停止所有游戏
	d.mgr.Stop()

	// 停止自己
	d.Actor.Stop()
}

// 加入房间 玩家历史房间优先 房间不存在则新建房间
func (d *RummyGameMgr) JoinRoomReq(head *pb.Head, req *pb.RummyJoinRoomReq, rsp *pb.RummyJoinRoomRsp) error {
	if act := d.mgr.GetActor(req.RoomId); act != nil {
		return act.SendMsg(head, req, rsp)
	}

	gameType := uint64(req.RoomId>>32) & 0xFF
	coinType := uint64(req.RoomId>>24) & 0xFF
	// 从match服务请求数据
	dst := framework.NewMatchRouter(gameType<<32|coinType, "MatchRummyRoom", "Query")
	newHead := framework.NewHead(dst, pb.RouterType_RouterTypeDataType, uint64(pb.DataType_DataTypeRummyRoom))
	data := &pb.RummyRoomData{}
	if err := framework.Request(newHead, req.RoomId, data); err != nil {
		return err
	}

	// 创建房间
	rr := rummy.NewRummyGame(data)
	if rr == nil {
		return uerror.NEW(pb.ErrorCode_CONFIG_NOT_FOUND, head, "rummy配置不存在: %d", data.GameId)
	}

	d.mgr.AddActor(rr)
	return rr.SendMsg(head, req, rsp)
}
