package ext_state

import (
	"poker_server/common/pb"
	"poker_server/library/mlog"
	"poker_server/server/room/internal/internal/rummy"
	"poker_server/server/room/internal/module/card"
)

type RummyZhuangState struct {
	BaseState
}

func (d *RummyZhuangState) OnEnter(nowMs int64, curState pb.GameState, extra interface{}) {
	d.Log(curState)
	game := extra.(*rummy.RummyGame)
	game.Change()
	game.FlushExpireTime(nowMs)

	cardData := make([]*pb.RummyZhuangCard, 0, game.GetPlayerCount())
	cards := rummy.GetRummyOneCards()
	rummy.Shuffle(cards, 2) //洗牌两次

	i := 0
	var big uint32
	var zhuangID uint64

	game.Walk(0, func(player *pb.RummyRoomPlayer) bool {
		cardItem := cards[i]
		cardData = append(cardData, &pb.RummyZhuangCard{PlayerId: player.PlayerId, Card: cardItem})

		mlog.Infof("玩家 %v 卡牌 %v", player.PlayerId, card.Card(cardItem).String())
		//比点
		if card.Card(cardItem).Rank() > card.Card(big).Rank() {
			big = cardItem
			zhuangID = player.PlayerId
		} else if card.Card(cardItem).Rank() == card.Card(big).Rank() {
			if cardItem > big {
				big = cardItem
				zhuangID = player.PlayerId
			}
		}
		i++
		return true
	})

	roomData := game.Data
	roomData.Common.ZhuangId = zhuangID

	//  发送定庄通知 广播所有人
	ntf := &pb.RummySetZhuangNtf{
		RoomId:    game.GetRoomId(),
		Cards:     cardData,
		ZhuangId:  zhuangID,
		TimeOut:   uint64(game.GetCurStateTTL()),
		TotalTime: game.Data.Common.TotalTime,
	}
	err := game.NotifyToClient(game.GetPlayerUidList(), pb.RummyEventType_RummySetZhuang, ntf)
	mlog.Infof("定庄过程 %v 庄家 %v err: %v ", cardData, zhuangID, err)
}

func (d *RummyZhuangState) OnExit(nowMs int64, curState pb.GameState, extra interface{}) {

}

func (d *RummyZhuangState) OnTick(nowMs int64, curState pb.GameState, extra interface{}) pb.GameState {
	return d.MoveStateByTime(nowMs, curState, extra)
}
