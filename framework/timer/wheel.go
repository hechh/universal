package timer

import (
	"sync/atomic"
	"universal/framework/util"
)

// 任务轮(用户对定时任务分类、并插入相应任务队列)
type Wheel struct {
	cursor   int64       // 游标
	bitTick  int64       // tick的bit数
	bitSize  int64       // size的bit数
	tick     int64       // 时间刻度
	size     int64       // 时间轮盘大小
	interval int64       // 间隔时间
	buckets  []*TaskList // 任务集合
}

func NewWheel(now, tick, size int64) *Wheel {
	return &Wheel{
		cursor:   now,
		bitTick:  tick,
		bitSize:  size,
		tick:     1<<tick - 1,
		size:     1<<size - 1,
		interval: 1<<(tick+size) - 1,
		buckets:  NewTaskBucket(1 << size),
	}
}

func (d *Wheel) Insert(task *Task) int {
	cursor := util.GetNowUnixMilli()
	// 判断是否过期
	if (cursor | d.tick) == (task.expire | d.tick) {
		return -1
	}
	// 判断是否匹配
	if (cursor | d.interval) != (task.expire | d.interval) {
		return 1
	}
	// 插入成功
	d.buckets[(task.expire>>d.bitTick)&d.size].Insert(task)
	return 0
}

func (d *Wheel) Pop() *Task {
	now := util.GetNowUnixMilli()
	cursor := atomic.LoadInt64(&d.cursor)
	if (cursor | d.tick) == (now | d.tick) {
		return nil
	}
	// 更新时间
	atomic.StoreInt64(&d.cursor, now)
	// 读取过期任务
	return d.buckets[(now>>d.bitTick)&d.size].Pop()
}
