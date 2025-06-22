package state

import (
	"poker_server/common/pb"
	"poker_server/library/mlog"
	"poker_server/server/room/internal/internal/rummy"
)

/*
	FUNC: 胜利者胡牌后 最后时间给所有人调整手牌
*/

type RummyFixState struct {
	BaseState
}

func (d *RummyFixState) OnEnter(nowMs int64, curState pb.GameState, extra interface{}) {
	d.Log(curState)
	game := extra.(*rummy.RummyGame)
	game.FlushExpireTime(nowMs)
	game.Change()

	// 通知玩家最后调整手牌
	ntf := &pb.RummyFixCardPlayersNtf{
		RoomId:      game.GetRoomId(),
		Players:     game.Data.Common.PlayerIds,
		WinId:       game.Data.Common.WinnerId,
		TimeOut:     game.Data.Common.TimeOut,
		CurPlayerId: game.Data.Common.WinnerId,
		TotalTime:   game.Data.Common.TotalTime,
	}
	err := game.NotifyToClient(game.GetPlayerUidList(), pb.RummyEventType_RummyFixCardPlayers, ntf)
	mlog.Infof("RummyFixState ntf all player RummyFixCardPlayersNtf err : %v", err)
}

func (d *RummyFixState) OnExit(nowMs int64, curState pb.GameState, extra interface{}) {
	// 回合结束未确认玩家 自动确认
	game := extra.(*rummy.RummyGame)
	game.Walk(0, func(player *pb.RummyRoomPlayer) bool {
		if player.State == pb.RummyPlayState_Rummy_PLAYSTATE_PLAY {
			game.SetPlayerLose(player.PlayerId, false)
		}
		return true
	})
}

func (d *RummyFixState) OnTick(nowMs int64, curState pb.GameState, extra interface{}) pb.GameState {

	game := extra.(*rummy.RummyGame)

	fastFinish := true
	//所有玩家确认完毕 快速结束确认回合
	game.Walk(0, func(player *pb.RummyRoomPlayer) bool {
		if player.State == pb.RummyPlayState_Rummy_PLAYSTATE_PLAY {
			fastFinish = false
			return false
		}
		return true
	})
	if fastFinish {
		return pb.GameState_Rummy_STAGE_SETTLE
	}

	return d.MoveStateByTime(nowMs, curState, extra)
}
