package timer

import (
	"runtime/debug"
	"time"
	"universal/framework/basic/async"
	"universal/framework/basic/util"
	"universal/framework/plog"
)

type Timer struct {
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
		exitNotify: make(chan struct{}, 1),
	}
	util.SafeGo(func(err interface{}) {
		plog.Fatal("%v\n stack: %s", err, string(debug.Stack()))
	}, ret.run)
	return ret
}

func (d *Timer) Stop() {
	d.exitNotify <- struct{}{}
}

func (d *Timer) addOver(ts ...*Task) {
	for _, tt := range ts {
		d.over.Push(tt)

		// 通知
		select {
		case d.overNotify <- struct{}{}:
		default:
		}
	}
}

func (d *Timer) Insert(tt *Task) {
	for i, wheel := range d.wheels {
		ret := wheel.Insert(tt)
		switch ret {
		case -1:
			d.addOver(tt)
			plog.Trace("Insert %d --> %d over task.expire: %d, task.once: %v", i, ret, tt.expire, tt.once)
			return
		case 0:
			plog.Trace("Insert %d --> %d task.expire: %d, task.once: %v", i, ret, tt.expire, tt.once)
			return
		}
	}
}

// 过滤过期任务执行
func (d *Timer) run() {
	tt := time.NewTicker(time.Duration(d.wheels[0].tick) * time.Millisecond)
	for {
		select {
		case <-tt.C:
			for _, wheel := range d.wheels {
				d.addOver(wheel.Pop()...)
			}
		case <-d.exitNotify:
			return
		}
	}
}

/*
// 过期任务执行
func (d *Timer) overHandle() {
	for wheel := d.wheels[0]; ; {
		select {
		case <-d.overNotify:
			for tt := d.over.Pop(); tt != nil; tt = d.over.Pop() {
				handler(wheel, tt.(*Task), d.Insert)
			}
		case <-d.exitNotify:
			return
		}
	}
}
*/

func handler(wheel *Wheel, tt *Task, f func(tt *Task)) {
	for tt != nil {
		// 遍历单向链表
		tmp := tt
		tt = tt.next
		tmp.next = nil
		// 执行任务
		now := util.GetNowUnixMilli()
		plog.Trace("task ---> now: %d, expire: %d, once: %v", now, tmp.expire, tmp.once)
		if now >= tmp.expire || wheel.GetIndex(now) == wheel.GetIndex(tmp.expire) {
			if !tmp.once {
				f(NewTask(tmp.handle, tmp.ttl, tmp.once))
			}
			tmp.handle()
		} else {
			f(tmp)
		}
	}
}
