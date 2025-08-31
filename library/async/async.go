package async

import (
	"sync"
	"sync/atomic"
	"universal/library/queue"
	"universal/library/safe"
)

type Async struct {
	sync.WaitGroup
	id     uint64
	status int32
	tasks  *queue.Queue[func()]
	notify chan struct{}
	exit   chan struct{}
}

func New() *Async {
	return &Async{
		tasks:  queue.New[func()](),
		notify: make(chan struct{}, 1),
		exit:   make(chan struct{}),
	}
}

func (d *Async) GetIdPointer() *uint64 {
	return &d.id
}

func (d *Async) GetId() uint64 {
	return atomic.LoadUint64(&d.id)
}

func (d *Async) SetId(id uint64) {
	atomic.StoreUint64(&d.id, id)
}

func (d *Async) Start() {
	if atomic.LoadInt32(&d.status) > 0 {
		return
	}
	atomic.AddInt32(&d.status, 1)
	d.Add(1)
	safe.Go(d.run)
}

func (d *Async) Stop() {
	if atomic.LoadInt32(&d.status) <= 0 {
		return
	}
	atomic.StoreInt32(&d.status, 0)
	atomic.StoreUint64(&d.id, 0)
	close(d.exit)
	d.Wait()
}

func (d *Async) Push(f func()) {
	d.tasks.Push(f)
	select {
	case d.notify <- struct{}{}:
	default:
	}
}

func (d *Async) run() {
	defer func() {
		for f := d.tasks.Pop(); f != nil; f = d.tasks.Pop() {
			safe.Recover(f)
		}
		d.Done()
	}()
	for {
		select {
		case <-d.notify:
			for f := d.tasks.Pop(); f != nil; f = d.tasks.Pop() {
				safe.Recover(f)
			}
		case <-d.exit:
			return
		}
	}
}
