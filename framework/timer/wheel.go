package timer

import (
	"sync/atomic"
	"time"
	"universal/framework/basic/async"
	"universal/framework/basic/util"
	"universal/framework/plog"
)

type Task struct {
	once   bool          // 是否一次性
	ttl    time.Duration // 定时时长
	expire int64         // 过期时长
	handle func()        // 任务
	next   *Task         // 下一个任务
}

type Wheel struct {
	cursor   int64          // 处理游标
	tickSize int64          // tick的bit数
	tick     int64          // 时间刻度
	interval int64          // 时间间隔
	size     int64          // 转盘大小
	buckets  []*async.Queue // 任务队列
}

func NewTask(f func(), ttl time.Duration, once bool) *Task {
	return &Task{
		handle: f,
		once:   once,
		ttl:    ttl,
		expire: util.GetNowTime().Add(ttl).UnixMilli(),
	}
}

func NewWheel(tick, size int64) *Wheel {
	return &Wheel{
		cursor:   util.GetNowUnixMilli(),
		tickSize: tick,
		tick:     1<<tick - 1,
		interval: 1<<(tick+size) - 1,
		size:     1<<size - 1,
		buckets:  NewTaskPool(1 << size),
	}
}

func NewTaskPool(size int) (rets []*async.Queue) {
	rets = make([]*async.Queue, size)
	for i := 0; i < size; i++ {
		rets[i] = async.NewQueue()
	}
	return
}

func (d *Wheel) GetIndex(val int64) int64 {
	return (val >> d.tickSize) & d.size
}

func (d *Wheel) Insert(task *Task) int {
	cursor := util.GetNowUnixMilli()
	// 判断是否匹配
	if (cursor | d.interval) != (task.expire | d.interval) {
		return 1
	}
	// 判断是否过期
	if (cursor | d.tick) == (task.expire | d.tick) {
		return -1
	}
	// 插入成功
	plog.Trace("Wheel.Insert tick: %fs, interval: %fs, cursor: %d, task: %d, expire: %d", float32(d.tick)/1000, float32(d.interval)/1000, d.GetIndex(cursor), d.GetIndex(task.expire), task.expire)
	d.buckets[d.GetIndex(task.expire)].Push(task)
	return 0
}

func (d *Wheel) Pop() (rets []*Task) {
	now := util.GetNowUnixMilli()
	cursor := atomic.SwapInt64(&d.cursor, now)
	begin, end := d.GetIndex(cursor), d.GetIndex(now)
	if end < begin {
		end += d.size
	}
	defer plog.Trace("Wheel.Pop tick: %fs, interval: %fs, begin: %d, end: %d, result: %d", float32(d.tick)/1000, float32(d.interval)/1000, begin, end, len(rets))
	// 读取过期
	for i := begin; i < end; i++ {
		old := d.buckets[i&d.size]
		if old.GetCount() <= 0 {
			continue
		}
		for tt := old.Pop(); tt != nil; tt = old.Pop() {
			rets = append(rets, tt.(*Task))
		}
	}
	return
}
