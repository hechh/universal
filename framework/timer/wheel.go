package timer

import "universal/library/mlog"

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

func NewWheels(count, shift, tick, now int64) (rets []*Wheel) {
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

// 获取过期任务
func (d *Wheel) Get(now int64) (rets []*Task) {
	if now <= d.cursor {
		return nil
	}

	// 获取所有过期的任务
	index := now >> d.shift
	cursor := d.cursor >> d.shift
	for i := cursor + 1; i < index; i++ {
		if task := d.buckets[i&d.mask]; task != nil {
			rets = append(rets, task)
			d.buckets[i&d.mask] = nil
		}
	}

	// 加载超时任务
	d.cursor = now
	return
}

// 添加任务
func (d *Wheel) Add(task *Task) bool {
	if task.expire <= d.cursor {
		return false
	}

	index := task.expire >> d.shift
	cursor := d.cursor >> d.shift
	if index <= cursor {
		return false
	}

	if index-cursor > d.mask+1 {
		mlog.Errorf("定时任务超出范围, taskId: %d, task: %v", *task.taskId, task)
		return true
	}

	// 插入
	task.next = d.buckets[index&d.mask]
	d.buckets[index&d.mask] = task
	return true
}
