package timer

import (
	"sort"
	"sync/atomic"
	"time"
	"universal/library/queue"
	"universal/library/uerror"
	"universal/library/util"
)

type Timer struct {
	now       int64
	startTime int64
	head      *Wheel
	tail      *Wheel
	wheels    []*Wheel
	tasks     *queue.Queue[*Task]
	exit      chan struct{}
	fatal     func(string, ...interface{})
}

func NewTimer(tick int64, fatal func(string, ...interface{})) *Timer {
	now := time.Now().UnixMilli()
	ret := &Timer{
		now:       now,
		startTime: now,
		wheels:    NewWheelList(tick, now),
		tasks:     queue.NewQueue[*Task](),
		exit:      make(chan struct{}),
		fatal:     fatal,
	}
	ret.head = ret.wheels[0]
	ret.tail = ret.wheels[len(ret.wheels)-1]
	util.SafeGo(fatal, ret.run)
	return ret
}

// 注册定时器
func (d *Timer) Register(taskId *uint64, f func(), ttl time.Duration, times int32) error {
	tt := int64(ttl / time.Millisecond)
	if tt>>d.head.shift <= 0 {
		return uerror.N(1, -1, "小于定时器最小时间间隔:%dms", 1<<d.head.shift)
	}
	if (tt >> d.tail.shift) > d.tail.mask {
		return uerror.N(1, -1, "定时器超出最大时间范围:%dms", d.tail.shift)
	}
	d.tasks.Push(&Task{
		taskId: taskId,
		task:   f,
		ttl:    tt,
		times:  times,
	})
	return nil
}

func (d *Timer) run() {
	tick := int64(1) << d.head.shift
	tt := time.NewTicker(time.Duration(tick) * time.Millisecond)
	defer tt.Stop()

	for {
		select {
		case <-tt.C:
			now := atomic.AddInt64(&d.now, tick)
			d.update(now)
			d.flush(now)
		case <-d.exit:
			return
		}
	}
}

// 刷入新定时器
func (d *Timer) flush(now int64) {
	news := []*Task{}
	for tt := d.tasks.Pop(); tt != nil; tt = d.tasks.Pop() {
		tt.expire = now + tt.ttl
		news = append(news, tt)
	}

	sort.Slice(news, func(i, j int) bool {
		return news[i].expire < news[j].expire
	})

	pos := 0
	lnews := len(news)
	for _, w := range d.wheels {
		for ; pos < lnews && w.IsMatch(news[pos]); pos++ {
			w.Insert(news[pos])
		}
		if lnews <= pos {
			break
		}
	}
}

func (d *Timer) update(now int64) {
	news := []*Task{}
	for _, w := range d.wheels {
		for tts := w.Get(now); tts != nil; {
			tt := tts
			tts = tts.next
			tt.next = nil
			if !d.head.IsExpire(tt) {
				news = append(news, tt)
			} else {
				if newTask := tt.handle(now, d.fatal); newTask != nil {
					news = append(news, newTask)
				}
			}
		}
		if !w.IsCarry() {
			break
		}
	}

	sort.Slice(news, func(i, j int) bool {
		return news[i].expire < news[j].expire
	})

	pos := 0
	lnews := len(news)
	for _, w := range d.wheels {
		for ; pos < lnews && w.IsMatch(news[pos]); pos++ {
			w.Insert(news[pos])
		}
		if lnews <= pos {
			break
		}
	}
}
