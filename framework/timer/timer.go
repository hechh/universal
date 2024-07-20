package timer

import (
	"sync"
	"sync/atomic"
	"time"
	"universal/framework/util"
	"unsafe"
)

const (
	MILLISECOND = int64(10 * time.Millisecond)
)

type Task struct {
	task   func() // 定时任务
	isOnce bool   // 是否一次性定时器
	expire int64  // 过期时间
	ttl    int64  // 定时时长
	next   *Task
}

type TaskList struct {
	head  *Task
	tail  *Task
	count int64
}

func NewTaskList() *TaskList {
	node := new(Task)
	return &TaskList{head: node, tail: node}
}

func (d *TaskList) Insert(task *Task) {
	prevNode := (*Task)(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&d.tail)), unsafe.Pointer(task)))
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&prevNode.next)), unsafe.Pointer(task))
	atomic.AddInt64(&d.count, 1)
}

func (d *TaskList) Pop() (tt *Task) {
	if tt = (*Task)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&d.head.next)))); tt != nil {
		d.head.next = nil
		d.head = tt
		atomic.AddInt64(&d.count, -1)
	}
	return
}

// 任务轮(用户对定时任务分类、并插入相应任务队列)
type Wheel struct {
	cursor    int64       // 游标
	tick      int64       // 时间刻度
	wheelSize int64       // 时间轮盘大小
	interval  int64       // 时间间隔
	buckets   []*TaskList // 任务集合
}

func NewWheel(tick, size int64) *Wheel {
	buckets := make([]*TaskList, size)
	for i := int64(0); i < size; i++ {
		buckets[i] = NewTaskList()
	}
	return &Wheel{
		tick:      tick,
		wheelSize: size,
		interval:  tick * size,
		buckets:   buckets,
	}
}

func timeToIndex(t int64, tick, size int64) int64 {
	return (t / tick) % size
}

func (d *Wheel) Insert(task *Task) (flag bool) {
	now := util.GetNowUnixNano() / MILLISECOND
	if flag = task.expire-now <= d.interval; flag {
		d.buckets[timeToIndex(task.expire, d.tick, d.wheelSize)].Insert(task)
	}
	return
}

func (d *Wheel) Pop(now int64) (list []*Task) {
	index := timeToIndex(now, d.tick, d.wheelSize)
	if index != d.cursor {
		d.cursor = index
		bucket := d.buckets[d.cursor]
		for item := bucket.Pop(); item != nil; item = bucket.Pop() {
			list = append(list, item)
		}
	}
	return
}

type Timer struct {
	sync.WaitGroup
	// 毫秒 --->>> tick：10ms, wheelSize: 100, interval: 1000ms=1s
	// 秒 --->>> tick: 1s, wheelSize: 120, interval: 120s=2min
	// 分 --->>> tick: 2m, wheelSize: 120, interval: 240m=4hour
	// 时 --->>> tick: 4h, wheelSize: 120, interval: 480h=20day
	wheels        [4]*Wheel     // 定时任务队列
	overTimeTasks *TaskList     // 超时任务队列
	overTimeCh    chan struct{} // 超时任务处理通知
	exitCh        chan struct{} // 退出定时器
}

func NewTimer() *Timer {
	ret := &Timer{
		wheels: [4]*Wheel{
			NewWheel(10, 100),
			NewWheel(1, 120),
			NewWheel(2, 120),
			NewWheel(4, 120),
		},
		overTimeTasks: NewTaskList(),
		overTimeCh:    make(chan struct{}, 2),
		exitCh:        make(chan struct{}, 1),
	}
	go ret.handle() // 处理超时任务
	go ret.run()    // 定时处理任务
	return ret
}

func (d *Timer) addTask(task *Task) {
	for _, wheel := range d.wheels {
		if wheel.Insert(task) {
			return
		}
	}
}

func (d *Timer) Insert(f func(), ttl time.Duration, isOnce bool) {
	d.addTask(&Task{
		task:   f,
		isOnce: isOnce,
		ttl:    int64(ttl),
		expire: util.GetNowUnixNano() / MILLISECOND,
	})
	// 定时的ttl超过20天，需要特殊处理。
	// todo
}

func (d *Timer) run() {
	tt := time.NewTicker(time.Duration(MILLISECOND))
	<-tt.C
	for {
		<-tt.C
		now := util.GetNowUnixNano() / MILLISECOND
		// 执行触发的
		for _, task := range d.wheels[0].Pop(now) {
			d.overTimeTasks.Insert(task)
		}
		// 转移任务队列
		for i := 1; i <= 3; i++ {
			wheel := d.wheels[i]
			for _, task := range wheel.Pop(now) {
				d.addTask(task)
			}
		}
	}
}

// 处理超时任务
func (d *Timer) handle() {
	d.Add(1)
	defer func() {
		for task := d.overTimeTasks.Pop(); task != nil; task = d.overTimeTasks.Pop() {
			task.task()
		}
		d.Done()
	}()

	for {
		select {
		case <-d.overTimeCh:
			for task := d.overTimeTasks.Pop(); task != nil; task = d.overTimeTasks.Pop() {
				// 执行定时任务
				task.task()

				// 循环定时任务，再一次插入任务队列
				if !task.isOnce {
					d.Insert(task.task, time.Duration(task.ttl), task.isOnce)
				}
			}
		case <-d.exitCh:
			return
		}
	}
}
