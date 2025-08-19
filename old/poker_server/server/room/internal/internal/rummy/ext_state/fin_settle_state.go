package ext_state

import (
	"math"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/cluster"
	"poker_server/library/mlog"
	"poker_server/server/room/internal/internal/rummy"
)

/*
	RummyFinSettleState 最终结算
*/

type RummyFinSettleState struct {
	BaseState
}

func (d *RummyFinSettleState) OnEnter(nowMs int64, curState pb.GameState, extra interface{}) {
	d.Log(curState)
	game := extra.(*rummy.RummyGame)
	game.FlushExpireTime(nowMs)
	game.Change()
	game.Data.Status = pb.RoomStatus_RoomStatusFinished

	// 终局赢家账单
	bill := make([]*pb.RummySettlePlayerInfo, 0, game.GetPlayerCount())

	// 总奖池与服务费
	winCoin := game.Data.Common.PrizePool
	rake := winCoin * game.Data.RoomCfg.RakeRate / 10000

	// 不同玩法 不同情况完成奖金分配
	if game.Data.RoomCfg.GameType == pb.GameType_GameTypeDR {
		dpMax := int64(math.MinInt)
		cnt := make(map[int64]int64, game.GetPlayerCount())
		game.Walk(0, func(player *pb.RummyRoomPlayer) bool {
			if player.State != pb.RummyPlayState_Rummy_PLAYSTATE_ELIMINATED {
				cnt[player.Total]++
				dpMax = max(dpMax, player.Total)
			}
			return true
		})

		if cnt[dpMax] > 1 { //多人分奖池
			rake += (winCoin - rake) % cnt[dpMax]

			game.Walk(0, func(player *pb.RummyRoomPlayer) bool {
				if player.Total == dpMax && player.State != pb.RummyPlayState_Rummy_PLAYSTATE_ELIMINATED {
					prize := (game.Data.Common.PrizePool - rake) / cnt[dpMax]
					player.Coin += uint64(prize)
					bill = append(bill, game.GetExtSettlePlayerInfo(player, prize-game.Data.RoomCfg.BaseScore))
				} else {
					bill = append(bill, game.GetExtSettlePlayerInfo(player, -game.Data.RoomCfg.BaseScore))
				}
				return true
			})
		} else { //独享奖池
			game.Walk(0, func(player *pb.RummyRoomPlayer) bool {
				if player.Total == dpMax && player.State != pb.RummyPlayState_Rummy_PLAYSTATE_ELIMINATED {
					prize := game.Data.Common.PrizePool - rake
					player.Coin += uint64(prize)
					bill = append(bill, game.GetExtSettlePlayerInfo(player, prize-game.Data.RoomCfg.BaseScore))
				} else {
					bill = append(bill, game.GetExtSettlePlayerInfo(player, -game.Data.RoomCfg.BaseScore))
				}
				return true
			})
		}

	} else { // pool rummy 平分或者剩者为王
		if game.Data.HalfState == pb.HalfState_Sure {
			rake += (winCoin - rake) % 2

			game.Walk(0, func(player *pb.RummyRoomPlayer) bool {
				if player.State != pb.RummyPlayState_Rummy_PLAYSTATE_ELIMINATED {
					prize := (game.Data.Common.PrizePool - rake) / 2
					player.Coin += uint64(prize)
					bill = append(bill, game.GetExtSettlePlayerInfo(player, prize-game.Data.RoomCfg.BaseScore))
				} else {
					bill = append(bill, game.GetExtSettlePlayerInfo(player, -game.Data.RoomCfg.BaseScore))
				}
				return true
			})
		} else { // 独享奖池
			game.Walk(0, func(player *pb.RummyRoomPlayer) bool {
				if player.State != pb.RummyPlayState_Rummy_PLAYSTATE_ELIMINATED {
					prize := game.Data.Common.PrizePool - rake
					player.Coin += uint64(prize)
					bill = append(bill, game.GetExtSettlePlayerInfo(player, prize-game.Data.RoomCfg.BaseScore))
				} else {
					bill = append(bill, game.GetExtSettlePlayerInfo(player, -game.Data.RoomCfg.BaseScore))
				}
				return true
			})
		}
	}

	//  发送小局结算通知
	ntf := &pb.RummyFinSettleNtf{
		RoomId:    game.GetRoomId(),
		Records:   game.Data.Records,
		TimeOut:   game.Data.Common.TimeOut,
		TotalTime: game.Data.Common.TotalTime,
		Rake:      rake,
		Bill:      bill,
	}
	err := game.NotifyToClient(game.GetPlayerUidList(), pb.RummyEventType_RummyFinish, ntf)
	mlog.Infof("game finish settle send ntf err: %v", err)

	// 对局记录显示所有玩家id
	playerIDs := make([]uint64, 0, len(game.Data.Common.Players))
	game.Walk(0, func(player *pb.RummyRoomPlayer) bool {
		playerIDs = append(playerIDs, player.PlayerId)
		return true
	})

	// 赛局结算 抽水数据异步落地
	head := &pb.Head{
		Src: framework.NewActorRouter(game),
		Dst: framework.NewDbRouter(uint64(pb.DataType_DataTypeRummySettle), "RummyExtSettleMatchPool", "Insert"),
	}

	//异步存储对局记录
	matchReq := &pb.RummyExtSettleMatchInsertReq{
		Data: &pb.RummyExtSettleMatch{ //game.Data.Records
			RoomId:    game.GetRoomId(),
			PlayerId:  playerIDs,
			GameType:  game.Data.RoomCfg.GameType,
			RoomType:  game.Data.RoomCfg.RoomType,
			CoinType:  game.Data.RoomCfg.CoinType,
			BaseScore: game.Data.RoomCfg.BaseScore,
			RakeRate:  game.Data.RoomCfg.RakeRate,
			CreatedAt: game.Data.Match.StartTime,
			FinishAt:  game.Data.Match.EndTime,
			RakeCoin:  rake,
			Bills:     bill,
			MatchInfo: game.Data.Records,
		},
	}
	err = cluster.Send(head, matchReq)
	if err != nil {
		mlog.Infof("RummyExtSettleMatchInsert Send Error: %v", err)
	}
}

func (d *RummyFinSettleState) OnExit(nowMs int64, curState pb.GameState, extra interface{}) {
	game := extra.(*rummy.RummyGame)
	game.Data.Common.GameFinish = true //单局游戏结束
	game.Data.IsGroup = false

	game.Data.Match.Match = 1
	game.Finish()
	game.Data.Records = nil
}
func (d *RummyFinSettleState) OnTick(nowMs int64, curState pb.GameState, extra interface{}) pb.GameState {
	return d.MoveStateByTime(nowMs, curState, extra)
}
