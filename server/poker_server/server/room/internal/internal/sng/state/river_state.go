package state

import (
	"encoding/json"
	"poker_server/common/config/repository/machine_config"
	"poker_server/common/config/repository/texas_config"
	"poker_server/common/config/repository/texas_test_config"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/library/mlog"
	"poker_server/library/util"
	"poker_server/server/room/internal/internal/sng"
	"poker_server/server/room/internal/module/card"
)

type RiverState struct {
	BaseState
}

func (d *RiverState) OnEnter(nowMs int64, curState pb.GameState, extra interface{}) {
	room := extra.(*sng.SngTexasGame)
	table := room.GetTable()
	machine := room.GetMachine()
	record := room.GetRecord()
	table.CurState = curState
	texasCfg := texas_config.MGetID(room.GameId)
	machineCfg := machine_config.MGetGameType(texasCfg.GameType)

	defer func() {
		buf, _ := json.Marshal(room.TexasRoomData)
		mlog.Debugf("roomId:%d,round:%d,River OnEnter: %s", room.RoomId, room.Table.Round, string(buf))
	}()

	// 发一张公共牌
	room.Deal(1, func(cursor uint32, cardVal uint32) {
		if texasCfg.IsTest {
			if testCfg := texas_test_config.MGetRound(texas_test_config.MGetRoundKey(room.Table.Round)); testCfg != nil {
				strCard := util.Index[string](testCfg.Publics, 4, "")
				cardVal = util.Index[uint32](card.StrToCard(strCard), 0, cardVal)
			}
		}
		record.DealRecord.List = append(record.DealRecord.List, &pb.TexasGameDealRecordInfo{
			DealType: pb.DealType_RIVER,
			Card:     cardVal,
			Cursor:   cursor,
		})
		table.GameData.PublicCardList = append(table.GameData.PublicCardList, cardVal)
	})

	// 设置玩家状态
	dealer := room.GetDealer()
	room.Walk(int(dealer.GameInfo.Position), func(usr *pb.TexasPlayerData) bool {
		room.UpdateBest(usr)

		switch usr.GameInfo.Operate {
		case pb.OperateType_ALL_IN, pb.OperateType_FOLD:
		default:
			usr.GameInfo.IsChange = false
			usr.GameInfo.GameState = curState
			usr.GameInfo.Operate = pb.OperateType_OperateNone
			usr.GameInfo.BetChips = 0
			usr.GameInfo.PreOperate = pb.OperateType_OperateNone
			usr.GameInfo.PreBetChips = 0
		}
		return true
	})

	// 判断是否全部比牌
	event := &pb.TexasDealEventNotify{
		RoomId:    room.RoomId,
		GameState: int32(curState),
		HandsCard: table.GameData.PublicCardList,
		PotPool:   table.GameData.PotPool,
		Duration:  machine.GetCurStateStartTime() + machineCfg.RiverDuration*1000 - nowMs,
	}
	table.GameData.MaxBetChips = 0
	cursor := room.GetNext(int(dealer.GameInfo.Position), curState, pb.OperateType_FOLD, pb.OperateType_ALL_IN)
	if cursor != nil {
		// 更新游标
		table.GameData.UidCursor = cursor.GameInfo.Position
		table.GameData.MinRaise = texasCfg.BigBlind
		next := room.GetNext(int(cursor.GameInfo.Position), curState, pb.OperateType_FOLD, pb.OperateType_ALL_IN)
		if next != nil && next.Uid != cursor.Uid {
			event.CurBetChairId = cursor.ChairId
		}
	}
	// 发送公共牌
	newHead := &pb.Head{
		Src: framework.NewSrcRouter(pb.RouterType_RouterTypeRoomId, room.RoomId),
		Cmd: uint32(pb.CMD_TEXAS_EVENT_NOTIFY),
	}
	framework.NotifyToClient(room.GetPlayerUidList(), newHead, sng.NewTexasEventNotify(pb.TexasEventType_EVENT_FLOP_CARD, event))
	room.Change()
}
