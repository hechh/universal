package state

import (
	"encoding/json"
	"poker_server/common/pb"
	"poker_server/framework/library/mlog"
	"poker_server/server/room/texas/domain"
)

type TexasPreflopState struct{ BaseState }

func (d *TexasPreflopState) OnEnter(nowMs int64, curState pb.TexasGameState, room domain.IRoom) {
	roomData := room.GetTexasRoomData()
	table := roomData.GetTable()
	record := room.GetRecord()
	table.CurState = curState
	table.GameData.GameState = curState

	defer func() {
		buf, _ := json.Marshal(roomData)
		mlog.Debugf("machine %d Preflop OnEnter: %s", roomData.RoomId, string(buf))
	}()

	// 发送第一轮手牌
	small := room.GetSmall()
	room.Walk(int(small.GameInfo.Position), func(usr *pb.TexasPlayerData) bool {
		usr.IsChange = false
		usr.GameInfo.GameState = curState
		// 发牌
		room.Deal(1, func(cursor uint32, card uint32) {
			record.Detail.DealList = append(record.Detail.DealList, &pb.TexasGamePokerDealRecord{
				DealType: pb.TexasDealType_TDT_HAND,
				Uid:      usr.Uid,
				Card:     card,
				Cursor:   cursor,
			})
			usr.GameInfo.HandCardList = append(usr.GameInfo.HandCardList, card)
		})
		return true
	})

	// 发送第二轮手牌 + 发送通知
	ttl := room.GetCurStateTTL() + room.GetStateStartTime() - nowMs
	cursorPlayer := room.GetCursor()
	room.Walk(int(small.GameInfo.Position), func(usr *pb.TexasPlayerData) bool {
		room.Deal(1, func(cursor uint32, card uint32) {
			record.Detail.DealList = append(record.Detail.DealList, &pb.TexasGamePokerDealRecord{
				DealType: pb.TexasDealType_TDT_HAND,
				Uid:      usr.Uid,
				Card:     card,
				Cursor:   cursor,
			})
			usr.GameInfo.HandCardList = append(usr.GameInfo.HandCardList, card)
			// 发送广播
			room.NotifyToPlayerClient(usr.Uid, pb.TexasEventType_EVENT_DEAL, &pb.TexasDealEventNotify{
				RoomId:        (roomData.RoomId),
				GameState:     (int32(curState)),
				HandsCard:     usr.GameInfo.HandCardList,
				CurBetChairId: (cursorPlayer.ChairId),
				PotPool:       (roomData.Table.GameData.PotPool),
				Duration:      (ttl),
			})
		})
		return true
	})
	room.Change()
}
