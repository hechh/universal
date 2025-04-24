package test

import (
	"sync/atomic"
	"testing"
	"time"
	"universal/library/mlog"
	"universal/library/timer"
)

var count int64

func Print() {
	atomic.AddInt64(&count, 1)
}

func TestTimer(t *testing.T) {

	tt := timer.NewTimer(4, 7, 4, mlog.Fatal)

	tmps := map[uint64]*uint64{}
	for i := uint64(1); i <= 100000; i++ {
		tmps[i] = &i
		tt.AddTaskFun(tmps[i], Print, 1*time.Second, 5)
	}

	time.Sleep(6 * time.Second)
	t.Log(atomic.LoadInt64(&count))
}
