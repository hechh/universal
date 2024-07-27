package test

import (
	"fmt"
	"testing"
	"time"
	"universal/framework/timer"
	"universal/framework/util"
)

func TestTimer(t *testing.T) {
	tt := timer.NewTimer()
	for i := 0; i < 10000; i++ {
		tt.Insert(func() {
			fmt.Println("-------> ", i)
		}, 1*time.Second, false)
	}
	time.Sleep(5 * time.Second)
	tt.Stop()
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
