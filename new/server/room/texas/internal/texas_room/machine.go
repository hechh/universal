package texas_room

import "poker_server/common/pb"

func (d *TexasRoom) GetStateStartTime() int64 {
	return d.machine.GetStateStartTime()
}

func (d *TexasRoom) SetStateStartTime(nowMs int64) {
	d.machine.SetStateStartTime(nowMs)
}

func (d *TexasRoom) GetCurState() pb.TexasGameState {
	return d.machine.GetCurState()
}

func (d *TexasRoom) GetNextState() pb.TexasGameState {
	switch d.GetCurState() {
	case pb.TexasGameState_TGS_INIT:
		return pb.TexasGameState_TGS_START
	case pb.TexasGameState_TGS_START:
		return pb.TexasGameState_TGS_PRE_FLOP
	case pb.TexasGameState_TGS_PRE_FLOP:
		return pb.TexasGameState_TGS_FLOP_ROUND
	case pb.TexasGameState_TGS_FLOP_ROUND:
		return pb.TexasGameState_TGS_TURN_ROUND
	case pb.TexasGameState_TGS_TURN_ROUND:
		return pb.TexasGameState_TGS_RIVER_ROUND
	case pb.TexasGameState_TGS_RIVER_ROUND:
		return pb.TexasGameState_TGS_END
	default:
		return pb.TexasGameState_TGS_INIT
	}
}

func (d *TexasRoom) GetCurStateTTL() int64 {
	switch d.GetCurState() {
	case pb.TexasGameState_TGS_START:
		if d.machineCfg != nil {
			return d.machineCfg.StartDuration * 1000 // 准备时间
		} else {
			return 15000 // 15秒准备时间
		}
	case pb.TexasGameState_TGS_PRE_FLOP:
		if d.machineCfg != nil {
			return d.machineCfg.PreFlopDuration * 1000 // 准备时间
		} else {
			return 70000 // 70秒下注时间
		}
	case pb.TexasGameState_TGS_FLOP_ROUND:
		if d.machineCfg != nil {
			return d.machineCfg.FlopDuration * 1000 // 准备时间
		} else {
			return 70000 // 70秒下注时间
		}
	case pb.TexasGameState_TGS_TURN_ROUND:
		if d.machineCfg != nil {
			return d.machineCfg.TurnDuration * 1000 // 准备时间
		} else {
			return 70000
		}
	case pb.TexasGameState_TGS_RIVER_ROUND:
		if d.machineCfg != nil {
			return d.machineCfg.RiverDuration * 1000 // 准备时间
		} else {
			return 70000
		}
	case pb.TexasGameState_TGS_END:
		if d.machineCfg != nil {
			return d.machineCfg.EndDuration * 1000 // 准备时间
		} else {
			return 10000 // 10秒结算时间
		}
	default:
		if d.machineCfg != nil {
			return d.machineCfg.DefaultDuration * 1000 // 准备时间
		} else {
			return 5000 // 默认5秒
		}
	}
}
