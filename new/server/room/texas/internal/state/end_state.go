package state

import (
	"encoding/json"
	"poker_server/common/pb"
	"poker_server/framework/library/mlog"
	"poker_server/server/room/texas/domain"
	"poker_server/server/room/texas/internal/base"
	"time"
)

type EndState struct{}

func (d *EndState) getUsers(room domain.IRoom, tmps map[uint64]*pb.TexasGameEndInfo, event *pb.TexasGameEventNotify) (users []*pb.TexasPlayerData) {
	table := room.GetTexasRoomData().GetTable()
	// 遍历结算
	room.Walk(0, func(usr *pb.TexasPlayerData) bool {
		// 获取结算玩家
		if usr.GameInfo.Operate != pb.TexasOperateType_TOT_ALL_IN && usr.GameInfo.Operate != pb.TexasOperateType_TOT_FOLD {
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
		if (table.GameData.IsCompare && usr.GameInfo.Operate != pb.TexasOperateType_TOT_FOLD) ||
			(usr.GameInfo.GameState == pb.TexasGameState_TGS_RIVER_ROUND && usr.GameInfo.Operate != pb.TexasOperateType_TOT_FOLD) {
			tmps[usr.Uid].Hands = usr.GameInfo.HandCardList
		}
		event.EndInfo = append(event.EndInfo, tmps[usr.Uid])
		return true
	})
	return
}

func (d *EndState) OnEnter(nowMs int64, curState pb.TexasGameState, room domain.IRoom) {
	roomData := room.GetTexasRoomData()
	table := roomData.GetTable()
	record := room.GetRecord()
	table.CurState = curState
	table.GameData.GameState = curState

	defer func() {
		buf, _ := json.Marshal(roomData)
		mlog.Debugf("machine %d End OnEnter: %s", roomData.RoomId, string(buf))
	}()

	// 获取结算玩家
	tmps := map[uint64]*pb.TexasGameEndInfo{}
	event := &pb.TexasGameEventNotify{
		RoomId:  roomData.RoomId,
		Round:   table.Round,
		PotPool: table.GameData.PotPool,
	}
	users := d.getUsers(room, tmps, event)

	// 更新主池
	room.UpdateMain(users...)

	// 奖励结算
	winners := map[uint64]struct{}{}
	serviceChips := room.Reward(len(table.GameData.UidList), func(uid uint64, addChips int64) {
		winners[uid] = struct{}{}
		// 统计赢家
		usr := table.Players[uid]
		usr.Chips += addChips
		endInfo := tmps[uid]
		endInfo.WinChips = addChips
		endInfo.Chips = usr.Chips
	})

	// 发送通知
	room.NotifyToClient(pb.TexasEventType_EVENT_GAME_END, event)

	// 添加结束日志
	record.EndTime = time.Now().UnixMilli()
	record.TotalPot = table.GameData.PotPool.TotalBetChips
	record.TotalService = serviceChips
	dealer := room.GetDealer()
	room.Walk(int(dealer.GameInfo.Position), func(usr *pb.TexasPlayerData) bool {
		item := &pb.TexasGamePlayerRecord{
			Uid:          usr.Uid,
			ChairId:      usr.ChairId,
			Chips:        usr.Chips,
			HandCardList: tmps[usr.Uid].Hands,
		}
		if _, ok := winners[usr.Uid]; ok {
			item.BestCardList = tmps[usr.Uid].Bests
			item.CardType = pb.TexasCardType(tmps[usr.Uid].CardType)
			item.WinChips = tmps[usr.Uid].WinChips
		}
		record.Detail.PlayerList = append(record.Detail.PlayerList, item)

		// 上报数据
		/*
			room.Report(&pb.TexasPlayerRecord{
				Uid:          usr.Uid,
				TableId:      roomData.Table.TableId,
				Round:        table.Round,
				GameType:     roomData.BaseInfo.GameType,
				RoomId:       roomData.RoomId,
				RoomStage:    roomData.BaseInfo.RoomStage,
				RoomName:     roomData.BaseInfo.Name,
				BeginTime:    record.BeginTime,
				EndTime:      record.EndTime,
				Chips:        usr.Chips,
				WinChips:     item.WinChips,
				ServiceChips: serviceChips / int64(len(winners)),
			})
		*/
		return true
	})

	// 统计数据
	roomData.TotalServiceChips += serviceChips
	roomData.TotalRuningWater += table.GameData.PotPool.TotalBetChips
	roomData.BaseInfo.RoomState = pb.TexasRoomState_TRS_END

	// 站起
	room.Walk(0, func(usr *pb.TexasPlayerData) bool {
		if usr.ChairId > 0 && (base.IsPlayerState(usr, pb.TexasPlayerState_TPS_QUIT_TABLE) || roomData.BaseInfo.BigBlind > usr.Chips) {
			delete(table.ChairInfo, usr.ChairId)
			room.NotifyToClient(pb.TexasEventType_EVENT_STAND_UP, &pb.TexasPlayerEventNotify{
				RoomId:  (roomData.RoomId),
				ChairId: (usr.ChairId),
			})
			usr.ChairId = 0
		}
		return true
	})
	room.Change()
}

func (t *EndState) OnTick(nowMs int64, curState pb.TexasGameState, room domain.IRoom) pb.TexasGameState {
	roomData := room.GetTexasRoomData()
	roomData.Table.GameData.GameState = curState

	ttl := room.GetCurStateTTL()
	beginMs := room.GetStateStartTime()
	if nowMs-beginMs < ttl {
		return curState
	}

	defer func() {
		buf, _ := json.Marshal(roomData)
		mlog.Debugf("machine %d End OnEnter: %s", roomData.RoomId, string(buf))
	}()

	return room.GetNextState()
}

func (t *EndState) OnExit(nowMs int64, curState pb.TexasGameState, room domain.IRoom) {
	roomData := room.GetTexasRoomData()
	buf, _ := json.Marshal(roomData)
	mlog.Debugf("machine %d End OnEnter: %s", roomData.RoomId, string(buf))
}
