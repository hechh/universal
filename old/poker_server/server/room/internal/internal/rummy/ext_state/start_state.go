package ext_state

import (
	"poker_server/common/pb"
	"poker_server/library/mlog"
	"poker_server/server/room/internal/internal/rummy"
)

type RummyStartState struct {
	BaseState
}

func (d *RummyStartState) OnEnter(nowMs int64, curState pb.GameState, extra interface{}) {
	d.Log(curState)
	game := extra.(*rummy.RummyGame)
	game.Data.Match.StartTime = nowMs
	game.Data.Status = pb.RoomStatus_RoomStatusPlaying

	// ext 玩法 连续比赛
	game.Data.IsGroup = true

	game.Walk(0, func(player *pb.RummyRoomPlayer) bool {
		if game.Data.Match.Match == 1 {
			// 玩家支付入场费
			game.CostTicket(player)
		}
		return true
	})

	ntf := &pb.RummyStartGameNtf{
		RoomId:    game.GetRoomId(),
		ZhuangId:  game.Data.Common.ZhuangId,
		CurMatch:  game.Data.Match.Match,
		PlayerIds: game.Data.Common.PlayerIds,
	}
	err := game.NotifyToClient(game.GetPlayerUidList(), pb.RummyEventType_RummyStartGame, ntf)
	mlog.Infof("游戏开始广播返回:%v ", err)
}

func (d *RummyStartState) OnExit(nowMs int64, curState pb.GameState, extra interface{}) {
}
func (d *RummyStartState) OnTick(nowMs int64, curState pb.GameState, extra interface{}) pb.GameState {
	return d.MoveStateByTime(nowMs, curState, extra)
}
