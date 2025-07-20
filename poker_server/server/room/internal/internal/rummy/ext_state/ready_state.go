package ext_state

import (
	"poker_server/common/pb"
	"poker_server/library/mlog"
	"poker_server/server/room/internal/internal/rummy"
)

type RummyReadyState struct {
	BaseState
}

func (d *RummyReadyState) OnEnter(nowMs int64, curState pb.GameState, extra interface{}) {
	game := extra.(*rummy.RummyGame)
	game.Change()
	mlog.Infof("%v游戏人数已就位 准备开始倒计时 %v", game.Data.RoomCfg.GameType, game.GetPlayerMap())
	game.FlushExpireTime(nowMs)

	ntf := &pb.RummyReadyStartGameNtf{
		RoomId:    game.GetRoomId(),
		TimeOut:   game.Data.Common.TimeOut,
		TotalTime: game.Data.Common.TotalTime,
	}

	err := game.NotifyToClient(game.GetPlayerUidList(), pb.RummyEventType_RummyReadyStartGame, ntf)
	mlog.Infof("游戏人数已就位 系统通知下发:%v", err)
}

func (d *RummyReadyState) OnExit(nowMs int64, curState pb.GameState, extra interface{}) {}

func (d *RummyReadyState) OnTick(nowMs int64, curState pb.GameState, extra interface{}) pb.GameState {
	game := extra.(*rummy.RummyGame)
	game.DelExitPlayer(nowMs)
	// 判断人数进入准备阶段
	if game.GetPlayerCount() < game.GetMinStartPlayers() {
		if game.Data.IsGroup {
			// 人数不够开启游戏 结束游戏
			return pb.GameState_RummyExt_STAGE_FIN_SETTLE
		} else {
			// 人数不够开启游戏 返回最初阶段
			return pb.GameState_RummyExt_STAGE_INIT
		}
	}

	return d.MoveStateByTime(nowMs, curState, extra)
}
