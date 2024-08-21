package test

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
	"universal/framework/basic/util"
	"universal/framework/timer"
)

var (
	count int64
)

func Print() {
	fmt.Println(util.GetNowUnixSecond(), "----", atomic.AddInt64(&count, 1))
}

func TestTimer01(t *testing.T) {
	tt := timer.NewTimer()
	tt.Insert(timer.NewTask(Print, 2*time.Second, false))
	tt.Insert(timer.NewTask(Print, 9*time.Second, false))
	tt.Insert(timer.NewTask(Print, 60*time.Minute, false))
	tt.Insert(timer.NewTask(Print, 49*time.Hour, false))
	time.Sleep(6 * time.Second)
	tt.Stop()
}

func TestTimer02(t *testing.T) {
	tt := timer.NewTimer()
	tt.Insert(timer.NewTask(Print, 3*time.Second, false))
	time.Sleep(10 * time.Second)
	t.Log(util.GetNowUnixMilli() / 1000)
	tt.Stop()
}

func TestInsert01(t *testing.T) {
	tt := timer.NewTimer()
	tt.Insert(timer.NewTask(Print, 1*time.Millisecond, true))
	tt.Insert(timer.NewTask(Print, 2*time.Millisecond, true))
	tt.Insert(timer.NewTask(Print, 68*time.Millisecond, true))
	time.Sleep(2 * time.Second)
	tt.Stop()
}

func TestInsert02(t *testing.T) {
	tt := timer.NewTimer()
	tt.Insert(timer.NewTask(Print, 65*time.Millisecond, true))
	tt.Insert(timer.NewTask(Print, 66*time.Millisecond, true))
	tt.Insert(timer.NewTask(Print, 67*time.Millisecond, true))
	time.Sleep(1 * time.Second)
	tt.Stop()
}
