package timer

import (
	"sync/atomic"
	"universal/framework/basic/util"
)

type Wheel struct {
	cursor   int64          // 处理游标
	tickSize int64          // tick的bit数
	tick     int64          // 时间刻度
	interval int64          // 时间间隔
	size     int64          // 转盘大小
	buckets  []atomic.Value // 任务队列
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
	queue := d.buckets[d.GetIndex(task.expire)].Load().(*TaskQueue)
	queue.Insert(task)
	return 0
}

func (d *Wheel) Pop() (rets []*Task) {
	end := d.GetIndex(util.GetNowUnixMilli())
	begin := d.GetIndex(atomic.LoadInt64(&d.cursor))
	if end < begin {
		end += d.size
	}
	// 判断是否存在过期任务队列
	if end == begin {
		return
	}
	// 读取过期
	for ; begin < end; begin++ {
		old := d.buckets[begin&d.size].Load().(*TaskQueue)
		if old.GetCount() <= 0 {
			continue
		}
		d.buckets[begin&d.size].Swap(NewTaskQueue())
		rets = append(rets, old.Remove())
	}
	return
}
