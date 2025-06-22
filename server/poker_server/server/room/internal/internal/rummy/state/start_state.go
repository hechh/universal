package state

import (
	"poker_server/common/pb"
	"poker_server/library/mlog"
	"poker_server/server/room/internal/internal/rummy"
)

/*
	FUNC: 开始游戏状态，进入游戏状态。
*/

type RummyStartState struct {
	BaseState
}

func (d *RummyStartState) OnEnter(nowMs int64, curState pb.GameState, extra interface{}) {
	d.Log(curState)
	game := extra.(*rummy.RummyGame)
	game.Data.Match.StartTime = nowMs

	ntf := &pb.RummyStartGameNtf{
		RoomId:   game.GetRoomId(),
		ZhuangId: game.Data.Common.ZhuangId,
		CurMatch: game.Data.Match.Match,
	}
	err := game.NotifyToClient(game.GetPlayerUidList(), pb.RummyEventType_RummyStartGame, ntf)
	mlog.Infof("游戏开始广播返回:%v ", err)
}

func (d *RummyStartState) OnExit(nowMs int64, curState pb.GameState, extra interface{}) {
}
func (d *RummyStartState) OnTick(nowMs int64, curState pb.GameState, extra interface{}) pb.GameState {
	return d.MoveStateByTime(nowMs, curState, extra)
}
