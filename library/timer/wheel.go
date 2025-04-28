package timer

import "universal/library/baselib/uerror"

type Wheel struct {
	mask    int64
	shift   int64   // 时间刻度偏移量
	cursor  int64   // 当前时间
	buckets []*Task // 任务队列
}

func newWheelList(count, shift, tick, now int64) (rets []*Wheel) {
	size := int64(1) << shift
	for i := count - 1; i >= 0; i-- {
		rets = append(rets, &Wheel{
			mask:    size - 1,
			shift:   shift*i + tick,
			cursor:  now,
			buckets: make([]*Task, size),
		})
	}
	return
}

// 定时器触发
func (d *Wheel) Get(now int64) *Task {
	// 判断是否可以获取到期任务
	if now>>d.shift <= d.cursor>>d.shift {
		return nil
	}

	// 加载超时任务
	d.cursor = now
	pos := (now >> d.shift) & d.mask
	task := d.buckets[pos]
	d.buckets[pos] = nil
	return task
}

// 添加定时任务
func (d *Wheel) Add(task *Task, cb func(interface{})) bool {
	// 已经过期
	expire := task.expire >> d.shift
	cursor := d.cursor >> d.shift
	if task.expire <= d.cursor || expire <= cursor {
		return false
	}
	// 超出了该定时器能够存储的范围
	if expire >= cursor+d.mask+1 {
		// todo
		err := uerror.New(1, -1, "不支持定时时长超出定时器最大范围: %d", task.ttl)
		cb(err)
		return true
	}
	// 插入
	task.next = d.buckets[expire&d.mask]
	d.buckets[expire&d.mask] = task
	return true
}
