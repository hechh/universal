package state

import (
	"poker_server/common/pb"
	"poker_server/library/mlog"
	"poker_server/server/room/internal/internal/rummy"
	"time"
)

/*
	FUNC: 打牌阶段:摸牌，出牌，结束牌，弃牌等。
*/

type RummyPlayState struct {
	BaseState
}

func (d *RummyPlayState) OnTick(nowMs int64, curState pb.GameState, extra interface{}) (ret pb.GameState) {
	game := extra.(*rummy.RummyGame)
	ret = curState

	if game.Data.Common.GameFinish {
		return game.GetNextState()
	}

	if game.Data.Common.TimeOut <= nowMs {
		mlog.Infof("room.ExpirTime %v", time.UnixMilli(game.Data.Common.TimeOut).Format("2006-01-02 15:04:05.000"))
		game.OnPlayerTimeout()
		game.FlushExpireTime(nowMs)
		ret = game.SetNextPlayer()
	} else {
		// 当前玩家结束 正常过渡下个玩家
		if game.Data.Common.OprPlayer.Step == pb.RummyRoundStep_PLAYSTEPSTEAFINISH {
			game.FlushExpireTime(nowMs)
			ret = game.SetNextPlayer()
		}
	}
	//超时检查 过渡回合执行权到下个玩家
	return
}

func (d *RummyPlayState) OnEnter(nowMs int64, curState pb.GameState, extra interface{}) {
	d.Log(curState)
	game := extra.(*rummy.RummyGame)
	game.Change()
	game.FlushExpireTime(nowMs)
	game.SetNextPlayer()
}

func (d *RummyPlayState) OnExit(nowMs int64, curState pb.GameState, extra interface{}) {}
