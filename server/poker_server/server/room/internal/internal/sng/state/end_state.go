package state

import (
	"encoding/json"
	"poker_server/common/config/repository/machine_config"
	"poker_server/common/config/repository/texas_config"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/library/mlog"
	"poker_server/server/room/internal/internal/sng"
	"time"
)

type EndState struct{}

func (d *EndState) getUsers(room *sng.SngTexasGame, tmps map[uint64]*pb.TexasGameEndInfo, event *pb.TexasGameEventNotify) (users []*pb.TexasPlayerData) {
	table := room.GetTable()
	// 遍历结算
	room.Walk(0, func(usr *pb.TexasPlayerData) bool {
		// 获取结算玩家
		if usr.GameInfo.Operate != pb.OperateType_ALL_IN && usr.GameInfo.Operate != pb.OperateType_FOLD {
			users = append(users, usr)
		}
		// 获取通知消息
		tmps[usr.Uid] = &pb.TexasGameEndInfo{
			Uid:      usr.Uid,
			ChairId:  usr.ChairId,
			Chips:    usr.Chips,
			CardType: int32(usr.GameInfo.BestCardType),
			Bests:    usr.GameInfo.BestCardList,
		}
		if (table.GameData.IsCompare && usr.GameInfo.Operate != pb.OperateType_FOLD) ||
			(usr.GameInfo.GameState == pb.GameState_SNG_TEXAS_RIVER_ROUND && usr.GameInfo.Operate != pb.OperateType_FOLD) {
			tmps[usr.Uid].Hands = usr.GameInfo.HandCardList
		}
		event.EndInfo = append(event.EndInfo, tmps[usr.Uid])
		return true
	})
	return
}

func (d *EndState) OnEnter(nowMs int64, curState pb.GameState, extra interface{}) {
	room := extra.(*sng.SngTexasGame)
	table := room.GetTable()
	table.CurState = curState
	record := room.GetRecord()
	texasCfg := texas_config.MGetID(room.GameId)

	defer func() {
		buf, _ := json.Marshal(room.TexasRoomData)
		mlog.Debugf("roomId:%d,round:%d,End OnEnter: %s", room.RoomId, room.Table.Round, string(buf))
	}()

	// 获取结算玩家
	tmps := map[uint64]*pb.TexasGameEndInfo{}
	event := &pb.TexasGameEventNotify{
		RoomId:  room.RoomId,
		Round:   table.Round,
		PotPool: table.GameData.PotPool,
	}
	users := d.getUsers(room, tmps, event)
	room.UpdateMain(users...)

	// 奖励结算
	winners := map[uint64]struct{}{}
	chips := map[uint64]int64{}
	srvs := map[uint64]int64{}
	serviceChips := room.Reward(winners, chips, srvs)
	for uid := range winners {
		usr := room.Table.Players[uid]
		usr.Chips += chips[usr.Uid]
		endInfo := tmps[usr.Uid]
		endInfo.WinChips = chips[usr.Uid]
		endInfo.Chips = usr.Chips
	}
	mlog.Debugf("-------->winners:%v, chips:%v, srvs:%v, room:%v", winners, chips, srvs, room.TexasRoomData)
	// 发送通知
	uids := room.GetPlayerUidList()
	newHead := &pb.Head{
		Src: framework.NewSrcRouter(pb.RouterType_RouterTypeRoomId, room.RoomId),
		Cmd: uint32(pb.CMD_TEXAS_EVENT_NOTIFY),
	}
	framework.NotifyToClient(uids, newHead, sng.NewTexasEventNotify(pb.TexasEventType_EVENT_GAME_END, event))

	// 添加结束日志
	room.TotalServiceChips += serviceChips
	room.TotalRuningWater += room.Table.GameData.PotPool.TotalBetChips
	record.EndTime = time.Now().UnixMilli()
	record.TotalPot = table.GameData.PotPool.TotalBetChips
	record.TotalService = serviceChips
	dealer := room.GetDealer()
	dst := framework.NewDbRouter(uint64(pb.DataType_DataTypeReport), "ReportDataMgr", "TexasPlayerFlowReport")
	head := framework.NewHead(dst, pb.RouterType_RouterTypeRoomId, room.RoomId)
	room.Walk(int(dealer.GameInfo.Position), func(usr *pb.TexasPlayerData) bool {
		usr.TotalIncr += chips[usr.Uid]
		item := &pb.TexasGamePlayerRecordInfo{
			Uid:          usr.Uid,
			ChairId:      usr.ChairId,
			Chips:        usr.Chips,
			HandCardList: tmps[usr.Uid].Hands,
		}
		if _, ok := winners[usr.Uid]; ok {
			item.BestCardList = tmps[usr.Uid].Bests
			item.CardType = pb.CardType(tmps[usr.Uid].CardType)
			item.WinChips = tmps[usr.Uid].WinChips
		}
		record.PlayerRecord.List = append(record.PlayerRecord.List, item)
		// 上报数据
		framework.Send(head, &pb.TexasPlayerFlowReport{
			Uid:          usr.Uid,
			RoomId:       room.RoomId,
			Round:        table.Round,
			GameType:     texasCfg.GameType,
			RoomType:     texasCfg.RoomType,
			CoinType:     texasCfg.CoinType,
			BeginTime:    record.BeginTime,
			EndTime:      record.EndTime,
			Chips:        usr.Chips,
			Incr:         chips[usr.Uid],
			ServiceChips: srvs[usr.Uid],
		})
		return true
	})

	dst.FuncName = "TexasGameReport"
	framework.Send(head, record)
	room.RoomState = pb.RoomStatus_RoomStatusFinished

	// 站起
	room.Walk(0, func(usr *pb.TexasPlayerData) bool {
		if usr.ChairId > 0 && (usr.PlayerState == pb.PlayerStatus_PlayerStatusQuitGame || texasCfg.BigBlind > usr.Chips) {
			framework.NotifyToClient(uids, newHead, sng.NewTexasEventNotify(pb.TexasEventType_EVENT_STAND_UP, &pb.TexasPlayerEventNotify{
				RoomId:     room.RoomId,
				ChairId:    usr.ChairId,
				Uid:        usr.Uid,
				PlayerInfo: room.GetPlayerInfo(usr.Uid),
			}))
			delete(table.ChairInfo, usr.ChairId)
			usr.ChairId = 0
		}
		return true
	})
	room.Change()
}

func (t *EndState) OnTick(nowMs int64, curState pb.GameState, extra interface{}) pb.GameState {
	room := extra.(*sng.SngTexasGame)
	machine := room.GetMachine()
	room.Table.CurState = curState

	texasCfg := texas_config.MGetID(room.GameId)
	machineCfg := machine_config.MGetGameType(texasCfg.GameType)

	if nowMs-machine.GetCurStateStartTime() < machineCfg.FlopDuration*1000 {
		return curState
	}

	defer func() {
		buf, _ := json.Marshal(room.TexasRoomData)
		mlog.Debugf("roomId:%d,round:%d,End OnTick: %s", room.RoomId, room.Table.Round, string(buf))
	}()
	return room.GetNextState()
}

func (t *EndState) OnExit(nowMs int64, curState pb.GameState, extra interface{}) {
	room := extra.(*sng.SngTexasGame)
	defer func() {
		buf, _ := json.Marshal(room.TexasRoomData)
		mlog.Debugf("roomId:%d,round:%d,End OnExit: %s", room.RoomId, room.Table.Round, string(buf))
	}()
}
