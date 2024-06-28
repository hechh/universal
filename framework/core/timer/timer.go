package timer

import (
	"sync/atomic"
	"time"
	"universal/framework/common/async"
	"unsafe"
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
	next   *Task  // 任务链表
}

type Bucket struct {
	head   *Task
	tail   *Task
	expire int64
}

func (d *Bucket) Push(task *Task) {
	// 将新增节点插入链表
	prevNode := (*Task)(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&d.tail)), unsafe.Pointer(task)))
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&prevNode.next)), unsafe.Pointer(task))
}

func (d *Bucket) Remove() *Task {
	if task := (*Task)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&d.head.next)))); task == nil {
		d.head.next = nil
		d.tail = d.head
		return task
	}
	return nil
}

// 任务轮(用户对定时任务分类、并插入相应任务队列)
type Wheel struct {
	cursor    int64     // 时间指针
	tick      int64     // 时间刻度
	wheelSize int64     // 时间轮盘大小
	interval  int64     // 总时长
	buckets   []*Bucket // 任务集合
}

// 获取当前游标
func (d *Wheel) getCursor() int64 {
	return atomic.LoadInt64(&d.cursor)
}

// 更新游标
func (d *Wheel) Next(ttl int64) bool {
	cur := d.getCursor()
	return !atomic.CompareAndSwapInt64(&d.cursor, cur, ttl/d.tick)
}

// 获取任务
func (d *Wheel) GetTask() *Task {
	bucket := d.buckets[d.getCursor()%d.wheelSize]
	task := bucket.Remove()
	atomic.StoreInt64(&bucket.expire, 0)
	return task
}

type Timer struct {
	timestamp int64         // 开启时间
	wheels    []*Wheel      // 转盘
	handles   *async.Queue  // 超时任务
	handleCh  chan struct{} // 触发
}

func (d *Timer) RegisterTimer(f func(), ttl int64, isOnce bool) {
	for _, wheel := range d.wheels {
		// 过滤不合适地wheel
		if wheel.interval < ttl {
			continue
		}

		// 定义task
		now := time.Now()
		task := &Task{
			handle: f,
			ttl:    ttl,
			isOnce: isOnce,
			expire: now.Add(time.Duration(ttl)).UnixMilli(),
		}

		// 获取当前cursor
		expire := truncate(now.UnixMilli(), wheel.tick)
		cursor := timeToIndex(expire, wheel.tick, wheel.wheelSize)
		index := timeToIndex(task.expire, wheel.tick, wheel.wheelSize)

		// 任务已经超时，立即执行
		if cursor == index {
			d.addHandle(task)
			return
		}

		// 任务尚未超时，放入wheel中等待触发
		bucket := wheel.buckets[index]
		bucket.Push(task)
		atomic.CompareAndSwapInt64(&bucket.expire, 0, expire)
		break
	}
}

func (d *Timer) run() {
	d.timestamp = time.Now().UnixMilli()
	tt := time.NewTicker(time.Duration(d.wheels[0].tick))
	for {
		select {
		case <-tt.C:
			for _, wheel := range d.wheels {
				// 判断是否切换了
				if !wheel.Next(time.Now().UnixMilli() - d.timestamp) {
					continue
				}

				// 执行定时任务
				d.addHandle(wheel.GetTask())
			}
		}
	}
}

func (d *Timer) addHandle(task *Task) {
	d.handles.Push(task)

	// 通知处理
	select {
	case d.handleCh <- struct{}{}:
	default:
	}
}

// 立即执行定时任务
func (d *Timer) handle() {
	for {
		select {
		case <-d.handleCh:
			for item := d.handles.Pop(); item != nil; item = d.handles.Pop() {
				vv := item.(*Task)
				for ; vv != nil; vv = vv.next {
					vv.handle()
					if !vv.isOnce {
						d.RegisterTimer(vv.handle, vv.ttl, vv.isOnce)
					}
				}
			}
		}
	}
}
