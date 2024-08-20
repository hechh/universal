package timer

import (
	"runtime/debug"
	"sync"
	"time"
	"universal/framework/basic/async"
	"universal/framework/basic/util"
	"universal/framework/plog"
)

type Timer struct {
	sync.WaitGroup
	wheels     [4]*Wheel     // 定时任务转盘
	over       *async.Queue  // 过期任务
	overNotify chan struct{} // 过期通知
	exitNotify chan struct{} // 退出通知
}

func NewTimer() *Timer {
	ret := &Timer{
		wheels:     [4]*Wheel{NewWheel(6, 7), NewWheel(13, 7), NewWheel(20, 7), NewWheel(27, 7)},
		over:       async.NewQueue(),
		overNotify: make(chan struct{}, 1),
		exitNotify: make(chan struct{}, 2),
	}
	util.SafeGo(func(err interface{}) {
		plog.Fatal("error: %v, stack: %s", err, string(debug.Stack()))
	}, ret.expireHandle)
	util.SafeGo(func(err interface{}) {
		plog.Fatal("error: %v, stack: %s", err, string(debug.Stack()))
	}, ret.overHandle)
	return ret
}

func (d *Timer) Insert(tt *Task) {
	for _, wheel := range d.wheels {
		switch wheel.Insert(tt) {
		case -1:
			d.over.Push(tt)
			select {
			case d.overNotify <- struct{}{}:
			default:
			}
			return
		case 0:
			return
		}
	}
}

// 过期任务执行
func (d *Timer) expireHandle() {
	tt := time.NewTicker(time.Duration(d.wheels[0].tick))
	for {
		select {
		case <-tt.C:
			for _, wheel := range d.wheels {
				for _, tt := range wheel.Pop() {
					d.over.Push(tt)
					select {
					case d.overNotify <- struct{}{}:
					default:
					}
				}
			}
		case <-d.exitNotify:
			return
		}
	}
}

// 过期任务执行
func (d *Timer) overHandle() {
	for {
		select {
		case <-d.overNotify:
			for tt := d.over.Pop(); tt != nil; tt = d.over.Pop() {
				task, ok := tt.(*Task)
				if !ok || task == nil {
					continue
				}
				// 执行定时任务
				for _, tt := range task.Handle() {
					d.Insert(tt)
				}
			}
		case <-d.exitNotify:
			return
		}
	}
}
