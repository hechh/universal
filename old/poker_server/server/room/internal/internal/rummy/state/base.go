package state

import (
	"poker_server/common/pb"
	"poker_server/library/mlog"
	"poker_server/server/room/internal/internal/rummy"
)

type BaseState struct{}

func (s *BaseState) Log(curState pb.GameState) {
	mlog.Infof("当前状态机状态:%v", curState)
}

// MoveStateByTime 倒计时结束自动切换阶段
func (s *BaseState) MoveStateByTime(nowMs int64, curState pb.GameState, extra interface{}) pb.GameState {

	game := extra.(*rummy.RummyGame)
	if game.Data.Common.TimeOut <= nowMs {
		// 开始游戏
		return game.GetNextState()
	}
	return curState
}
