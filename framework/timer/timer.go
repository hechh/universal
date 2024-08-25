package timer

import (
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
	"universal/framework/basic/async"
	"universal/framework/basic/util"
	"universal/framework/plog"
)

var (
	generator uint64
	object    *Timer = NewTimer()
)

type Task struct {
	id     uint64        // 唯一id
	handle func()        // 任务
	ttl    time.Duration // 定时时长
	once   bool          // 是否一次性
	expire int64         // 过期时长
}

type Wheel struct {
	cursor   int64          // 处理游标
	tickSize int64          // tick的bit数
	tick     int64          // 时间刻度
	interval int64          // 时间间隔
	size     int64          // 转盘大小
	buckets  []*async.Queue // 任务队列
}

type Timer struct {
	sync.WaitGroup
	wheels     [4]*Wheel     // 定时任务转盘
	over       *async.Queue  // 过期任务
	updateOver chan struct{} // 过期通知
	exitRun    chan struct{} // 退出通知
	exitHandle chan struct{} // 退出通知
}

func Insert(tt *Task) {
	object.Insert(tt)
}

func Stop() {
	object.Stop()
}

func NewTask(f func(), ttl time.Duration, once bool) *Task {
	return &Task{
		id:     atomic.AddUint64(&generator, 1),
		handle: f,
		ttl:    ttl,
		once:   once,
		expire: util.GetNowTime().Add(ttl).UnixMilli(),
	}
}

func NewDelayTask(f func(), ttl time.Duration, once bool, interval time.Duration) *Task {
	return &Task{
		id:     atomic.AddUint64(&generator, 1),
		handle: f,
		ttl:    ttl,
		once:   once,
		expire: util.GetNowTime().Add(interval).UnixMilli(),
	}
}

func NewWheel(tick, size int64) *Wheel {
	ll := (1 << size)
	rets := make([]*async.Queue, ll)
	for i := 0; i < ll; i++ {
		rets[i] = async.NewQueue()
	}
	return &Wheel{
		cursor:   util.GetNowUnixMilli(),
		tickSize: tick,
		tick:     1<<tick - 1,
		interval: 1<<(tick+size) - 1,
		size:     1<<size - 1,
		buckets:  rets,
	}
}

func NewTimer() *Timer {
	ret := &Timer{
		wheels:     [4]*Wheel{NewWheel(6, 7), NewWheel(13, 7), NewWheel(20, 7), NewWheel(27, 7)},
		over:       async.NewQueue(),
		updateOver: make(chan struct{}, 1),
		exitRun:    make(chan struct{}, 1),
		exitHandle: make(chan struct{}, 1),
	}
	// 定时触发
	util.SafeGo(func(err interface{}) {
		plog.Fatal("%v\n stack: %s", err, string(debug.Stack()))
	}, ret.run)
	// 执行定时任务
	ret.Add(1)
	util.SafeGo(func(err interface{}) {
		plog.Fatal("%v\n stack: %s", err, string(debug.Stack()))
	}, ret.handle)
	return ret
}

func (d *Task) Update(now int64) {
	ttl := int64(d.ttl / time.Millisecond)
	for d.expire <= now {
		d.expire += ttl
	}
}

func (d *Timer) Stop() {
	d.exitRun <- struct{}{}
	d.exitHandle <- struct{}{}
	d.Wait()
}

func (d *Timer) addOver(tt ...*Task) {
	d.over.Push(tt)
	select {
	case d.updateOver <- struct{}{}:
	default:
	}
}

func (d *Timer) Insert(tt *Task) {
	for _, wheel := range d.wheels {
		switch wheel.Insert(tt) {
		case -1:
			d.addOver(tt)
			return
		case 0:
			return
		}
	}
}

// 过滤过期任务执行
func (d *Timer) run() {
	tt := time.NewTicker(time.Duration(d.wheels[0].tick) * time.Millisecond)
	for {
		select {
		case <-d.exitRun:
			return
		case <-tt.C:
			overs := []*Task{}
			now := util.GetNowUnixMilli()
			for _, wl := range d.wheels {
				for _, item := range wl.Pop(now) {
					if !d.wheels[0].IsExpired(now, item.expire) {
						d.Insert(item)
					} else {
						if !item.once {
							item.Update(now)
							d.Insert(item)
						}
						overs = append(overs, item)
					}
				}
			}
			d.addOver(overs...)
		}
	}
}

// 过期任务执行
func (d *Timer) handle() {
	handleF := func() {
		for tt := d.over.Pop(); tt != nil; tt = d.over.Pop() {
			switch tmp := tt.(type) {
			case *Task:
				tmp.handle()
			case []*Task:
				for _, vv := range tmp {
					vv.handle()
				}
			}
		}
	}
	defer func() {
		handleF()
		d.Done()
	}()
	for {
		select {
		case <-d.updateOver:
			handleF()
		case <-d.exitHandle:
			return
		}
	}
}

func (d *Wheel) GetIndex(val int64) int64 {
	return (val >> d.tickSize) & d.size
}

// 判断是否过期
func (d *Wheel) IsExpired(now, expire int64) bool {
	return (now >> d.tickSize) > (expire >> d.tickSize)
}

// 插入任务
func (d *Wheel) Insert(tt *Task) int {
	cursor := util.GetNowUnixMilli()
	if d.IsExpired(cursor, tt.expire) {
		return -1
	}
	if tt.expire-cursor >= d.interval {
		return 1
	}
	// 插入成功
	//plog.Trace("%d: index: %d, Task--->id: %d, now: %d, expire: %d, ttl: %d, once: %v", d.index, d.GetIndex(tt.expire), tt.id, cursor, tt.expire, tt.ttl, tt.once)
	d.buckets[d.GetIndex(tt.expire)].Push(tt)
	return 0
}

// 获取过期或者需要迁移的任务
func (d *Wheel) Pop(now int64) (rets []*Task) {
	cursor := atomic.SwapInt64(&d.cursor, now)
	begin, end := d.GetIndex(cursor), d.GetIndex(now)
	if end < begin {
		end += d.size
	}
	// 读取过期
	//plog.Trace("%d: now: %d, begin: %d, end: %d, count: %d", d.index, now, begin, end, d.buckets[begin].GetCount())
	for i := begin; i < end; i++ {
		for tt := d.buckets[i&d.size].Pop(); tt != nil; tt = d.buckets[i&d.size].Pop() {
			rets = append(rets, tt.(*Task))
		}
	}
	//plog.Trace("%d: now: %d, rets: %d", d.index, now, len(rets))
	return
}
