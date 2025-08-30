package recycle

import (
	"sync"
	"universal/library/queue"
	"universal/library/safe"
)

var (
	gc *Recycle
)

type IDestroy interface {
	Close()
}

type DestroyFunc func()

func (f DestroyFunc) Close() {
	f()
}

type Recycle struct {
	sync.WaitGroup
	tasks  *queue.Queue[IDestroy]
	notify chan struct{} // 消耗通知
	exit   chan struct{} // 退出
}

func Init() {
	gc = &Recycle{
		tasks:  queue.New[IDestroy](),
		notify: make(chan struct{}, 1),
		exit:   make(chan struct{}, 1),
	}
	gc.Add(1)
	safe.Go(run)
}

func Close() {
	close(gc.exit)
	gc.Wait()
}

func WrapIDestroy(f func()) IDestroy {
	return DestroyFunc(f)
}

func Destroy(fs ...IDestroy) {
	for _, f := range fs {
		gc.tasks.Push(f)
	}
	select {
	case gc.notify <- struct{}{}:
	default:
	}
}

func run() {
	defer func() {
		for f := gc.tasks.Pop(); f != nil; f = gc.tasks.Pop() {
			safe.Recover(f.Close)
		}
		gc.Done()
	}()
	for {
		select {
		case <-gc.notify:
			for f := gc.tasks.Pop(); f != nil; f = gc.tasks.Pop() {
				safe.Recover(f.Close)
			}
		case <-gc.exit:
			return
		}
	}
}
