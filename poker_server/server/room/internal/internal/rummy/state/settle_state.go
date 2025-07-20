package state

import (
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/cluster"
	"poker_server/library/mlog"
	"poker_server/library/util"
	"poker_server/server/room/internal/internal/rummy"
)

/*
	FUNC: 结算状态:实现结算和结算展示。
*/

type RummySettleState struct {
	BaseState
}

func (d *RummySettleState) OnEnter(nowMs int64, curState pb.GameState, extra interface{}) {
	d.Log(curState)
	game := extra.(*rummy.RummyGame)
	game.FlushExpireTime(nowMs)
	game.Change()
	var playerSettles = make([]*pb.RummySettlePlayerInfo, 0, game.GetPlayerCount())
	winner := game.Data.Common.Players[game.Data.Common.WinnerId]
	game.Walk(0, func(player *pb.RummyRoomPlayer) bool {
		if player.State != pb.RummyPlayState_Rummy_PLAYSTATE_WIN {
			player.Total = max(game.Data.RoomCfg.MinBuyIn, game.Data.RoomCfg.MinBuyIn+int64(player.Coin)-player.JoinCoin) //更新total

			playerSettle := game.GetSettlePlayerInfo(player)

			playerSettles = append(playerSettles, playerSettle)
		}
		return true
	})

	// 结算奖池
	winCoin := game.Data.RoomCfg.BaseScore * game.Data.Common.OprPlayer.ScorePool
	// 抽水金额
	rake := winCoin * game.Data.RoomCfg.RakeRate / 10000
	winner.Coin += uint64(winCoin - rake)
	winner.Total = max(game.Data.RoomCfg.MinBuyIn, game.Data.RoomCfg.MinBuyIn+int64(winner.Coin)-winner.JoinCoin) //更新total

	// 赢家拿下奖池
	winnerSettle := game.GetSettlePlayerInfo(winner)
	winnerSettle.Coin = winCoin
	winnerSettle.Score = 0
	playerSettles = append(playerSettles, winnerSettle)

	mlog.Infof("赢家: %v 手牌: %v 结算记录:%v", winner.PlayerId, winner.Private.CardGroup, playerSettles)
	game.Data.Match.EndTime = nowMs

	//  发送结算通知
	ntf := &pb.RummySettleNtf{
		RoomId:    game.GetRoomId(),
		GhostCard: game.Data.Common.GhostCard,
		Players:   playerSettles,
		TimeOut:   game.Data.Common.TimeOut,
		Rake:      rake,
		TotalTime: game.Data.Common.TotalTime,
		ScorePoll: game.Data.Common.OprPlayer.ScorePool,
	}
	err := game.NotifyToClient(game.GetPlayerUidList(), pb.RummyEventType_RummySettle, ntf)
	mlog.Infof("game settle send ntf err: %v", err)

	// 赛局结算 抽水数据异步落地
	head := &pb.Head{
		Src: framework.NewActorRouter(game),
		Dst: framework.NewDbRouter(uint64(pb.DataType_DataTypeRummySettle), "RummySettleMatchPool", "Insert"),
	}
	matchReq := &pb.RummySettleMatchInsertReq{
		Data: []*pb.RummySettleMatch{
			{
				RoomId:     game.GetRoomId(),
				PlayerId:   game.Data.Common.PlayerIds,
				Match:      game.Data.Match.Match,
				GameType:   game.Data.RoomCfg.GameType,
				RoomType:   game.Data.RoomCfg.RoomType,
				CoinType:   game.Data.RoomCfg.CoinType,
				BaseScore:  game.Data.RoomCfg.BaseScore,
				RakeRate:   game.Data.RoomCfg.RakeRate,
				CreatedAt:  game.Data.Match.StartTime,
				FinishAt:   game.Data.Match.EndTime,
				RakeCoin:   rake,
				SettleCoin: winCoin,
				Players:    playerSettles,
			},
		},
	}

	util.DestructRoomId(game.GetRoomId())

	err = cluster.Send(head, matchReq)
	if err != nil {
		mlog.Infof("RummySettleMatchPool Insert Error: %v", err)
	}

	// 赛局玩家流水
	head = &pb.Head{
		Src: framework.NewActorRouter(game),
		Dst: framework.NewDbRouter(uint64(pb.DataType_DataTypeRummySettle), "RummySettlePool", "Insert"),
	}

	settleReq := &pb.RummySettleInsertReq{
		Data: make([]*pb.RummySettleData, 0, len(playerSettles)),
	}

	for i := range playerSettles { // 底分 人数 时间 输赢
		settleReq.Data = append(settleReq.Data, &pb.RummySettleData{
			PlayerId:    playerSettles[i].PlayerId,
			RoomId:      game.GetRoomId(),
			Groups:      playerSettles[i].CardGroup,
			GhostCard:   game.Data.Common.GhostCard,
			HandScore:   playerSettles[i].Score,
			Coin:        playerSettles[i].Coin,
			State:       playerSettles[i].State,
			CreatedAt:   nowMs,
			PlayerCount: uint32(game.GetPlayerCount()),
			MatchId:     game.Data.Match.Match,
			GameType:    game.Data.RoomCfg.GameType,
			RoomType:    game.Data.RoomCfg.RoomType,
			CoinType:    game.Data.RoomCfg.CoinType,
			BaseScore:   game.Data.RoomCfg.BaseScore,
			RakeCoin:    rake,
		})
	}

	err = cluster.Send(head, settleReq)
	if err != nil {
		mlog.Infof("RummySettlePool Insert Error: %v", err)
	}

	game.Data.Status = pb.RoomStatus_RoomStatusFinished
}

func (d *RummySettleState) OnExit(nowMs int64, curState pb.GameState, extra interface{}) {
	game := extra.(*rummy.RummyGame)
	game.Data.Common.GameFinish = true
}
func (d *RummySettleState) OnTick(nowMs int64, curState pb.GameState, extra interface{}) pb.GameState {
	return d.MoveStateByTime(nowMs, curState, extra)
}
