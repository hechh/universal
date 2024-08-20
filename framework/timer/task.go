package timer

import (
	"sync/atomic"
	"time"
	"universal/framework/basic/util"
	"unsafe"
)

type Task struct {
	once   bool   // 是否一次性
	ttl    int64  // 定时时长
	expire int64  // 过期时长
	handle func() // 任务
	next   *Task
}

type TaskQueue struct {
	head  *Task
	tail  *Task
	count int64
}

func NewTask(f func(), ttl int64, once bool) *Task {
	return &Task{
		handle: f,
		once:   once,
		ttl:    ttl,
		expire: util.GetNowTime().Add(time.Duration(ttl)).UnixMilli(),
	}
}

func NewTaskQueue() *TaskQueue {
	node := new(Task)
	return &TaskQueue{head: node, tail: node}
}

func NewTaskPool(size int) (rets []atomic.Value) {
	rets = make([]atomic.Value, size)
	for i := 0; i < size; i++ {
		rets[i].Store(NewTaskQueue())
	}
	return
}

func (d *Task) Handle() (rets []*Task) {
	for d != nil {
		item := d
		d = d.next
		item.next = nil

		item.handle()
		if !item.once {
			rets = append(rets, NewTask(item.handle, item.ttl, item.once))
		}
	}
	return
}

func (d *TaskQueue) Insert(task *Task) {
	prevNode := (*Task)(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&d.tail)), unsafe.Pointer(task)))
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&prevNode.next)), unsafe.Pointer(task))
	atomic.AddInt64(&d.count, 1)
}

func (d *TaskQueue) GetCount() int64 {
	return atomic.LoadInt64(&d.count)
}

func (d *TaskQueue) Remove() (tt *Task) {
	tt = d.head.next
	d.head.next = nil
	return
}
