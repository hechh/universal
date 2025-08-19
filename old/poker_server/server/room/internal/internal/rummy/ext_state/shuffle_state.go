package ext_state

import (
	"poker_server/common/pb"
	"poker_server/library/mlog"
	"poker_server/server/room/internal/internal/rummy"
)

/*
	FUNC: 玩牌阶段 牌堆摸空触发洗牌后 返回玩牌阶段。
*/

type RummyShuffleState struct {
	BaseState
}

func (d *RummyShuffleState) OnEnter(nowMs int64, curState pb.GameState, extra interface{}) {
	d.Log(curState)
	game := extra.(*rummy.RummyGame)
	game.FlushExpireTime(nowMs)

	game.Data.Private.Cards = game.Data.Common.OprPlayer.OutCards[:len(game.Data.Common.OprPlayer.OutCards)-1]
	game.Data.Common.OprPlayer.OutCards = game.Data.Common.OprPlayer.OutCards[len(game.Data.Common.OprPlayer.OutCards)-1:]
	game.Data.Private.CardIdx = 0

	ntf := &pb.RummyShuffleNtf{
		RoomId:    game.GetRoomId(),
		DrawCount: uint32(len(game.Data.Private.Cards)),
		ShowCard:  game.Data.Common.OprPlayer.OutCards[0],
	}
	err := game.NotifyToClient(game.GetPlayerUidList(), pb.RummyEventType_RummyShuffle, ntf)
	mlog.Debugf("RummyShuffleState NotifyToClient:%v", err)
}

func (d *RummyShuffleState) OnExit(nowMs int64, curState pb.GameState, extra interface{}) {}

func (d *RummyShuffleState) OnTick(nowMs int64, curState pb.GameState, extra interface{}) pb.GameState {
	return d.MoveStateByTime(nowMs, curState, extra)
}
