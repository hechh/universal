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
	tt    = timer.NewTimer()
)

func Print() {
	fmt.Println(util.GetNowUnixSecond(), "----", atomic.AddInt64(&count, 1))
}

func TestMain(m *testing.M) {
	m.Run()
}

func BenchmarkTimer(b *testing.B) {
	print := func() {
		fmt.Println(util.GetNowUnixSecond(), "----", atomic.AddInt64(&count, 1))
	}
	for i := 0; i < b.N; i++ {
		tt.Insert(timer.NewTask(print, 1*time.Second, true))
	}
}

func TestTimer01(t *testing.T) {
	tt := timer.NewTimer()
	for i := 0; i < 5000; i++ {
		tt.Insert(timer.NewTask(func() {
			fmt.Println(util.GetNowUnixSecond(), "------------>", i)
		}, 1*time.Second, false))
	}
	time.Sleep(6 * time.Second)
	tt.Stop()
}

func TestTimer02(t *testing.T) {
	tt := timer.NewTimer()
	tt.Insert(timer.NewTask(Print, 2*time.Second, false))
	time.Sleep(5 * time.Second)
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
