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

func TestMain(m *testing.M) {
	m.Run()
}

func TestTimer01(t *testing.T) {
	for i := 0; i < 50000; i++ {
		timer.Insert(timer.NewTask(func() {
			fmt.Println(util.GetNowUnixSecond(), "------------>", i)
		}, 3*time.Second, false))
	}
	time.Sleep(7 * time.Second)
	timer.Stop()
}

func TestTimer02(t *testing.T) {
	timer.Insert(timer.NewTask(Print, 2*time.Second, false))
	time.Sleep(5 * time.Second)
	timer.Stop()
}

func TestInsert01(t *testing.T) {
	timer.Insert(timer.NewTask(Print, 1*time.Millisecond, true))
	timer.Insert(timer.NewTask(Print, 2*time.Millisecond, true))
	timer.Insert(timer.NewTask(Print, 68*time.Millisecond, true))
	time.Sleep(2 * time.Second)
	timer.Stop()
}

func TestInsert02(t *testing.T) {
	timer.Insert(timer.NewTask(Print, 65*time.Millisecond, true))
	timer.Insert(timer.NewTask(Print, 66*time.Millisecond, true))
	timer.Insert(timer.NewTask(Print, 67*time.Millisecond, true))
	time.Sleep(1 * time.Second)
	timer.Stop()
}
