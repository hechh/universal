package test

import (
	"sync/atomic"
	"testing"
	"time"
	"universal/library/timer"
)

func TestTimer(t *testing.T) {
	var count int64
	printFun := func() {
		atomic.AddInt64(&count, 1)
	}

	tt := timer.NewTimer(4, 7, 4, nil)

	tmps := map[uint64]*uint64{}
	for i := uint64(1); i <= 100000; i++ {
		tmps[i] = &i
		tt.AddTaskFun(tmps[i], printFun, 1*time.Second, 5)
	}

	time.Sleep(6 * time.Second)
	t.Log(atomic.LoadInt64(&count))
}
