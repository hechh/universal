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
	tt := timer.NewTimer()
	id := new(uint64)
	*id = 10

	tt.Insert(id, 0, Print, 30*time.Millisecond, 10)

	tim := time.NewTicker(10 * time.Millisecond)
	for i := 0; i < 30; i++ {
		<-tim.C
		tt.Update()
	}
}

func TestTT(t *testing.T) {
	id := new(uint64)
	*id = 10
	tt := timer.NewTask(id, 300, Print, uint64(time.Second/(10*time.Millisecond)), 5)
	for i := uint64(1); i <= 10; i++ {
		tt.Handle(i)
	}
}
