package state

import (
	"poker_server/common/config/repository/texas_config"
	"poker_server/common/pb"
	"poker_server/library/mlog"
	"poker_server/library/random"
	"poker_server/server/room/internal/internal/texas"
)

type InitState struct{}

// 初始化
func (d *InitState) OnEnter(nowMs int64, curState pb.GameState, extra interface{}) {
	room := extra.(*texas.TexasGame)
	room.Table.CurState = curState
	room.RoomState = pb.RoomStatus_RoomStatusWait
	texasCfg := texas_config.MGetID(room.GameId)

	defer func() {
		mlog.Debugf("roomId:%d,round:%d,Init OnEnter: %s", room.RoomId, room.Table.Round, room.String())
	}()

	newGame := &pb.TexasGameData{PotPool: &pb.TexasPotPoolData{}}
	if room.Table.GameData != nil {
		newGame.DealerChairId = room.Table.GameData.DealerChairId
	}
	if newGame.DealerChairId <= 0 {
		newGame.DealerChairId = uint32(random.Int32n(texasCfg.MaxPlayerCount)) + 1
	}
	room.Table.GameData = newGame

	room.Shuffle(2)
	room.Deal(uint32(random.Int32n(5)), nil)
	for _, usr := range room.Table.Players {
		usr.GameInfo = &pb.TexasPlayerGameInfo{GameState: curState}
	}
}

func (d *InitState) OnTick(nowMs int64, curState pb.GameState, extra interface{}) pb.GameState {
	room := extra.(*texas.TexasGame)
	texasCfg := texas_config.MGetID(room.GameId)
	startTime := room.GetMachine().GetCurStateStartTime()
	if room.HasFinished() {
		return curState
	}
	if startTime+10*60*1000 <= nowMs || texasCfg.RoomKeepLive*60+room.CreateTime <= nowMs/1000 {
		room.Finish()
		return curState
	}
	if len(room.Table.ChairInfo) < 2 {
		return curState
	}
	defer func() {
		mlog.Debugf("roomId:%d,round:%d,Init OnTick: %s", room.RoomId, room.Table.Round, room.String())
	}()
	return pb.GameState_TEXAS_START
}

func (d *InitState) OnExit(nowMs int64, curState pb.GameState, extra interface{}) {
	defer func() {
		room := extra.(*texas.TexasGame)
		mlog.Debugf("roomId:%d,round:%d,Init OnExit: %s", room.RoomId, room.Table.Round, room.String())
	}()
}
