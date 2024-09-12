package timer

import (
	"sync/atomic"
	"time"
	"universal/framework/basic/async"
	"universal/framework/basic/util"
	"universal/framework/plog"
)

const (
	TVR_BITS = 8
	TVR_SIZE = (1 << TVR_BITS)
	TVR_MASK = (TVR_SIZE - 1)
	TVN_BITS = 6
	TVN_SIZE = (1 << TVN_BITS)
	TVN_MASK = (TVN_SIZE - 1)
	INTERVAL = 10 * time.Millisecond
)

type Task struct {
	id     *uint64 // 定时器唯一ID
	f      func()  // 定时任务
	ttl    uint64  // 超时时间
	expire uint64  // 超时时间
	times  int64   // 重复执行多少次。-1表示执行无数次
}

type Timer struct {
	handle *async.Async      // 定时任务处理程序
	cursor uint64            // 游标
	list   *async.Queue      // 添加任务队列
	tick   []*async.Queue    // 最小刻度转盘
	wheels [4][]*async.Queue // 时间轮
	timer  time.Ticker       // 定时器
	exit   chan struct{}     // 退出
}

func INDEX(expire, n uint64) uint64 {
	return (expire >> (TVR_BITS + n*TVN_BITS)) & TVN_MASK
}

func NewTimer(tick time.Duration) *Timer {
	return &Timer{
		handle: async.NewAsync(),
		list:   async.NewQueue(),
		tick:   async.NewQueuePool(TVR_SIZE),
		wheels: [4][]*async.Queue{async.NewQueuePool(TVN_SIZE), async.NewQueuePool(TVN_SIZE), async.NewQueuePool(TVN_SIZE), async.NewQueuePool(TVN_SIZE)},
		timer:  *time.NewTicker(tick),
		exit:   make(chan struct{}, 1),
	}
}

// 执行定时任务
func (d *Task) Update(cur uint64) bool {
	if d.times > 0 {
		d.times--
	}
	d.expire = cur + d.ttl
	return d.times != 0
}

// 注册定时器
func (d *Timer) RegisterTimer(id *uint64, f func(), ttl, delay time.Duration, times int64) {
	d.list.Push(&Task{
		id:     id,
		f:      f,
		ttl:    uint64(ttl / INTERVAL),
		expire: atomic.LoadUint64(&d.cursor) + uint64((ttl+delay)/INTERVAL),
		times:  times,
	})
}

// 插入定时任务
func (d *Timer) insert(cur uint64, tt *Task) {
	diff := tt.expire - cur
	if diff < TVR_SIZE {
		pos := tt.expire & TVR_MASK
		d.tick[pos].Push(tt)
	} else if (diff >> TVR_BITS) < TVN_SIZE {
		pos := INDEX(tt.expire, 0)
		d.wheels[0][pos].Push(tt)
	} else if (diff >> (TVR_BITS + TVN_BITS)) < TVN_SIZE {
		pos := INDEX(tt.expire, 1)
		d.wheels[1][pos].Push(tt)
	} else if (diff >> (TVR_BITS + 2*TVN_BITS)) < TVN_SIZE {
		pos := INDEX(tt.expire, 2)
		d.wheels[2][pos].Push(tt)
	} else if (diff >> (TVR_BITS + 3*TVN_BITS)) < TVN_SIZE {
		pos := INDEX(tt.expire, 3)
		d.wheels[3][pos].Push(tt)
	}
}

// 定时任务缓存队列插入
func (d *Timer) move(cur uint64) {
	for tt := d.list.Pop(); tt != nil; tt = d.list.Pop() {
		vv := tt.(*Task)
		// 判断定时器是否有效
		if *vv.id <= 0 || vv.times == 0 {
			continue
		}
		// 判断是否过期
		if cur < vv.expire {
			d.insert(cur, vv)
			continue
		}
		// 处理过期定时任务
		d.handle.Push(vv.f)
		// 更新定时任务
		if vv.Update(cur) {
			d.insert(cur, vv)
		}
	}
}

// 处理过期定时任务
func (d *Timer) expire(cur uint64) {
	// 处理超时任务
	pos := cur & TVR_MASK
	for tt := d.tick[pos].Pop(); tt != nil; tt = d.tick[pos].Pop() {
		vv := tt.(*Task)
		// 判断定时器是否有效
		if *vv.id <= 0 || vv.times == 0 {
			continue
		}
		// 处理过期定时任务
		d.handle.Push(vv.f)
		// 更新定时任务
		if vv.Update(cur) {
			d.insert(cur, vv)
		}
	}
}

// 迁移定时器任务
func (d *Timer) shift(cur uint64) {
	// 迁移定时器任务
	for i, wl := range d.wheels {
		pos := INDEX(cur, uint64(i))
		for tt := wl[pos].Pop(); tt != nil; tt = wl[pos].Pop() {
			vv := tt.(*Task)
			// 判断定时器是否有效
			if *vv.id <= 0 || vv.times == 0 {
				continue
			}
			// 迁移任务
			d.insert(cur, vv)
		}
	}
}

func (d *Timer) Stop() {
	d.exit <- struct{}{}
	d.handle.Stop()
}

func (d *Timer) Start() {
	d.handle.Start()
	util.SafeGo(plog.Catch, func() {
		for {
			select {
			case <-d.timer.C:
				cur := atomic.AddUint64(&d.cursor, 1)
				d.move(cur)
				d.shift(cur)
				d.expire(cur)
			case <-d.exit:
				return
			}
		}
	})
}

func (d *Timer) StartTest(times int) {
	defer d.Stop()
	d.handle.Start()
	for i := 0; i < times; i++ {
		<-d.timer.C
		cur := atomic.AddUint64(&d.cursor, 1)
		d.move(cur)
		d.shift(cur)
		d.expire(cur)
	}
}
