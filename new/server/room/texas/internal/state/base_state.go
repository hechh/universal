package state

import (
	"encoding/json"
	"poker_server/common/pb"
	"poker_server/framework/library/mlog"
	"poker_server/server/room/texas/domain"
	"poker_server/server/room/texas/internal/base"
)

type BaseState struct{}

func (d *BaseState) OnTick(nowMs int64, curState pb.TexasGameState, room domain.IRoom) pb.TexasGameState {
	roomData := room.GetTexasRoomData()
	table := roomData.GetTable()
	record := room.GetRecord()

	// 是否直接比牌
	if table.GameData.IsCompare {
		if nowMs-room.GetStateStartTime() < domain.HOLD_TIME {
			return curState
		}
		return room.GetNextState()
	}

	// 玩家超时处理
	usr := room.GetCursor()
	aprev := room.GetPrev(int(usr.GameInfo.Position), curState, pb.TexasOperateType_TOT_FOLD)
	next := room.GetNext(int(usr.GameInfo.Position), curState, pb.TexasOperateType_TOT_FOLD, pb.TexasOperateType_TOT_ALL_IN)
	if duration := room.GetCurStateTTL() + room.GetStateStartTime() - nowMs; duration <= 0 {
		base.SetPlayerState(usr, pb.TexasPlayerState_TPS_QUIT_TABLE)
		if next == nil {
			if aprev != nil && aprev.GameInfo.BetChips == table.GameData.MaxBetChips && usr.GameInfo.BetChips == table.GameData.MaxBetChips {
				room.Operate(usr, pb.TexasOperateType_TOT_CHECK, 0)
			} else {
				room.Operate(usr, pb.TexasOperateType_TOT_FOLD, 0)
			}
		} else {
			if aprev != nil && aprev.GameInfo.BetChips == table.GameData.MaxBetChips && usr.GameInfo.BetChips == table.GameData.MaxBetChips && next.GameInfo.BetChips == usr.GameInfo.BetChips {
				room.Operate(usr, pb.TexasOperateType_TOT_CHECK, 0)
			} else {
				room.Operate(usr, pb.TexasOperateType_TOT_FOLD, 0)
			}
		}
	}

	// 是否有操作
	if !usr.IsChange {
		return curState
	}

	defer func() {
		buf, _ := json.Marshal(roomData)
		switch curState {
		case pb.TexasGameState_TGS_PRE_FLOP:
			mlog.Debugf("machine %d Preflop OnTick: %s", roomData.RoomId, string(buf))
		case pb.TexasGameState_TGS_FLOP_ROUND:
			mlog.Debugf("machine %d Flop OnTick: %s", roomData.RoomId, string(buf))
		case pb.TexasGameState_TGS_TURN_ROUND:
			mlog.Debugf("machine %d Turn OnTick: %s", roomData.RoomId, string(buf))
		case pb.TexasGameState_TGS_RIVER_ROUND:
			mlog.Debugf("machine %d River OnTick: %s", roomData.RoomId, string(buf))
		}
	}()
	// 添加操作记录
	record.Detail.OperateList = append(record.Detail.OperateList, &pb.TexasGameOperateRecord{
		GameState: curState,
		Uid:       usr.Uid,
		Operate:   usr.GameInfo.Operate,
		BetChips:  usr.GameInfo.BetChips,
	})

	// 先处理玩家操作，进行状态转移
	usr.IsChange = false
	room.SetStateStartTime(nowMs)
	event := &pb.TexasBetEventNotify{
		RoomId:      roomData.RoomId,
		ChairId:     usr.ChairId,
		Chips:       usr.Chips,
		OperateType: int32(usr.GameInfo.Operate),
		BetChips:    usr.GameInfo.BetChips,
		PotPool:     table.GameData.PotPool,
		MinRaise:    table.GameData.MinRaise,
		MaxBetChips: table.GameData.MaxBetChips,
		Duration:    room.GetStateStartTime() + room.GetCurStateTTL() - nowMs,
	}
	defer room.NotifyToClient(pb.TexasEventType_EVENT_BET, event)

	// 状态转移
	bigPlayer := room.GetBig()
	anext := room.GetNext(int(usr.GameInfo.Position), curState, pb.TexasOperateType_TOT_FOLD)
	switch usr.GameInfo.Operate {
	case pb.TexasOperateType_TOT_FOLD:
		// 无任何活跃玩家
		if next == nil {
			if aprev == nil || aprev.Uid == anext.Uid {
				return pb.TexasGameState_TGS_END
			} else {
				table.GameData.IsCompare = true
				return room.GetNextState()
			}
		}
		// 只剩一个all in玩家 或者 没有任何all in玩家
		if aprev == nil || aprev.Uid == next.Uid {
			return pb.TexasGameState_TGS_END
		}
		// 筹码追平，结束当前论下注
		if aprev.GameInfo.BetChips == usr.GameInfo.BetChips && usr.GameInfo.BetChips == next.GameInfo.BetChips && next.GameInfo.Operate != pb.TexasOperateType_TOT_NONE {
			return room.GetNextState()
		}
		// 下一个玩家发话
		table.GameData.UidCursor = next.GameInfo.Position
		event.NextChairId = next.ChairId

	case pb.TexasOperateType_TOT_ALL_IN:
		// 无任何活跃玩家
		if next == nil {
			if aprev == nil {
				return pb.TexasGameState_TGS_END
			} else {
				table.GameData.IsCompare = true
				return room.GetNextState()
			}
		}
		// 下一个玩家发话
		table.GameData.UidCursor = next.GameInfo.Position
		event.NextChairId = next.ChairId
	default:
		// 无任何活跃玩家
		if next == nil {
			if aprev == nil {
				return pb.TexasGameState_TGS_END
			} else {
				table.GameData.IsCompare = true
				return room.GetNextState()
			}
		}
		// 结束当前轮下注
		prev := room.GetPrev(int(usr.GameInfo.Position), curState, pb.TexasOperateType_TOT_FOLD, pb.TexasOperateType_TOT_ALL_IN)
		if !(table.GameData.GameState == pb.TexasGameState_TGS_PRE_FLOP && roomData.BaseInfo.BigBlind == usr.GameInfo.BetChips && next.Uid == bigPlayer.Uid) &&
			prev.GameInfo.BetChips == usr.GameInfo.BetChips && next.GameInfo.BetChips == usr.GameInfo.BetChips && next.GameInfo.Operate != pb.TexasOperateType_TOT_NONE {
			return room.GetNextState()
		}
		// 继续执行当前轮
		table.GameData.UidCursor = next.GameInfo.Position
		event.NextChairId = next.ChairId
	}
	return curState
}

func (d *BaseState) OnExit(nowMs int64, curState pb.TexasGameState, room domain.IRoom) {
	room.UpdateSide(room.GetPlayers(curState)...)
}
