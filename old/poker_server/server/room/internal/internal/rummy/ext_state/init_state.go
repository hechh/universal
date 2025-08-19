package ext_state

import (
	"poker_server/common/pb"
	"poker_server/server/room/internal/internal/rummy"
)

type RummyInitstate struct {
	BaseState
}

func (d *RummyInitstate) OnEnter(nowMs int64, curState pb.GameState, extra interface{}) {
	game := extra.(*rummy.RummyGame)

	// 重置房间状态
	game.Data.Stage = curState

	// init 初始化游戏
	if game.Data.Common.GameFinish {
		game.ExtReset() //重置房间
		game.UpReadyPlayer()
		game.Data.Status = pb.RoomStatus_RoomStatusWait
	}

	game.DelExitPlayer(nowMs)
}

func (d *RummyInitstate) OnTick(nowMs int64, curState pb.GameState, extra interface{}) pb.GameState {
	game := extra.(*rummy.RummyGame)
	// 回收退出游戏用户
	game.DelExitPlayer(nowMs)

	if game.IsFinish {
		return curState
	}

	// 判断人数进入准备阶段
	if game.GetPlayerCount() >= game.GetMinStartPlayers() {
		return game.GetNextState()
	} else {
		if game.Data.IsGroup {
			return pb.GameState_RummyExt_STAGE_FIN_SETTLE
		} else {
			return curState
		}
	}
}

func (d *RummyInitstate) OnExit(nowMs int64, curState pb.GameState, extra interface{}) {}
