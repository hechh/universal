package state

import (
	"encoding/json"
	"poker_server/common/pb"
	"poker_server/framework/library/mlog"
	"poker_server/server/room/texas/domain"
)

type RiverState struct {
	BaseState
}

func (d *RiverState) OnEnter(nowMs int64, curState pb.TexasGameState, room domain.IRoom) {
	roomData := room.GetTexasRoomData()
	table := roomData.GetTable()
	record := room.GetRecord()
	table.CurState = curState
	table.GameData.GameState = curState

	defer func() {
		buf, _ := json.Marshal(roomData)
		mlog.Debugf("machine %d River OnEnter: %s", roomData.RoomId, string(buf))
	}()

	// 发一张公共牌
	room.Deal(1, func(cursor uint32, card uint32) {
		record.Detail.DealList = append(record.Detail.DealList, &pb.TexasGamePokerDealRecord{
			DealType: pb.TexasDealType_TDT_RIVER,
			Card:     card,
			Cursor:   cursor,
		})
		table.GameData.PublicCardList = append(table.GameData.PublicCardList, card)
	})

	// 设置玩家状态
	cmp := room.GetCompare()
	dealer := room.GetDealer()
	room.Walk(int(dealer.GameInfo.Position), func(usr *pb.TexasPlayerData) bool {
		room.UpdateBest(usr, cmp, table.GameData.PublicCardList)

		switch usr.GameInfo.Operate {
		case pb.TexasOperateType_TOT_ALL_IN, pb.TexasOperateType_TOT_FOLD:
		default:
			usr.IsChange = false
			usr.GameInfo.GameState = curState
			usr.GameInfo.Operate = pb.TexasOperateType_TOT_NONE
			usr.GameInfo.BetChips = 0
			usr.GameInfo.PreOperate = pb.TexasOperateType_TOT_NONE
			usr.GameInfo.PreBetChips = 0
		}
		return true
	})

	// 判断是否全部比牌
	event := &pb.TexasDealEventNotify{
		RoomId:    (roomData.RoomId),
		GameState: (int32(curState)),
		HandsCard: table.GameData.PublicCardList,
		PotPool:   (table.GameData.PotPool),
		Duration:  (room.GetStateStartTime() + room.GetCurStateTTL() - nowMs),
	}
	table.GameData.MaxBetChips = 0
	cursor := room.GetNext(int(dealer.GameInfo.Position), curState, pb.TexasOperateType_TOT_FOLD, pb.TexasOperateType_TOT_ALL_IN)
	if cursor != nil {
		// 更新游标
		table.GameData.UidCursor = cursor.GameInfo.Position
		table.GameData.MinRaise = roomData.BaseInfo.BigBlind
		next := room.GetNext(int(cursor.GameInfo.Position), curState, pb.TexasOperateType_TOT_FOLD, pb.TexasOperateType_TOT_ALL_IN)
		if next != nil && next.Uid != cursor.Uid {
			event.CurBetChairId = (cursor.ChairId)
		}
	}
	// 发送公共牌
	room.NotifyToClient(pb.TexasEventType_EVENT_FLOP_CARD, event)
	room.Change()
}
