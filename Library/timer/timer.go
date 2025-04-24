package timer

import (
	"sync/atomic"
	"time"
	"universal/library/baselib/queue"
	"universal/library/baselib/safe"
	"universal/library/baselib/uerror"
)

type Timer struct {
	tick    int64             // 最小时间间隔
	now     int64             // 当前时间
	tasks   *queue.Queue      // 待插入定时任务队列
	wheels  []*Wheel          // 时间轮
	notify  chan struct{}     // 插入通知
	exit    chan struct{}     // 定时器退出通知
	errorCb func(interface{}) // 错误通知
}

func NewTimer(count, shift, tick int64, cb func(interface{})) *Timer {
	now := time.Now().UnixMilli()
	tt := &Timer{
		tick:    1 << tick,
		now:     now,
		tasks:   queue.NewQueue(),
		wheels:  newWheelList(count, shift, tick, now),
		notify:  make(chan struct{}, 3),
		exit:    make(chan struct{}),
		errorCb: cb,
	}
	safe.SafeGo(cb, tt.run)
	return tt
}

// 添加定时任务
func (d *Timer) AddTaskFun(taskId *uint64, f func(), ttl time.Duration, times int32) error {
	ttl = ttl / time.Millisecond
	return d.AddTask(&Task{
		taskId: taskId,
		task:   f,
		ttl:    int64(ttl),
		expire: time.Now().UnixMilli() + int64(ttl),
		times:  times,
	})
}

func (d *Timer) AddTask(t *Task) error {
	if t.ttl < d.tick {
		return uerror.New(1, -1, "定时器最小时间间隔:%dms", d.tick)
	}
	d.tasks.Push(t)
	// 避免阻塞
	select {
	case d.notify <- struct{}{}:
	default:
	}
	return nil
}

func (d *Timer) run() {
	timer := time.NewTicker(time.Duration(d.tick) * time.Millisecond)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			// 更新时间
			now := atomic.AddInt64(&d.now, d.tick)
			// 处理超时任务
			d.update(now)
			// 处理缓存队列
			d.insert(now)
		case <-d.exit:
			return
		case <-d.notify:
			d.insert(atomic.LoadInt64(&d.now))
		}
	}
}

func (d *Timer) insert(now int64) {
	for item := d.tasks.Pop(); item != nil; item = d.tasks.Pop() {
		// 获取待插入任务
		tt, ok := item.(*Task)
		if !ok || tt == nil {
			continue
		}
		// 插入不成功
		if !d.dispatch(0, tt) {
			d.dispatch(0, tt.Do(now, d.errorCb))
		}
	}
}

func (d *Timer) update(now int64) {
	for i, wheel := range d.wheels {
		// 获取任务
		task := wheel.Get(now)
		if task == nil {
			continue
		}
		// 处理任务
		for item := task; item != nil; {
			tt := item
			item = item.next
			if !d.dispatch(i+1, tt) {
				d.dispatch(0, tt.Do(now, d.errorCb))
			}
		}
	}
}

func (d *Timer) dispatch(pos int, t *Task) bool {
	if t == nil {
		return true
	}

	for ; pos < len(d.wheels); pos++ {
		if d.wheels[pos].Add(t, d.errorCb) {
			return true
		}
	}
	return false
}
