package timer

import (
	"sync"
	"time"
	"universal/framework/util"
)

// 毫秒 --->>> tick：10ms, wheelSize: 100, interval: 1000ms=1s
// 秒 --->>> tick: 1s, wheelSize: 120, interval: 120s=2min
// 分 --->>> tick: 2m, wheelSize: 120, interval: 240m=4hour
// 时 --->>> tick: 4h, wheelSize: 120, interval: 480h=20day

type Timer struct {
	sync.WaitGroup
	wheels        [4]*Wheel     // 定时任务队列
	overTimeTasks *TaskList     // 超时任务队列
	overTimeCh    chan struct{} // 超时任务处理通知
	exitHandleCh  chan struct{} // 退出定时器
	exitRunCh     chan struct{} // 退出定时器
}

func NewTimer() *Timer {
	now := util.GetNowUnixMilli()
	ret := &Timer{
		wheels:        [4]*Wheel{NewWheel(now, 6, 7), NewWheel(now, 13, 7), NewWheel(now, 20, 7), NewWheel(now, 27, 7)},
		overTimeTasks: NewTaskList(),
		overTimeCh:    make(chan struct{}, 2),
		exitHandleCh:  make(chan struct{}, 1),
		exitRunCh:     make(chan struct{}, 1),
	}
	go ret.handle() // 处理超时任务
	go ret.run()    // 定时处理任务
	return ret
}

func (d *Timer) Insert(f func(), ttl time.Duration, isOnce bool) {
	d.add(&Task{
		task:   f,
		isOnce: isOnce,
		ttl:    int64(ttl),
		expire: util.GetNowTime().Add(ttl).UnixMilli(),
	})
}

func (d *Timer) Stop() {
	d.exitRunCh <- struct{}{}
	d.exitHandleCh <- struct{}{}
	d.Wait()
}

func (d *Timer) run() {
	d.Add(1)
	tt := time.NewTicker(time.Duration(d.wheels[0].tick) * time.Millisecond)
	defer func() {
		tt.Stop()
		d.Done()
	}()

	for {
		select {
		case <-tt.C:
			// 执行触发的
			for _, wheel := range d.wheels {
				if task := wheel.Pop(); task != nil {
					d.overTimeTasks.Insert(task)
					select {
					case d.overTimeCh <- struct{}{}:
					default:
					}
				}
			}
		case <-d.exitRunCh:
			return
		}
	}
}

func (d *Timer) add(task *Task) {
	for _, wheel := range d.wheels {
		switch wheel.Insert(task) {
		case -1:
			d.overTimeTasks.Insert(task)
			select {
			case d.overTimeCh <- struct{}{}:
			default:
			}
			return
		case 0:
			return
		}
	}
}

// 处理超时任务
func (d *Timer) handle() {
	wheel := d.wheels[0]
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
				for task != nil {
					item := task
					task = task.next
					item.next = nil
					now := util.GetNowUnixMilli()
					// 判断是否过期
					if item.expire-now <= wheel.tick {
						item.task()
						if !item.isOnce {
							d.Insert(item.task, time.Duration(item.ttl), item.isOnce)
						}
					} else {
						d.add(item)
					}
				}
			}
		case <-d.exitHandleCh:
			return
		}
	}
}
