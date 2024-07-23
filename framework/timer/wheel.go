package timer

import (
	"sync/atomic"
)

// 任务轮(用户对定时任务分类、并插入相应任务队列)
type Wheel struct {
	cursor    int64 // 游标
	tickShift int64 // bit数
	tick      int64 // 时间刻度
	size      int64 // 时间轮盘大小
	interval  int64
	buckets   []*TaskBucket // 任务集合
}

func NewWheel(now int64, tick, size int64) *Wheel {
	return &Wheel{
		cursor:    now,
		tickShift: tick,
		tick:      1<<tick - 1,
		size:      1<<size - 1,
		interval:  1<<(tick+size) - 1,
		buckets:   NewTaskBucketList(size),
	}
}

func (d *Wheel) Insert(task *Task) int {
	cursor := atomic.LoadInt64(&d.cursor)
	// 判断是否过期
	if (cursor | d.tick) == (task.expire | d.tick) {
		return -1
	}
	// 判断是否匹配
	if (cursor | d.interval) != (task.expire | d.interval) {
		return 1
	}
	// 插入成功
	d.buckets[(task.expire>>d.tickShift)&d.size].Insert(task)
	return 0
}

func (d *Wheel) Pop(now int64) (list []*Task) {
	cursor := atomic.LoadInt64(&d.cursor)
	if (cursor | d.tick) == (now | d.tick) {
		return
	}
	// 更新时间
	atomic.StoreInt64(&d.cursor, now)
	// 读取过期任务
	end := (now >> d.tickShift) & d.size
	begin := (cursor>>d.tickShift + 1) & d.size
	for i := begin; i <= end; i++ {
		bucket := d.buckets[i&d.size]
		for item := bucket.Pop(); item != nil; item = bucket.Pop() {
			list = append(list, item)
		}
	}
	return
}
