package basic

import (
	"sync/atomic"
)

// actor唯一id分配器
var _AsyncIDSeed uint64

const (
	AsyncStatusRun  = 1
	AsyncStatusStop = 2
)

type Async struct {
	id     uint64        // 唯一id
	status int32         // actor运行状态
	tasks  *Queue        // 任务队列
	pushCh chan struct{} // 消耗通知
	exitCh chan struct{} // 退出
}

func NewAsync() *Async {
	return &Async{
		id:     atomic.AddUint64(&_AsyncIDSeed, 1),
		tasks:  NewQueue(),
		pushCh: make(chan struct{}, 1),
		exitCh: make(chan struct{}, 0),
	}
}

// 添加任务
func (d *Async) Push(task func()) {
	if atomic.CompareAndSwapInt32(&d.status, AsyncStatusRun, AsyncStatusRun) {
		d.tasks.Push(task)
		// 避免阻塞
		select {
		case d.pushCh <- struct{}{}:
		default:
		}
	}
}

// 开始actor任务协程
func (d *Async) Start() {
	if atomic.CompareAndSwapInt32(&d.status, AsyncStatusRun, AsyncStatusRun) {
		return
	}
	atomic.StoreInt32(&d.status, AsyncStatusRun)
	go d.run()
}

// 停止actor任务协程
func (d *Async) Stop() {
	if atomic.CompareAndSwapInt32(&d.status, AsyncStatusStop, AsyncStatusStop) {
		return
	}
	atomic.StoreInt32(&d.status, AsyncStatusStop)
	// 等待停止
	d.exitCh <- struct{}{}
}

func (d *Async) run() {
	consume := func() {
		for data := d.tasks.Pop(); data != nil; data = d.tasks.Pop() {
			(data.(func()))()
		}
	}
	for {
		select {
		case <-d.pushCh:
			consume()
		case <-d.exitCh:
			// 处理未完成的任务
			consume()
			return
		}
	}
}
