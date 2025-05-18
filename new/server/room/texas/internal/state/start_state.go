package state

import (
	"encoding/json"
	"fmt"
	"poker_server/common/pb"
	"poker_server/framework/library/mlog"
	"poker_server/server/room/texas/domain"
	"sort"
)

type TexasStartState struct{}

func (d *TexasStartState) OnEnter(nowMs int64, curState pb.TexasGameState, room domain.IRoom) {
	roomData := room.GetTexasRoomData()
	table := roomData.GetTable()
	table.CurState = curState

	defer func() {
		room.Change()
		buf, _ := json.Marshal(roomData)
		mlog.Debugf("machine %d Start OnEnter: %s", roomData.RoomId, string(buf))
	}()

	// 获取加入游戏玩家
	users := []*pb.TexasPlayerData{}
	for _, uid := range table.ChairInfo {
		users = append(users, table.Players[uid])
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].ChairId < users[j].ChairId
	})

	// 初始化游戏玩家
	flag := true
	var dealer *pb.TexasPlayerData
	for i, usr := range users {
		usr.GameInfo.InPlaying = true
		table.GameData.UidList = append(table.GameData.UidList, usr.Uid)
		usr.GameInfo.Position = uint32(i)
		usr.GameInfo.GameState = pb.TexasGameState(curState)
		if flag && table.GameData.DealerChairId < usr.ChairId {
			dealer = usr
			flag = false
		}
	}
	if dealer == nil {
		dealer = users[0]
	}

	// 初始化游戏状态
	table.GameData.GameState = pb.TexasGameState(curState)
	table.GameData.DealerChairId = dealer.ChairId
	table.GameData.SmallChairId = users[int(dealer.GameInfo.Position+1)%len(users)].ChairId
	table.GameData.BigChairId = users[int(dealer.GameInfo.Position+2)%len(users)].ChairId
	table.GameData.UidCursor = uint32(int(dealer.GameInfo.Position+3) % len(users))
	table.Round++

	// 设置大小盲
	room.Operate(room.GetSmall(), pb.TexasOperateType_TOT_BET_SMALL_BLIND, roomData.BaseInfo.SmallBlind)
	room.Operate(room.GetBig(), pb.TexasOperateType_TOT_BET_BIG_BLIND, roomData.BaseInfo.BigBlind)
}

func (d *TexasStartState) OnExit(nowMs int64, curState pb.TexasGameState, room domain.IRoom) {
	roomData := room.GetTexasRoomData()
	roomData.Table.CurState = curState

	buf, _ := json.Marshal(roomData)
	mlog.Debugf("machine %d Start OnExit: %s", roomData.RoomId, string(buf))
}

func (d *TexasStartState) OnTick(nowMs int64, curState pb.TexasGameState, room domain.IRoom) pb.TexasGameState {
	roomData := room.GetTexasRoomData()
	table := roomData.GetTable()
	table.CurState = curState
	table.GameData.GameState = curState

	defer func() {
		room.Change()
		buf, _ := json.Marshal(roomData)
		mlog.Debugf("machine %d Start OnTick: %s", roomData.RoomId, string(buf))
	}()

	roomData.TotalJoinCount = int64(len(table.Players))
	roomData.BaseInfo.RoomState = pb.TexasRoomState_TRS_PLAYING

	// 通知客户端
	room.NotifyToClient(pb.TexasEventType_EVENT_GAME_START, &pb.TexasGameEventNotify{
		RoomId:        (roomData.RoomId),
		Round:         (table.Round),
		BigChair:      (room.GetBig().ChairId),
		SmallChair:    (room.GetSmall().ChairId),
		DealerChair:   (room.GetDealer().ChairId),
		SmallChip:     (uint32(roomData.BaseInfo.SmallBlind)),
		BigChip:       (uint32(roomData.BaseInfo.BigBlind)),
		CurBetChairId: (room.GetCursor().ChairId),
		PotPool:       (table.GameData.PotPool),
		Duration:      (room.GetCurStateTTL() + room.GetStateStartTime() - nowMs),
	})

	// 添加游戏日志
	room.SetRecord(&pb.TexasGameRecord{
		TableId:   table.TableId,
		Round:     table.Round,
		GameType:  roomData.BaseInfo.GameType,
		RoomStage: roomData.BaseInfo.RoomStage,
		Blind:     fmt.Sprintf("%d/%d", roomData.BaseInfo.SmallBlind, roomData.BaseInfo.BigBlind),
		BeginTime: nowMs,
		Detail: &pb.TexasGameRecordDetail{
			OperateList: []*pb.TexasGameOperateRecord{
				{GameState: curState, Uid: room.GetSmall().Uid, Operate: pb.TexasOperateType_TOT_BET_SMALL_BLIND, BetChips: roomData.BaseInfo.SmallBlind},
				{GameState: curState, Uid: room.GetBig().Uid, Operate: pb.TexasOperateType_TOT_BET_BIG_BLIND, BetChips: roomData.BaseInfo.BigBlind},
			},
		},
	})

	return pb.TexasGameState_TGS_PRE_FLOP
}
