package util

import "poker_server/common/pb"

func GetCurStateTTL(machineCfg *pb.MachineConfig, curState pb.GameState) int64 {
	switch curState {
	case pb.GameState_TEXAS_START:
		if machineCfg != nil {
			return machineCfg.StartDuration * 1000 // 准备时间
		} else {
			return 15000 // 15秒准备时间
		}
	case pb.GameState_TEXAS_PRE_FLOP:
		if machineCfg != nil {
			return machineCfg.PreFlopDuration * 1000 // 准备时间
		} else {
			return 70000 // 70秒下注时间
		}
	case pb.GameState_TEXAS_FLOP_ROUND:
		if machineCfg != nil {
			return machineCfg.FlopDuration * 1000 // 准备时间
		} else {
			return 70000 // 70秒下注时间
		}
	case pb.GameState_TEXAS_TURN_ROUND:
		if machineCfg != nil {
			return machineCfg.TurnDuration * 1000 // 准备时间
		} else {
			return 70000
		}
	case pb.GameState_TEXAS_RIVER_ROUND:
		if machineCfg != nil {
			return machineCfg.RiverDuration * 1000 // 准备时间
		} else {
			return 70000
		}
	case pb.GameState_TEXAS_END:
		if machineCfg != nil {
			return machineCfg.EndDuration * 1000 // 准备时间
		} else {
			return 10000 // 10秒结算时间
		}
	case pb.GameState_SNG_TEXAS_END:
		return 60 * 1000
	default:
		if machineCfg != nil {
			return machineCfg.DefaultDuration * 1000 // 准备时间
		} else {
			return 5000 // 默认5秒
		}
	}
}
