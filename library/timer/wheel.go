package timer

import "universal/library/util"

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
	cursor  int64
	buckets []*Task
}

func NewWheelList(tick int64, now int64) []*Wheel {
	return []*Wheel{
		{mask: (1 << 12) - 1, shift: tick, cursor: now, buckets: make([]*Task, 1<<12)},
		{mask: (1 << 5) - 1, shift: tick + 12, cursor: now, buckets: make([]*Task, 1<<5)},
		{mask: (1 << 5) - 1, shift: tick + 17, cursor: now, buckets: make([]*Task, 1<<5)},
		{mask: (1 << 5) - 1, shift: tick + 22, cursor: now, buckets: make([]*Task, 1<<5)},
		{mask: (1 << 5) - 1, shift: tick + 27, cursor: now, buckets: make([]*Task, 1<<5)},
	}
}

// 是否进位
func (w *Wheel) IsCarry() bool {
	return (w.cursor>>w.shift)&w.mask <= 0
}

// 是否过期
func (w *Wheel) IsExpire(tt *Task) bool {
	return tt.expire <= w.cursor || (tt.expire>>w.shift) <= (w.cursor>>w.shift)
}

// 是否匹配
func (w *Wheel) IsMatch(tt *Task) bool {
	return (tt.expire>>w.shift)-(w.cursor>>w.shift) <= w.mask
}

// 插入数据
func (w *Wheel) Insert(tt *Task) {
	pos := (tt.expire >> w.shift) & w.mask
	tt.next = w.buckets[pos]
	w.buckets[pos] = tt
}

// 获取过期定时任务
func (w *Wheel) Get(now int64) *Task {
	pos := (now >> w.shift) & w.mask
	ret := w.buckets[pos]
	w.buckets[pos] = nil
	w.cursor = now
	return ret
}

func (tt *Task) handle(now int64, fatal func(string, ...interface{})) *Task {
	if tt.taskId == nil || *tt.taskId <= 0 || tt.times == 0 {
		return nil
	}
	util.SafeRecover(fatal, tt.task)
	if tt.times > 0 {
		tt.times--
	}
	if tt.times != 0 {
		tt.expire = now + tt.ttl
		return tt
	}
	return nil
}
