package timer

import (
	"sync/atomic"
	"time"
	"universal/framework/basic/async"
	"universal/framework/basic/util"
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

func INDEX(expire, n uint64) uint64 {
	return (expire >> (TVR_BITS + n*TVN_BITS)) & TVN_MASK
}

type Task struct {
	id     *uint64 // 定时器唯一ID
	f      func()  // 定时任务
	ttl    uint64  // 超时时间
	expire uint64  // 超时时间
	times  int64   // 重复执行多少次。-1表示执行无数次
}

type Timer struct {
	cursor uint64
	list   *async.Queue      // 添加任务队列
	tick   []*async.Queue    // 最小刻度转盘
	wheels [4][]*async.Queue // 时间轮
	timer  time.Ticker       // 定时器
	exit   chan struct{}     // 退出
}

func NewTask(id *uint64, expire uint64, f func(), ttl uint64, times int64) *Task {
	return &Task{
		id:     id,
		f:      f,
		ttl:    ttl,
		expire: expire,
		times:  times,
	}
}
func NewTimer() *Timer {
	return &Timer{
		list: async.NewQueue(),
		tick: async.NewQueuePool(TVR_SIZE),
		wheels: [4][]*async.Queue{
			async.NewQueuePool(TVN_SIZE),
			async.NewQueuePool(TVN_SIZE),
			async.NewQueuePool(TVN_SIZE),
			async.NewQueuePool(TVN_SIZE),
		},
		timer: *time.NewTicker(INTERVAL),
		exit:  make(chan struct{}, 1),
	}
}

func (d *Task) Handle(cur uint64) bool {
	if *d.id <= 0 || d.times == 0 {
		return false
	}
	// 执行定时任务
	d.f()
	if d.times > 0 {
		d.times--
	}
	d.expire = cur + d.ttl
	return d.times != 0
}

func (d *Timer) GetCursor() uint64 {
	return atomic.LoadUint64(&d.cursor)
}

func (d *Timer) AddCursor() uint64 {
	return atomic.AddUint64(&d.cursor, 1)
}

func (d *Timer) Insert(id *uint64, delay time.Duration, f func(), ttl time.Duration, times int64) {
	d.list.Push(&Task{
		id:     id,
		f:      f,
		ttl:    uint64(ttl / INTERVAL),
		expire: d.GetCursor() + uint64((ttl+delay)/INTERVAL),
		times:  times,
	})
}

func (d *Timer) insert(diff uint64, tt *Task) {
	if diff < TVR_SIZE {
		d.tick[tt.expire&TVR_MASK].Push(tt)
	} else if (diff >> TVR_BITS) < TVN_SIZE {
		d.wheels[0][INDEX(tt.expire, 0)].Push(tt)
	} else if (diff >> (TVR_BITS + TVN_BITS)) < TVN_SIZE {
		d.wheels[1][INDEX(tt.expire, 1)].Push(tt)
	} else if (diff >> (TVR_BITS + 2*TVN_BITS)) < TVN_SIZE {
		d.wheels[2][INDEX(tt.expire, 2)].Push(tt)
	} else if (diff >> (TVR_BITS + 3*TVN_BITS)) < TVN_SIZE {
		d.wheels[3][INDEX(tt.expire, 3)].Push(tt)
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
			d.insert(vv.expire-cur, vv)
			continue
		}
		// 处理过期定时任务
		if vv.Handle(cur) {
			d.insert(vv.expire-cur, vv)
		}
	}
}

// 处理过期定时任务
func (d *Timer) expire(cur uint64) {
	q := d.tick[cur&TVR_MASK]
	// 处理超时任务
	for tt := q.Pop(); tt != nil; tt = q.Pop() {
		vv := tt.(*Task)
		if vv.Handle(cur) {
			d.insert(vv.expire-cur, vv)
		}
	}
	// 迁移定时器任务
	for i, wl := range d.wheels {
		index := INDEX(cur, uint64(i))
		for tt := wl[index].Pop(); tt != nil; tt = wl[index].Pop() {
			vv := tt.(*Task)
			// 判断定时器是否有效
			if *vv.id <= 0 || vv.times == 0 {
				continue
			}
			// 迁移任务
			d.insert(vv.expire-cur, vv)
		}
	}
}

func (d *Timer) Update() {
	cur := d.AddCursor()
	d.expire(cur)
	d.move(cur)
}

func (d *Timer) Start() {
	util.SafeGo(nil, func() {
		for {
			select {
			case <-d.timer.C:
				d.Update()
			case <-d.exit:
				return
			}
		}
	})
}
