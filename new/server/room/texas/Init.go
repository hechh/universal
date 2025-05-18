package texas

import (
	"poker_server/common/pb"
	"poker_server/server/room/texas/internal/manager"
	"poker_server/server/room/texas/internal/state"
)

var (
	roomMgr = NewTexasRoomMgr()
)

func Init() {
	manager.RegisterState(pb.TexasGameState_TGS_INIT, &state.TexasInitState{})
	manager.RegisterState(pb.TexasGameState_TGS_START, &state.TexasStartState{})
	manager.RegisterState(pb.TexasGameState_TGS_PRE_FLOP, &state.TexasPreflopState{})
	manager.RegisterState(pb.TexasGameState_TGS_FLOP_ROUND, &state.TexasFlopState{})
	manager.RegisterState(pb.TexasGameState_TGS_TURN_ROUND, &state.TurnState{})
	manager.RegisterState(pb.TexasGameState_TGS_RIVER_ROUND, &state.RiverState{})
	manager.RegisterState(pb.TexasGameState_TGS_END, &state.EndState{})
}
