package timer

import (
	"time"
	"universal/framework/common/async"
)

const (
	SECOND = 1 * int64(time.Second)
	MINUTE = 60 * SECOND
	HOUR   = 60 * MINUTE
	DAY    = 24 * HOUR
)

func truncate(x, m int64) int64 {
	if m <= 0 {
		return x
	}
	return x - x%m
}

func timeToMs(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func msToTime(t int64) time.Time {
	return time.Unix(0, t*int64(time.Millisecond)).UTC()
}

func timeToIndex(t int64, tick, size int64) int64 {
	return (t / tick) % size
}

// 定时任务
type Task struct {
	expire int64  // 过期时间
	ttl    int64  // 定时时长
	handle func() // 定时任务
	isOnce bool   // 是否一次性定时器
}

type Bucket struct {
	*async.Queue
	expire int64
	uuid   int64
}

// 任务轮(用户对定时任务分类、并插入相应任务队列)
type Wheel struct {
	tick      int64     // 时间刻度
	wheelSize int64     // 时间轮盘大小
	interval  int64     // 总时长
	buckets   []*Bucket // 任务集合
}

type Timer struct {
	wheels   map[int32]*Wheel // 转盘
	handles  *async.Queue     // 超时任务
	handleCh chan struct{}    // 触发
	adds     *async.Queue
	addCh    chan struct{} // 循环定时触发任务
}

func (d *Timer) RegisterTimer(f func(), ttl int64, isOnce bool) {
	d.RegisterTask(&Task{
		isOnce: isOnce,
		handle: f,
		ttl:    ttl,
		expire: time.Now().UnixMilli(),
	})
}

func (d *Timer) RegisterTask(vv *Task) {
	// 获取当前cursor
	wheel := d.wheels[getWheelType(vv.ttl)]
	cursor := timeToIndex(time.Now().UnixMilli(), wheel.tick, wheel.wheelSize)
	index := timeToIndex(vv.expire, wheel.tick, wheel.wheelSize)

	// 任务已经超时，立即执行
	if cursor != index {
		d.Handle(vv)
		return
	}

	// 任务尚未超时，放入wheel中等待触发
	bucket := wheel.buckets[index]
	bucket.Push(vv)
}

func getWheelType(ttl int64) int32 {
	if ttl <= SECOND {
		return 1
	} else if ttl <= MINUTE {
		return 2
	} else if ttl <= HOUR {
		return 3
	} else if ttl <= DAY {
		return 4
	}
	return 5
}

func (d *Timer) Handle(task *Task) {
	d.handles.Push(task)

	select {
	case d.handleCh <- struct{}{}:
	default:
	}
}

func (d *Timer) run() {
	for {
		select {
		case <-d.handleCh:
			for item := d.handles.Pop(); item != nil; item = d.handles.Pop() {
				vv := item.(*Task)
				vv.handle()
				if !vv.isOnce {
					d.addCh <- vv
				}
			}
		}
	}
}
