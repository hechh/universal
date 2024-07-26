package test

import (
	"fmt"
	"testing"
	"time"
	"universal/framework/timer"
	"universal/framework/util"
)

func TestTimer(t *testing.T) {
	now := util.GetNowUnixMilli()
	tt := timer.NewTimer(now)
	ff := func() {
		fmt.Println("--------ff----------")
	}
	tt.Insert(ff, 1*time.Second, false)
	time.Sleep(5 * time.Second)
}

func TestWheel(t *testing.T) {
	wh := timer.NewWheel(util.GetNowUnixMilli(), 6, 13)
	ff := func() {
		fmt.Println("--------ff----------")
	}

	tt := time.NewTicker(100 * time.Millisecond)
	for i := 0; i < 10; i++ {
		<-tt.C
		t.Log(wh.Insert(timer.NewTask(ff, int64(10*time.Millisecond), true)))

		for _, task := range wh.Pop(util.GetNowUnixMilli()) {
			task.Handle()
		}
	}
}
