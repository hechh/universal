package state

import (
	"encoding/json"
	"fmt"
	"poker_server/common/config/repository/machine_config"
	"poker_server/common/config/repository/texas_config"
	"poker_server/common/pb"
	"poker_server/framework"
	"poker_server/library/mlog"
	"poker_server/server/room/internal/internal/texas"
)

type StartState struct{}

func (d *StartState) OnEnter(nowMs int64, curState pb.GameState, extra interface{}) {
	room := extra.(*texas.TexasGame)
	room.Table.CurState = curState
	room.RoomState = pb.RoomStatus_RoomStatusPlaying

	defer func() {
		buf, _ := json.Marshal(room.TexasRoomData)
		mlog.Debugf("roomId:%d,round:%d,Start OnEnter: %s", room.RoomId, room.Table.Round, string(buf))
	}()

	users := room.GetGamePlayers()
	flag := true
	var dealer *pb.TexasPlayerData
	for i, usr := range users {
		usr.TotalTimes++
		usr.GameInfo.Position = uint32(i)
		usr.GameInfo.GameState = curState
		room.Table.GameData.UidList = append(room.Table.GameData.UidList, usr.Uid)
		if flag && room.Table.GameData.DealerChairId < usr.ChairId {
			dealer = usr
			flag = false
		}
	}
	if dealer == nil {
		dealer = users[0]
	}

	texasCfg := texas_config.MGetID(room.GameId)
	room.Table.GameData.DealerChairId = dealer.ChairId
	room.Table.GameData.SmallChairId = users[int(dealer.GameInfo.Position+1)%len(users)].ChairId
	room.Table.GameData.BigChairId = users[int(dealer.GameInfo.Position+2)%len(users)].ChairId
	room.Table.GameData.UidCursor = uint32(int(dealer.GameInfo.Position+3) % len(users))
	room.Table.Round++
	room.Operate(room.GetSmall(), pb.OperateType_BET_SMALL_BLIND, texasCfg.SmallBlind)
	room.Operate(room.GetBig(), pb.OperateType_BET_BIG_BLIND, texasCfg.BigBlind)
	room.Change()
}

func (d *StartState) OnTick(nowMs int64, curState pb.GameState, extra interface{}) pb.GameState {
	room := extra.(*texas.TexasGame)
	room.Table.CurState = curState

	defer func() {
		buf, _ := json.Marshal(room.TexasRoomData)
		mlog.Debugf("roomId:%d,round:%d,Start OnTick: %s", room.RoomId, room.Table.Round, string(buf))
	}()

	texasCfg := texas_config.MGetID(room.GameId)
	machineCfg := machine_config.MGetGameType(texasCfg.GameType)
	uids := room.GetPlayerUidList()
	newHead := &pb.Head{
		Src: framework.NewSrcRouter(pb.RouterType_RouterTypeRoomId, room.RoomId),
		Cmd: uint32(pb.CMD_TEXAS_EVENT_NOTIFY),
	}
	framework.NotifyToClient(uids, newHead, texas.NewTexasEventNotify(pb.TexasEventType_EVENT_GAME_START, &pb.TexasGameEventNotify{
		RoomId:        room.RoomId,
		Round:         room.Table.Round,
		BigChair:      room.Table.GameData.BigChairId,
		SmallChair:    room.Table.GameData.SmallChairId,
		DealerChair:   room.Table.GameData.DealerChairId,
		SmallChip:     uint32(texasCfg.SmallBlind),
		BigChip:       uint32(texasCfg.BigBlind),
		CurBetChairId: room.GetCursor().ChairId,
		PotPool:       room.Table.GameData.PotPool,
		Duration:      room.GetMachine().GetCurStateStartTime() + machineCfg.StartDuration*1000 - nowMs,
	}))

	// 添加游戏日志
	room.SetRecord(&pb.TexasGameReport{
		Round:     room.Table.Round,
		RoomId:    room.RoomId,
		GameType:  texasCfg.GameType,
		RoomType:  texasCfg.RoomType,
		Blind:     fmt.Sprintf("%d/%d", texasCfg.SmallBlind, texasCfg.BigBlind),
		BeginTime: nowMs,
		OperateRecord: &pb.TexasGameOperateRecord{
			List: []*pb.TexasGameOperateRecordInfo{
				{GameState: curState, Uid: room.GetSmall().Uid, Operate: pb.OperateType_BET_SMALL_BLIND, BetChips: texasCfg.SmallBlind},
				{GameState: curState, Uid: room.GetBig().Uid, Operate: pb.OperateType_BET_BIG_BLIND, BetChips: texasCfg.BigBlind},
			},
		},
	})

	return pb.GameState_TEXAS_PRE_FLOP
}

func (d *StartState) OnExit(nowMs int64, curState pb.GameState, extra interface{}) {
	room := extra.(*texas.TexasGame)
	defer func() {
		buf, _ := json.Marshal(room.TexasRoomData)
		mlog.Debugf("roomId:%d,round:%d,Start OnExit: %s", room.RoomId, room.Table.Round, string(buf))
	}()
}
