package ext_state

import (
	"math"
	"poker_server/common/pb"
	"poker_server/library/mlog"
	"poker_server/server/room/internal/internal/rummy"
)

/*
FUNC: 多局游戏 单局结算状态:实现结算和结算展示。
*/
type RummySettleState struct {
	BaseState
}

func (d *RummySettleState) OnEnter(nowMs int64, curState pb.GameState, extra interface{}) {
	d.Log(curState)
	game := extra.(*rummy.RummyGame)
	game.FlushExpireTime(nowMs)
	game.Change()

	scoreTmp := int64(0)
	record := &pb.RummyMatchInfo{
		Players: make([]*pb.RummySettlePlayerInfo, 0, game.GetPlayerCount()),
	}
	// 更新玩家total和淘汰状态
	winner := game.Data.Common.Players[game.Data.Common.WinnerId]
	game.Walk(0, func(player *pb.RummyRoomPlayer) bool {
		if player.State != pb.RummyPlayState_Rummy_PLAYSTATE_WIN {

			switch game.Data.RoomCfg.GameType {
			case pb.GameType_GameTypeDR:
				player.Total -= player.Private.Score
				scoreTmp += player.Private.Score
			default:
				player.Total += player.Private.Score
				if player.Total >= int64(game.Data.RoomCfg.OutLimit) {
					player.State = pb.RummyPlayState_Rummy_PLAYSTATE_ELIMINATED
				}
			}

			player_info := game.GetPlayerSettleInfo(player)
			record.Players = append(record.Players, player_info)
		}
		return true
	})

	switch game.Data.RoomCfg.GameType {
	case pb.GameType_GameTypeDR:
		winner.Total += scoreTmp
	}
	winner_info := game.GetPlayerSettleInfo(winner)
	record.Players = append(record.Players, winner_info)
	// 更新玩家total和淘汰状态 结束

	// 保存小局记录
	if game.Data.Records == nil {
		size := uint32(5)
		if game.Data.RoomCfg.Deals > 0 {
			size = game.Data.RoomCfg.Deals
		}
		game.Data.Records = make([]*pb.RummyMatchInfo, 0, size)
	}
	game.Data.Records = append(game.Data.Records, record)

	// 判断总结算或瓜分奖池
	playerLen := game.GetContinues()
	game.Data.HalfState = pb.HalfState_Init
	if playerLen <= 1 {
		game.Data.IsEnd = true
	} else if playerLen == math.MaxInt {
		game.Data.HalfState = pb.HalfState_Allow
	}

	if !game.Data.IsEnd {
		//  发送小局结算通知
		ntf := &pb.RummyExtSettleNtf{
			RoomId:    game.GetRoomId(),
			GhostCard: game.Data.Common.GhostCard,
			Records:   game.Data.Records,
			TimeOut:   game.Data.Common.TimeOut,
			TotalTime: game.Data.Common.TotalTime,
			HalfState: game.Data.HalfState,
		}
		err := game.NotifyToClient(game.GetPlayerUidList(), pb.RummyEventType_RummySettle, ntf)
		mlog.Infof("game settle send ntf err: %v", err)
	}
}

func (d *RummySettleState) OnExit(nowMs int64, curState pb.GameState, extra interface{}) {}
func (d *RummySettleState) OnTick(nowMs int64, curState pb.GameState, extra interface{}) pb.GameState {
	game := extra.(*rummy.RummyGame)
	if game.Data.HalfState == pb.HalfState_Sure || game.Data.IsEnd {
		return pb.GameState_RummyExt_STAGE_FIN_SETTLE
	}
	return d.MoveStateByTime(nowMs, curState, extra)
}
