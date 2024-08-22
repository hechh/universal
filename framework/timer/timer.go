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
	updateOver chan struct{} // 过期通知
	exitRun    chan struct{} // 退出通知
	exitHandle chan struct{} // 退出通知
}

func NewTimer() *Timer {
	ret := &Timer{
		wheels:     [4]*Wheel{NewWheel(1, 6, 7), NewWheel(2, 13, 7), NewWheel(3, 20, 7), NewWheel(4, 27, 7)},
		over:       async.NewQueue(),
		updateOver: make(chan struct{}, 1),
		exitRun:    make(chan struct{}, 1),
		exitHandle: make(chan struct{}, 1),
	}
	// 定时触发
	util.SafeGo(func(err interface{}) {
		plog.Fatal("%v\n stack: %s", err, string(debug.Stack()))
	}, ret.run)
	// 执行定时任务
	ret.Add(1)
	util.SafeGo(func(err interface{}) {
		plog.Fatal("%v\n stack: %s", err, string(debug.Stack()))
	}, ret.handle)
	return ret
}

func (d *Timer) Stop() {
	d.exitRun <- struct{}{}
	d.exitHandle <- struct{}{}
	d.Wait()
}

func (d *Timer) addOver(tt ...*Task) {
	d.over.Push(tt)
	select {
	case d.updateOver <- struct{}{}:
	default:
	}
}

func (d *Timer) Insert(tt *Task) {
	for _, wheel := range d.wheels {
		switch wheel.Insert(tt) {
		case -1:
			d.addOver(tt)
			return
		case 0:
			return
		}
	}
}

// 过滤过期任务执行
func (d *Timer) run() {
	tt := time.NewTicker(time.Duration(d.wheels[0].tick) * time.Millisecond)
	for {
		select {
		case <-d.exitRun:
			return
		case <-tt.C:
			for _, wheel := range d.wheels {
				now := util.GetNowUnixMilli()
				for _, tt := range wheel.Pop(now) {
					if wheel.IsExpired(now, tt.expire) {
						if !tt.once {
							tt.refresh(now)
							d.Insert(tt)
						}
						d.addOver(tt)
					} else {
						d.Insert(tt)
					}
				}
			}
		}
	}
}

// 过期任务执行
func (d *Timer) handle() {
	handleF := func() {
		for tt := d.over.Pop(); tt != nil; tt = d.over.Pop() {
			switch tmp := tt.(type) {
			case *Task:
				tmp.handle()
			case []*Task:
				for _, vv := range tmp {
					vv.handle()
				}
			}
		}
	}
	defer func() {
		handleF()
		d.Done()
	}()
	for {
		select {
		case <-d.updateOver:
			handleF()
		case <-d.exitHandle:
			return
		}
	}
}
