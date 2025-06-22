package state

import (
	"poker_server/common/pb"
	"poker_server/library/mlog"
	"poker_server/server/room/internal/internal/rummy"
	"poker_server/server/room/internal/module/card"
)

/*
	FUNC: 发牌状态:负责洗牌和发牌，翻癞。
*/

type RummyDealState struct {
	BaseState
}

func (d *RummyDealState) OnEnter(nowMs int64, curState pb.GameState, extra interface{}) {
	d.Log(curState)
	game := extra.(*rummy.RummyGame)
	game.FlushExpireTime(nowMs)
	game.Change()
	cards := rummy.GetRummyCardsByCfg(game.Data.RoomCfg.Decks, game.Data.RoomCfg.Jokers)

	rummy.Shuffle(cards, 3) //洗牌3次

	idx := uint32(0)
	roomData := game.Data
	ghostCard := card.Card(cards[idx]).AddWild()
	roomData.Common.GhostCard = ghostCard
	idx++ //翻赖 wild end

	for i := idx; i < uint32(len(cards)); i++ {
		if card.Card(cards[i]).Rank() == card.Card(ghostCard).Rank() {
			cards[i] = card.Card(cards[i]).AddWild() //相同卡翻赖处理
		}
	}
	// 翻赖结束

	showCard := cards[idx] //明牌一张
	idx++

	roomData.Common.OprPlayer.OutCards = make([]uint32, 0, len(cards)-game.GetPlayerCount()*13)
	roomData.Common.OprPlayer.OutCards = append(roomData.Common.OprPlayer.OutCards, showCard) //更新明牌堆
	roomData.Private.Cards = cards                                                            //暗牌堆
	roomData.Common.ShowCard = showCard                                                       //出牌
	roomData.Common.ShowCard2 = 0

	players := game.GetPlayerMap()
	// 广播观战ntf 坐桌玩家黑名单过滤
	black_players := make(map[uint64]int, len(players))

	game.Walk(0, func(player *pb.RummyRoomPlayer) bool {
		if player.State == pb.RummyPlayState_Rummy_PLAYSTATE_PLAY {
			black_players[player.PlayerId]++
			handCards := make([]uint32, 0, 13) //13张手牌
			handCards = append(handCards, cards[idx:idx+13]...)

			// 初始手牌排序
			player.Private.HandCards = handCards
			player.Private.CardGroup = rummy.NewCardGroup(handCards)
			_, Score := rummy.CheckRCG(player.Private.HandCards, player.Private.CardGroup)
			player.Private.Score = Score
			idx += 13
			mlog.Infof("玩家 %v 卡牌 %v", player, card.CardList(handCards).String())

			// 发牌通知
			ntf := &pb.RummyDealNtf{
				RoomId:      game.GetRoomId(),
				HandCards:   player.Private.HandCards,
				GhostCard:   roomData.Common.GhostCard,
				ShowCard:    showCard,
				ShowCard2:   0,
				TimeOut:     roomData.Common.TimeOut,
				Groups:      player.Private.CardGroup,
				GroupsScore: player.Private.Score,
				TotalTime:   roomData.Common.TotalTime,
			}
			err := game.NotifyToClient([]uint64{player.PlayerId}, pb.RummyEventType_RummyDeal, ntf)
			mlog.Debugf("NotifyToClient %v err:%v", player.PlayerId, err)
		}
		return true
	})

	// 观战发送空发牌事件
	for playerId := range players {
		_, ok := black_players[playerId]
		if !ok && players[playerId].Health == pb.RummyPlayHealth_Rummy_NORMAL {
			ntf := &pb.RummyDealNtf{
				RoomId:    game.GetRoomId(),
				GhostCard: roomData.Common.GhostCard,
				ShowCard:  showCard,
				ShowCard2: 0,
				TimeOut:   roomData.Common.TimeOut,
			}
			err := game.NotifyToClient([]uint64{players[playerId].PlayerId}, pb.RummyEventType_RummyDeal, ntf)
			mlog.Debugf("NotifyToClient %v err:%v", players[playerId].PlayerId, err)
		}
	}

	roomData.Private.CardIdx = idx
}

func (d *RummyDealState) OnExit(nowMs int64, curState pb.GameState, extra interface{}) {
}

func (d *RummyDealState) OnTick(nowMs int64, curState pb.GameState, extra interface{}) pb.GameState {
	return d.MoveStateByTime(nowMs, curState, extra)
}
