package stat

import (
	"poker_server/framework/actor"
	"poker_server/library/mlog"
	"reflect"
	"time"
)

type Statistics struct {
	actor.Actor
}

func NewStatistics() *Statistics {
	ret := &Statistics{}
	ret.Actor.Register(ret)
	ret.Actor.ParseFunc(reflect.TypeOf(ret))
	ret.Actor.Start()
	actor.Register(ret)
	return ret
}

func (d *Statistics) Analysis(st *CmdStat) {
	time.Sleep(2 * time.Second)
	st.Wait()
	var overs, sum, min, max int64
	for _, ret := range st.players {
		diff := ret.endTime - ret.startTime
		if min == 0 {
			min = diff
		}
		if diff > 200 {
			overs++
		}
		if min > diff {
			min = diff
		}
		if max < diff {
			max = diff
		}
		sum += diff
	}
	avg := sum / int64(st.total)
	mlog.Infof("AVG:%dms, min:%dms, max:%dms, 超时请求:%d, 总请求:%d", avg, min, max, overs, st.total)
}
