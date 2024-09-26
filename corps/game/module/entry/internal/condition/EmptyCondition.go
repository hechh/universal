package condition

import (
	"corps/common/cfgData"
	"corps/pb"
)

type EmptyCondition struct{}

func (d *EmptyCondition) Update(data *pb.EntryCondition, cfg *cfgData.EntryCfg, times uint32, params ...uint32) uint32 {
	if data.Times >= cfg.MaxLimit {
		return 0
	}
	if cfg.IsTotal > 0 {
		data.Times = times
	} else {
		data.Times++
	}

	return 1
}
