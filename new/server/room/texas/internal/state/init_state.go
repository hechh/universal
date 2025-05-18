package state

import (
	"encoding/json"
	"poker_server/common/pb"
	"poker_server/framework/library/mlog"
	"poker_server/server/room/texas/domain"
	"poker_server/server/room/texas/internal/base"
)

type TexasInitState struct{}

// 初始化
func (d *TexasInitState) OnEnter(nowMs int64, curState pb.TexasGameState, extra domain.IRoom) {
	// 设置房间状态
	roomData := extra.GetTexasRoomData()
	roomData.BaseInfo.RoomState = pb.TexasRoomState_TRS_WAITSTART
	table := roomData.GetTable()
	table.GameData.GameState = curState

	defer func() {
		buf, _ := json.Marshal(roomData)
		mlog.Debugf("machine %d Init OnEnter: %s", roomData.RoomId, string(buf))
	}()

	// 初始化游戏
	newGameData := &pb.TexasGameData{
		GameState:     curState,
		DealerChairId: table.GameData.DealerChairId,
		PotPool:       &pb.TexasPotPoolData{},
	}
	if newGameData.DealerChairId <= 0 {
		newGameData.DealerChairId = uint32(base.Int32n(roomData.BaseInfo.MaxPlayerCount) + 1)
	}
	table.GameData = newGameData

	// 洗牌
	extra.Shuffle(1)
	extra.Deal(uint32(base.Int32n(5)), nil)

	// 初始化玩家状态
	for _, usr := range table.Players {
		usr.GameInfo = &pb.TexasPlayerGameInfo{GameState: curState}
	}

	// 变更
	extra.Change()
}

func (d *TexasInitState) OnExit(nowMs int64, curState pb.TexasGameState, room domain.IRoom) {
	roomData := room.GetTexasRoomData()
	buf, _ := json.Marshal(roomData)
	mlog.Debugf("machine %d Init OnExit: %s", roomData.RoomId, string(buf))
}

func (d *TexasInitState) OnTick(nowMs int64, curState pb.TexasGameState, room domain.IRoom) pb.TexasGameState {
	roomData := room.GetTexasRoomData()
	table := roomData.GetTable()
	table.GameData.GameState = curState

	// 判断房间是否还有效
	if roomData.BaseInfo.FinishTime*1000 <= nowMs || len(roomData.Table.ChairInfo) < 2 {
		return curState
	}

	// 获取玩家
	for _, usr := range table.Players {
		if base.IsPlayerState(usr, pb.TexasPlayerState_TPS_QUIT_ROOM) {
			delete(table.Players, usr.Uid)
		}
	}
	room.Change()
	return pb.TexasGameState_TGS_START
}
