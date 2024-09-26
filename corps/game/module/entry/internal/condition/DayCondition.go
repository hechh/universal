package condition

import (
	"corps/base"
	"corps/common/cfgData"
	"corps/pb"
)

type DayCondition struct{}

// 通用条件逻辑
func (d *DayCondition) Update(data *pb.EntryCondition, cfg *cfgData.EntryCfg, times uint32, subTypes ...uint32) (flag uint32) {
	now := base.GetNow()
	// 判断是否跨天
	// 同一天内，多次登录，不计入条件触发
	if base.IsSameDay(now, data.UpdateTime) {
		return
	}
	// 判断是否达到触发上限
	if data.Times >= cfg.MaxLimit {
		return
	}
	// 更新词条进度
	data.UpdateTime = now
	data.Process += 1
	if data.Process < cfg.CondParam {
		return
	}
	// 累计属性
	data.Process -= cfg.CondParam
	data.Times++
	return 1
}
