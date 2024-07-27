package timer

import (
	"sync/atomic"
	"time"
	"universal/framework/util"
	"unsafe"
)

type Task struct {
	task   func() // 定时任务
	isOnce bool   // 是否一次性定时器
	expire int64  // 过期时间
	ttl    int64  // 定时时长
	next   *Task
}

func (d *Task) Handle() {
	d.task()
}

type TaskList struct {
	head  *Task
	tail  *Task
	count int64
}

func (d *TaskList) GetCount() int64 {
	return atomic.LoadInt64(&d.count)
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

func NewTask(f func(), ttl int64, isOnce bool) *Task {
	return &Task{
		task:   f,
		isOnce: isOnce,
		ttl:    ttl,
		expire: util.GetNowTime().Add(time.Duration(ttl)).UnixMilli(),
	}
}

func NewTaskList() *TaskList {
	node := new(Task)
	return &TaskList{head: node, tail: node}
}

func NewTaskBucket(size int64) []*TaskList {
	rets := make([]*TaskList, size)
	for i := int64(0); i < size; i++ {
		rets[i] = NewTaskList()
	}
	return rets
}
