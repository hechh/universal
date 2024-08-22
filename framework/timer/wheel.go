package timer

import (
	"sync/atomic"
	"time"
	"universal/framework/basic/async"
	"universal/framework/basic/util"
)

var (
	generator uint64
)

type Task struct {
	id     uint64        // 唯一id
	handle func()        // 任务
	ttl    time.Duration // 定时时长
	once   bool          // 是否一次性
	expire int64         // 过期时长
}

type Wheel struct {
	index    int            // 序号
	cursor   int64          // 处理游标
	tickSize int64          // tick的bit数
	tick     int64          // 时间刻度
	interval int64          // 时间间隔
	size     int64          // 转盘大小
	buckets  []*async.Queue // 任务队列
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

func (d *Task) refresh(now int64) {
	ttl := int64(d.ttl / time.Millisecond)
	for d.expire <= now {
		d.expire += ttl
	}
}

func NewWheel(index int, tick, size int64) *Wheel {
	ll := (1 << size)
	rets := make([]*async.Queue, ll)
	for i := 0; i < ll; i++ {
		rets[i] = async.NewQueue()
	}
	return &Wheel{
		index:    index,
		cursor:   util.GetNowUnixMilli(),
		tickSize: tick,
		tick:     1<<tick - 1,
		interval: 1<<(tick+size) - 1,
		size:     1<<size - 1,
		buckets:  rets,
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
func (d *Wheel) Pop() (rets []*Task) {
	now := util.GetNowUnixMilli()
	cursor := atomic.SwapInt64(&d.cursor, now)
	begin, end := d.GetIndex(cursor), d.GetIndex(now)
	if end < begin {
		end += d.size
	}
	//plog.Trace("%d: now: %d, begin: %d, end: %d, count: %d", d.index, now, begin, end, d.buckets[begin].GetCount())
	// 读取过期
	for i := begin; i < end; i++ {
		old := d.buckets[i&d.size]
		for tt := old.Pop(); tt != nil; tt = old.Pop() {
			rets = append(rets, tt.(*Task))
		}
	}
	//plog.Trace("%d: now: %d, rets: %d", d.index, now, len(rets))
	return
}
