package stat

import (
	"sync"
	"sync/atomic"
	"time"
)

type Result struct {
	uid       uint64
	cmd       uint32
	startTime int64
	endTime   int64
}

type CmdStat struct {
	sync.WaitGroup
	total   int32
	players map[uint64]*Result
}

func NewCmdStat(cmd uint32, uids ...uint64) *CmdStat {
	players := make(map[uint64]*Result)
	for _, uid := range uids {
		players[uid] = &Result{uid: uid, cmd: cmd}
	}
	return &CmdStat{players: players}
}

func (r *Result) Start() {
	atomic.StoreInt64(&r.startTime, time.Now().UnixMilli())
}

func (r *Result) Finish() {
	atomic.StoreInt64(&r.endTime, time.Now().UnixMilli())
}

func (r *Result) IsFinish() bool {
	return atomic.LoadInt64(&r.endTime) > 0
}

func (r *CmdStat) Get(uid uint64) *Result {
	return r.players[uid]
}

func (r *CmdStat) Done() {
	atomic.AddInt32(&r.total, 1)
	r.WaitGroup.Done()
}
