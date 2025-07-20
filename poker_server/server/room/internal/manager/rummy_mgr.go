package manager

import (
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/actor"
	"poker_server/framework/cluster"
	"poker_server/library/mlog"
	"poker_server/library/safe"
	"poker_server/library/uerror"
	"poker_server/library/util"
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
	mlog.Infof("RummyGameMgr关闭成功")
}

// 加入房间 玩家历史房间优先 房间不存在则新建房间
func (d *RummyGameMgr) JoinRoomReq(head *pb.Head, req *pb.RummyJoinRoomReq, rsp *pb.RummyJoinRoomRsp) error {
	if act := d.mgr.GetActor(req.RoomId); act != nil {
		return act.SendMsg(head, req, rsp)
	}

	types := util.DestructRoomId(req.RoomId)
	// 从match服务请求数据
	newHead := &pb.Head{
		Src: framework.NewSrcRouter(d.GetId(), d.GetActorName()),
		Dst: framework.NewMatchRouter(uint64(types.GetGameType())<<32|uint64(types.GetCoinType()), "MatchRummyRoom", "Query"),
	}
	data := &pb.RummyRoomData{}
	if err := cluster.Request(newHead, req.RoomId, data); err != nil {
		return err
	}

	// 创建房间
	rr := rummy.NewRummyGame(data)
	if rr == nil {
		return uerror.New(1, pb.ErrorCode_CONFIG_NOT_FOUND, "rummy配置不存在: %d", data.GameId)
	}

	d.mgr.AddActor(rr)
	return rr.SendMsg(head, req, rsp)
}

func (d *RummyGameMgr) Remove(roomID uint64) {
	if act := d.mgr.GetActor(roomID); act != nil {
		types := util.DestructRoomId(roomID)

		head := &pb.Head{
			Src: framework.NewSrcRouter(roomID, "RummyGameMgr", "RemoveCompleted"),
			Dst: framework.NewMatchRouter(uint64(types.GetGameType())<<32|uint64(types.GetCoinType()), "MatchRummyRoomMgr", "Remove"),
		}
		req := &pb.RummyRemoveRoomReq{RoomId: roomID}
		err := cluster.Send(head, req)
		mlog.Infof("RummyGameMgr Remove Actor framework SendMsg err:%v ", err)
	}
}

func (d *RummyGameMgr) RemoveCompleted(head *pb.Head, rsp *pb.RummyRemoveRoomRsp) error {
	if act := d.mgr.GetActor(rsp.RoomId); act != nil {
		d.mgr.DelActor(rsp.RoomId)

		// 等待消息处理完成，然后关闭连接
		d.Add(1)
		safe.Go(func() {
			act.Stop()
			d.Done()
		})
	}
	return nil
}
