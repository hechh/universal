package state

import (
	"encoding/json"
	"poker_server/common/config/repository/machine_config"
	"poker_server/common/config/repository/texas_config"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/framework/cluster"
	"poker_server/library/mlog"
	"poker_server/server/room/internal/internal/sng"
	"poker_server/server/room/internal/internal/sng/util"
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
	winners := map[uint64]int64{}
	chips := map[uint64]int64{}
	srvs := map[uint64]int64{}
	serviceChips := room.Reward(winners, chips, srvs)
	for uid, winVal := range winners {
		usr := room.Table.Players[uid]
		usr.Chips += winVal
		endInfo := tmps[usr.Uid]
		endInfo.IsWinner = true
		endInfo.WinChips = winVal + chips[uid]
		endInfo.Chips = usr.Chips
	}
	mlog.Debugf("texas_end_state chips:%v, srvs:%v, winners:%v", chips, srvs, winners)
	rankList := room.UpdateRankList()

	// 发送通知
	uids := room.GetPlayerUidList()
	newHead := &pb.Head{
		Src: framework.NewSrcRouter(room.RoomId, room.GetActorName()),
		Cmd: uint32(pb.CMD_TEXAS_EVENT_NOTIFY),
	}
	cluster.SendToClient(newHead, sng.NewTexasEventNotify(pb.TexasEventType_EVENT_GAME_END, event), uids...)

	// 添加结束日志
	room.TotalServiceChips += serviceChips
	room.TotalRuningWater += room.Table.GameData.PotPool.TotalBetChips
	record.EndTime = time.Now().UnixMilli()
	record.TotalPot = table.GameData.PotPool.TotalBetChips
	record.TotalService = serviceChips
	dealer := room.GetDealer()
	head := &pb.Head{
		Src: framework.NewSrcRouter(room.RoomId, room.GetActorName()),
		Dst: framework.NewDbRouter(uint64(pb.DataType_DataTypeReport), "ReportDataMgr", "TexasPlayerFlowReport"),
	}
	room.Walk(int(dealer.GameInfo.Position), func(usr *pb.TexasPlayerData) bool {
		usr.TotalIncr += (chips[usr.Uid] + winners[usr.Uid])
		record.PlayerRecord.List = append(record.PlayerRecord.List, &pb.TexasGamePlayerRecordInfo{
			Uid:          usr.Uid,
			ChairId:      usr.ChairId,
			Chips:        usr.Chips,
			WinChips:     (chips[usr.Uid] + winners[usr.Uid]),
			ServiceChips: srvs[usr.Uid],
			HandCardList: usr.GameInfo.HandCardList,
			CardType:     usr.GameInfo.BestCardType,
			BestCardList: usr.GameInfo.BestCardList,
		})
		// 上报数据
		cluster.Send(head, &pb.TexasPlayerFlowReport{
			Uid:          usr.Uid,
			RoomId:       room.RoomId,
			Round:        table.Round,
			GameType:     texasCfg.GameType,
			RoomType:     texasCfg.RoomType,
			CoinType:     texasCfg.CoinType,
			BeginTime:    record.BeginTime,
			EndTime:      record.EndTime,
			Chips:        usr.Chips,
			Incr:         (chips[usr.Uid] + winners[usr.Uid]),
			ServiceChips: srvs[usr.Uid],
		})
		return true
	})

	head.Dst.FuncName = "TexasGameReport"
	cluster.Send(head, record)
	room.RoomState = pb.RoomStatus_RoomStatusFinished

	// 站起 + 排名
	ranks := room.GetRankList(rankList)
	room.Walk(0, func(usr *pb.TexasPlayerData) bool {
		if usr.ChairId > 0 && (usr.PlayerState == pb.PlayerStatus_PlayerStatusQuitGame || texasCfg.BigBlind > usr.Chips) {
			cluster.SendToClient(newHead, sng.NewTexasEventNotify(pb.TexasEventType_EVENT_STAND_UP, &pb.TexasPlayerEventNotify{
				RoomId:     room.RoomId,
				ChairId:    usr.ChairId,
				Uid:        usr.Uid,
				PlayerInfo: room.GetPlayerInfo(usr.Uid),
				RankInfo:   ranks[usr.Uid],
			}), uids...)
			delete(table.ChairInfo, usr.ChairId)
			usr.ChairId = 0
		}
		return true
	})
	if len(room.Table.ChairInfo) <= 1 {
		for chairId, uid := range room.Table.ChairInfo {
			cluster.SendToClient(newHead, sng.NewTexasEventNotify(pb.TexasEventType_EVENT_STAND_UP, &pb.TexasPlayerEventNotify{
				RoomId:     room.RoomId,
				ChairId:    chairId,
				Uid:        uid,
				PlayerInfo: room.GetPlayerInfo(uid),
				RankInfo:   ranks[uid],
			}), uids...)
		}
	}
	room.Change()
}

func (t *EndState) OnTick(nowMs int64, curState pb.GameState, extra interface{}) pb.GameState {
	room := extra.(*sng.SngTexasGame)
	machine := room.GetMachine()
	room.Table.CurState = curState

	texasCfg := texas_config.MGetID(room.GameId)
	machineCfg := machine_config.MGetGameType(texasCfg.GameType)

	if nowMs-machine.GetCurStateStartTime() < util.GetCurStateTTL(machineCfg, curState) {
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
