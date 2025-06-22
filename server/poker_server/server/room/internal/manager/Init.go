package manager

import (
	"poker_server/common/pb"
	rummy "poker_server/server/room/internal/internal/rummy/state"
	sngstate "poker_server/server/room/internal/internal/sng/state"
	"poker_server/server/room/internal/internal/texas/state"
	"poker_server/server/room/internal/module/machine"
)

var (
	texasMgr = NewTexasGameMgr()
	rummyMgr = NewRummyGameMgr()
)

func Init() error {
	return nil
}

func Close() {
	texasMgr.Stop()
	rummyMgr.Stop()
}

func init() {
	machine.RegisterState(pb.GameState_TEXAS_INIT, &state.InitState{})
	machine.RegisterState(pb.GameState_TEXAS_START, &state.StartState{})
	machine.RegisterState(pb.GameState_TEXAS_PRE_FLOP, &state.PreflopState{})
	machine.RegisterState(pb.GameState_TEXAS_FLOP_ROUND, &state.FlopState{})
	machine.RegisterState(pb.GameState_TEXAS_TURN_ROUND, &state.TurnState{})
	machine.RegisterState(pb.GameState_TEXAS_RIVER_ROUND, &state.RiverState{})
	machine.RegisterState(pb.GameState_TEXAS_END, &state.EndState{})

	machine.RegisterState(pb.GameState_SNG_TEXAS_INIT, &sngstate.InitState{})
	machine.RegisterState(pb.GameState_SNG_TEXAS_START, &sngstate.StartState{})
	machine.RegisterState(pb.GameState_SNG_TEXAS_PRE_FLOP, &sngstate.PreflopState{})
	machine.RegisterState(pb.GameState_SNG_TEXAS_FLOP_ROUND, &sngstate.FlopState{})
	machine.RegisterState(pb.GameState_SNG_TEXAS_TURN_ROUND, &sngstate.TurnState{})
	machine.RegisterState(pb.GameState_SNG_TEXAS_RIVER_ROUND, &sngstate.RiverState{})
	machine.RegisterState(pb.GameState_SNG_TEXAS_END, &sngstate.EndState{})

	machine.RegisterState(pb.GameState_Rummy_STAGE_INIT, &rummy.RummyInitstate{})
	machine.RegisterState(pb.GameState_Rummy_STAGE_READY_START, &rummy.RummyReadyState{})
	machine.RegisterState(pb.GameState_Rummy_STAGE_START, &rummy.RummyStartState{})
	machine.RegisterState(pb.GameState_Rummy_STAGE_ZHUANG, &rummy.RummyZhuangState{})
	machine.RegisterState(pb.GameState_Rummy_STAGE_DEAL, &rummy.RummyDealState{})
	machine.RegisterState(pb.GameState_Rummy_STAGE_PLAYING, &rummy.RummyPlayState{})
	machine.RegisterState(pb.GameState_Rummy_STAGE_FIX_CARD, &rummy.RummyFixState{})
	machine.RegisterState(pb.GameState_Rummy_STAGE_SETTLE, &rummy.RummySettleState{})
	machine.RegisterState(pb.GameState_Rummy_STAGE_SHUFFLE, &rummy.RummyShuffleState{})
}
