package test

import (
	"fmt"
	"testing"
	"time"
	"universal/framework/basic/util"
	"universal/framework/timer"
)

func Print() {
	fmt.Println("---->", util.GetNowUnixMilli())
}

func TestTimer(t *testing.T) {
	id := uint64(10)
	tt := timer.NewTimer(timer.INTERVAL)
	// 注册定时器任务
	for i := 0; i < 100000; i++ {
		tt.RegisterTimer(&id, Print, 1000*time.Millisecond, 0, 10)
	}
	tt.StartTest(300)
}
