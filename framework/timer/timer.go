package timer

import (
	"sort"
	"sync/atomic"
	"time"
	"universal/library/async"
	"universal/library/mlog"
	"universal/library/uerror"
)

type Task struct {
	taskId *uint64
	task   func()
	ttl    int64
	expire int64
	times  int32
	next   *Task
}

type Wheel struct {
	mask    int64
	shift   int64
	smask   int64
	cursor  int64
	buckets []*Task
}

type Timer struct {
	tick      int64               // 最小时间间隔
	now       int64               // 当前时间
	startTime int64               // 启动时间
	size      int                 // 时间轮数量
	wheels    []*Wheel            // 时间轮
	tasks     *async.Queue[*Task] // 待插入定时任务队列
	notify    chan struct{}       // 插入通知
	exit      chan struct{}       // 定时器退出通知
}

func NewTimer(count, shift, tick int64) *Timer {
	now := time.Now().UnixMilli()
	ret := &Timer{
		tick:      1 << tick,
		now:       now,
		startTime: now,
		size:      int(count),
		wheels:    newWheels(count, shift, tick, now),
		tasks:     async.NewQueue[*Task](),
		notify:    make(chan struct{}, 1),
		exit:      make(chan struct{}),
	}
	async.SafeGo(mlog.Fatalf, ret.run)
	return ret
}

func newWheels(count, shift, tick, now int64) (rets []*Wheel) {
	size := int64(1) << shift
	for i := count - 1; i >= 0; i-- {
		rets = append(rets, &Wheel{
			mask:    size - 1,
			shift:   shift*i + tick,
			smask:   1<<(shift*i+tick) - 1,
			cursor:  now,
			buckets: make([]*Task, size),
		})
	}
	return
}

// 任务是否过期
func (d *Wheel) IsExpire(tt *Task) bool {
	return tt.expire <= d.cursor || tt.expire|d.smask <= d.cursor|d.smask
}

// 添加任务
func (d *Wheel) Add(task *Task) bool {
	if d.IsExpire(task) {
		return false
	}

	// 插入
	pos := (task.expire >> d.shift) & d.mask
	task.next = d.buckets[pos]
	d.buckets[pos] = task
	return true
}

// 获取过期任务
func (d *Wheel) Get(now int64) *Task {
	if now <= d.cursor || now|d.smask <= d.cursor|d.smask {
		return nil
	}

	// 加载超时任务
	pos := (now >> d.shift) & d.mask
	tt := d.buckets[pos]
	d.buckets[pos] = nil
	d.cursor = now
	return tt
}

// 注册定时器
func (d *Timer) Register(taskId *uint64, f func(), ttl time.Duration, times int32) error {
	tt := int64(ttl / time.Millisecond)
	if tt < d.tick {
		return uerror.New(1, -1, "定时器最小时间间隔:%dms", d.tick)
	}

	max := d.wheels[0]
	if (tt >> max.shift) > max.mask {
		return uerror.New(1, -1, "定时器超出最大时间范围:%dms", max.shift)
	}

	return d.add(&Task{
		taskId: taskId,
		task:   f,
		ttl:    tt,
		expire: atomic.LoadInt64(&d.now) + tt,
		times:  times,
	})
}

func (d *Timer) add(task *Task) error {
	d.tasks.Push(task)
	select {
	case d.notify <- struct{}{}:
	default:
	}
	return nil
}

func (d *Timer) run() {
	news := []*Task{}
	tt := time.NewTicker(time.Duration(d.tick) * time.Millisecond)
	defer func() {
		tt.Stop()
		news = news[:0]
		d.refresh(news)
		d.sync(news)
	}()

	for {
		select {
		case <-tt.C:
			atomic.AddInt64(&d.now, d.tick)
			news = news[:0]
			news = d.refresh(news)
			news = d.sync(news)
			d.insert(news)
		case <-d.notify:
			news = news[:0]
			news = d.sync(news)
			d.insert(news)
		case <-d.exit:
			return
		}
	}
}

func (d *Timer) refresh(news []*Task) []*Task {
	min := d.wheels[d.size-1]
	for i, wheel := range d.wheels {
		for item := wheel.Get(atomic.LoadInt64(&d.now)); item != nil; {
			tt := item
			item = item.next
			tt.next = nil
			if min.IsExpire(tt) {
				if newTask := d.handle(tt); newTask != nil {
					news = append(news, newTask)
				}
			} else {
				d.dispatcher(i+1, tt)
			}
		}
	}
	return news
}

func (d *Timer) sync(news []*Task) []*Task {
	min := d.wheels[d.size-1]
	for item := d.tasks.Pop(); item != nil; item = d.tasks.Pop() {
		if min.IsExpire(item) {
			if newTask := d.handle(item); newTask != nil {
				news = append(news, newTask)
			}
		} else {
			news = append(news, item)
		}
	}
	return news
}

func (d *Timer) handle(tt *Task) *Task {
	if *tt.taskId <= 0 || tt.times == 0 {
		return nil
	}
	async.SafeRecover(mlog.Fatalf, tt.task)
	if tt.times > 0 {
		tt.times--
	}
	if tt.times == 0 {
		return nil
	}
	tt.expire = atomic.LoadInt64(&d.now) + tt.ttl
	return tt
}

func (d *Timer) insert(tts []*Task) {
	if len(tts) <= 0 {
		return
	}

	now := atomic.LoadInt64(&d.now)
	sort.Slice(tts, func(i, j int) bool {
		return (tts[i].expire - now) > (tts[j].expire - now)
	})

	pos := 0
	for _, tt := range tts {
		for ; pos < d.size; pos++ {
			if d.wheels[pos].Add(tt) {
				break
			}
		}
	}
}

func (d *Timer) dispatcher(pos int, tt *Task) {
	for i := pos; i < d.size; i++ {
		if d.wheels[i].Add(tt) {
			return
		}
	}
}
