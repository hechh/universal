package test

import (
	"fmt"
	"hego/framework/basic"
	"hego/framework/timer"
	"testing"
	"time"
)

func Print() {
	fmt.Println("---->", basic.GetNowUnixSecond())
}

func TestTimer(t *testing.T) {
	id := uint64(10)
	tt := timer.NewTimer(timer.INTERVAL)
	// 注册定时器任务
	//for i := 0; i < 1000; i++ {
	tt.RegisterTimer(&id, Print, 1000*time.Millisecond, 0, 10)
	//}
	tt.StartTest(600)
}
