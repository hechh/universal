package rummy

import (
	"poker_server/common/pb"
	"poker_server/library/mlog"
	"poker_server/library/uerror"
	"time"
)

// extra 玩法不需要readyreq 通过matchsvr game room直接创建房间
// destory ready

// RummyGiveUpReq pool等分支玩法 退出Rummy游戏
func (d *RummyGame) RummyGiveUpReq(head *pb.Head, req *pb.RummyGiveUpReq, rsp *pb.RummyGiveUpRsp) error {
	player, ok := d.Data.Common.Players[head.Uid]
	if !ok || player == nil {
		return uerror.New(1, pb.ErrorCode_RUMMY_PLAYER_NOT_IN_ROOM, "玩家不在房间内: %d", head.Uid)
	}

	d.giveUpGame(head.Uid)

	player.Health = pb.RummyPlayHealth_Rummy_QUIT
	player.State = pb.RummyPlayState_Rummy_PLAYSTATE_ELIMINATED //淘汰状态

	//推送退出消息
	ntf := &pb.RummyQuitRoomNtf{
		RoomId:    d.GetRoomId(),
		PlayerId:  head.Uid,
		LeaveType: pb.RummyLeaveType_Rummy_LEAVE_TYPE_QUIT,
	}
	err := d.NotifyToClient(d.GetPlayerUidList(), pb.RummyEventType_RummyQuitRoom, ntf)
	mlog.Infof("RummyGiveUpRsp QuitRoomReq ntf send: %v", err)

	rsp.RoomId = d.GetRoomId()
	rsp.Charge = int64(d.Data.Common.Players[head.Uid].Coin)
	d.Data.Common.Players[head.Uid].Coin = 0
	d.Change()
	return nil
}

// RummyTryHalfPoolReq Pool玩法玩家平分奖池
func (d *RummyGame) RummyTryHalfPoolReq(head *pb.Head, req *pb.RummyTryHalfPoolReq, rsp *pb.RummyTryHalfPoolRsp) error {
	if !(d.Data.HalfState == pb.HalfState_Allow || d.Data.HalfState == pb.HalfState_Try) {
		return uerror.New(1, pb.ErrorCode_RUMMY_PLAYER_DATA_INVAILD, "不满足平分奖池条件")
	}

	if d.Data.HalfState == pb.HalfState_Allow {
		d.Data.HalfState = pb.HalfState_Try
		duration := int64(10000)
		if d.Data.MachineCfg.HalfDuration != 0 {
			duration = d.Data.MachineCfg.HalfDuration * 1000
		}

		tmp := time.Now().UnixMilli() + duration
		if d.Data.Common.TimeOut < tmp {
			d.Data.Common.TimeOut = tmp
			d.Data.Common.TotalTime = duration
		}
	}
	d.Data.Common.Players[head.Uid].State = pb.RummyPlayState_Rummy_PLAYSTATE_Half

	count := 0
	d.Walk(0, func(player *pb.RummyRoomPlayer) bool {
		if player.State != pb.RummyPlayState_Rummy_PLAYSTATE_Half && player.State != pb.RummyPlayState_Rummy_PLAYSTATE_ELIMINATED {
			count++ // 统计所有玩家
		}
		return true
	})
	if count == 0 {
		d.Data.HalfState = pb.HalfState_Sure
	}

	ntf := &pb.RummyTryHalfPoolNtf{
		RoomId:  d.GetRoomId(),
		Players: d.Data.Common.Players,
		TimeOut: d.Data.Common.TimeOut,
	}

	err := d.NotifyToClient(d.GetPlayerUidList(), pb.RummyEventType_RummyFixCardPlayers, ntf)
	mlog.Infof("RummyTryHalfPoolReq NotifyToClient %v", err)
	return nil
}
