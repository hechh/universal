package timer

import (
	"time"
	"universal/library/async"
	"universal/library/mlog"
	"universal/library/uerror"
)

type Timer struct {
	tick   int64               // 最小时间间隔
	now    int64               // 当前时间
	wheels []*Wheel            // 时间轮
	tasks  *async.Queue[*Task] // 待插入定时任务队列
	notify chan struct{}       // 插入通知
	exit   chan struct{}       // 定时器退出通知
}

func NewTimer(count, shift, tick int64) *Timer {
	now := time.Now().UnixMilli()
	ret := &Timer{
		tick:   1 << tick,
		now:    now,
		wheels: NewWheels(count, shift, tick, now),
		notify: make(chan struct{}, 1),
		exit:   make(chan struct{}),
	}
	async.SafeGo(mlog.Fatalf, ret.run)
	return ret
}

// 注册定时器
func (d *Timer) Register(taskId *uint64, f func(), ttl time.Duration, times int32) error {
	tt := int64(ttl / time.Millisecond)
	if tt < d.tick {
		return uerror.New(1, -1, "定时器最小时间间隔:%dms", d.tick)
	}

	max := d.wheels[0]
	if (tt >> max.shift) > max.mask {
		return uerror.New(1, -1, "定时器超出最大时间范围:%dms", max.shift)
	}

	d.add(&Task{
		taskId: taskId,
		task:   f,
		ttl:    int64(ttl),
		expire: time.Now().UnixMilli() + int64(ttl),
		times:  times,
	})
	return nil
}

func (d *Timer) add(task *Task) {
	d.tasks.Push(task)

	select {
	case d.notify <- struct{}{}:
	default:
	}
}

func (d *Timer) run() {
}
