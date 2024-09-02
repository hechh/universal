package timer

import (
	"runtime/debug"
	"sync"
	"time"
	"universal/framework/basic/async"
	"universal/framework/basic/util"
	"universal/framework/plog"
)

var (
	object *Timer = NewTimer()
)

func Insert(tt *Task) {
	object.Insert(tt)
}

func Stop() {
	object.Stop()
}

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
		wheels:     [4]*Wheel{NewWheel(6, 7), NewWheel(13, 7), NewWheel(20, 7), NewWheel(27, 7)},
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
			overs := []*Task{}
			now := util.GetNowUnixMilli()
			for _, wl := range d.wheels {
				for _, item := range wl.Pop(now) {
					if !d.wheels[0].IsExpired(now, item.expire) {
						d.Insert(item)
					} else {
						if !item.once {
							item.Update(now)
							d.Insert(item)
						}
						overs = append(overs, item)
					}
				}
			}
			d.addOver(overs...)
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
