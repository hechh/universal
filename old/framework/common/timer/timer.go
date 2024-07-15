package timer

import (
	"sync/atomic"
	"time"
	"unsafe"
)

type Task struct {
	task   func() // 定时任务
	isOnce bool   // 是否一次性定时器
	expire int64  // 过期时间
	ttl    int64  // 定时时长
	next   *Task
}

type TaskList struct {
	head *Task
	tail *Task
}

func (d *TaskList) Insert(task *Task) {
	prevNode := (*Task)(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&d.tail)), unsafe.Pointer(task)))
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&prevNode.next)), unsafe.Pointer(task))
}

func (d *TaskList) Pop() (tt *Task) {
	if tt = (*Task)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&d.head.next)))); tt != nil {
		d.head.next = nil
		d.head = tt
	}
	return
}

// 任务轮(用户对定时任务分类、并插入相应任务队列)
type Wheel struct {
	cursor   int64       // 游标
	tick     int64       // 时间刻度
	size     int64       // 时间轮盘大小
	interval int64       // 时间间隔
	buckets  []*TaskList // 任务集合
}

func timeToIndex(t int64, tick, size int64) int64 {
	return (t / tick) % size
}

func (d *Wheel) Insert(now int64, task *Task) (flag bool) {
	if flag = (task.expire - now) <= d.interval; flag {
		d.buckets[timeToIndex(task.expire, d.tick, d.size)].Insert(task)
	}
	return
}

func (d *Wheel) Pop(now int64) (list []*Task) {
	if index := timeToIndex(now, d.tick, d.size); index != d.cursor {
		d.cursor = index
		bucket := d.buckets[d.cursor]
		for item := bucket.Pop(); item != nil; item = bucket.Pop() {
			list = append(list, item)
		}
	}
	return
}

type Timer struct {
	wheels    []*Wheel      // 转盘
	addList   *TaskList     // 添加任务
	runList   *TaskList     // 超时任务
	addNotify chan struct{} // 通知
	runNotify chan struct{} // 触发
}

func (d *Timer) RegisterTimer(f func(), ttl int64, isOnce bool) {
	d.addTask(&Task{
		task:   f,
		ttl:    ttl,
		isOnce: isOnce,
		expire: time.Now().Add(time.Duration(ttl)).UnixMilli(),
	})
}

func (d *Timer) insert(now int64, item *Task) {
	for _, wheel := range d.wheels {
		if wheel.Insert(now, item) {
			return
		}
	}
}

func (d *Timer) addTask(tt *Task) {
	d.addList.Insert(tt)

	select {
	case d.addNotify <- struct{}{}:
	default:
	}
}

func (d *Timer) run() {
	for tt := time.NewTicker(time.Duration(d.wheels[0].tick)); ; {
		select {
		case <-d.addNotify:
			for item := d.addList.Pop(); item != nil; item = d.addList.Pop() {
				now := time.Now().UnixMilli()
				if now >= item.expire {
					d.addRun(item)
				} else {
					d.insert(now, item)
				}
			}
		case <-tt.C:
			for _, wheel := range d.wheels {
				now := time.Now().UnixMilli()
				for _, item := range wheel.Pop(now) {
					if now >= item.expire {
						d.addRun(item)
					} else {
						d.addTask(item)
					}
				}
			}
		}
	}
}

func (d *Timer) addRun(vv *Task) {
	d.runList.Insert(vv)

	// 通知处理
	select {
	case d.runNotify <- struct{}{}:
	default:
	}
}

// 立即执行定时任务
func (d *Timer) handle() {
	for {
		select {
		case <-d.runNotify:
			for item := d.runList.Pop(); item != nil; item = d.runList.Pop() {
				item.task()
				// 循环定时任务
				if !item.isOnce {
					item.expire = time.Now().Add(time.Duration(item.ttl)).UnixMilli()
					d.addTask(item)
				}
			}
		}
	}
}
