package state

import (
	"poker_server/common/pb"
	"poker_server/server/room/internal/internal/rummy"
)

/*
	FUNC: 初始状态:负责牌局初始化和等待玩家准备。
*/

type RummyInitstate struct {
	BaseState
}

func (d *RummyInitstate) OnEnter(nowMs int64, curState pb.GameState, extra interface{}) {
	game := extra.(*rummy.RummyGame)

	// 重置房间状态
	game.Data.Stage = curState

	// init 初始化游戏
	if game.Data.Common.GameFinish {
		game.Reset() //重置房间
		game.UpReadyPlayer()
	}
}

func (d *RummyInitstate) OnTick(nowMs int64, curState pb.GameState, extra interface{}) pb.GameState {
	game := extra.(*rummy.RummyGame)
	// 回收退出游戏用户
	game.DelExitPlayer(0, nowMs)

	// 判断人数进入准备阶段
	if game.GetPlayerCount() >= game.GetMinStartPlayers() {
		return game.GetNextState()
	}
	return curState
}

func (d *RummyInitstate) OnExit(nowMs int64, curState pb.GameState, extra interface{}) {

}
