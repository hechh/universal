package condition

import (
	"corps/common/cfgData"
	"corps/pb"
)

type BaseCondition struct{}

// 通用条件逻辑
func (d *BaseCondition) Update(data *pb.EntryCondition, cfg *cfgData.EntryCfg, times uint32, subTypes ...uint32) (addTimes uint32) {
	// 判断条件是否服务
	if len(subTypes) > 0 {
		if len(subTypes) != len(cfg.SubType) {
			return
		}
		for i := 0; i < len(subTypes); i++ {
			if subTypes[i] != cfg.SubType[i] {
				return
			}
		}
	}
	// 判断是否达到触发上限
	if data.Times >= cfg.MaxLimit {
		return
	}
	uOldTimes := data.Times
	//取历史值
	if cfg.IsTotal > 0 {
		data.Process = times
		data.Times = 0
	} else {
		// 更新词条进度
		data.Process += times
	}

	if data.Process >= cfg.CondParam {
		data.Times += data.Process / cfg.CondParam
		data.Process = data.Process % cfg.CondParam

		//最大值兼容
		if data.Times >= cfg.MaxLimit {
			data.Times = cfg.MaxLimit
			data.Process = 0
		}

		if data.Times > uOldTimes {
			addTimes = data.Times - uOldTimes
		}
	}

	return
}
