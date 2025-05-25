package timer

type Task struct {
	taskId *uint64
	task   func()
	ttl    int64
	expire int64
	times  int32
	next   *Task
}

func (t *Task) Do(now int64) *Task {
	if *t.taskId > 0 {
		// 执行定时任务
		t.task()
		// 刷新超时时间
		t.expire = now + t.ttl
		t.next = nil
		//  减少执行次数
		if t.times > 0 {
			t.times--
		}
		if t.times != 0 {
			return t
		}
	}
	return nil
}

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
	index := now >> d.shift
	cursor := d.cursor >> d.shift
	// 判断是否可以获取到期任务
	if index <= cursor {
		return nil
	}

	// 加载超时任务
	d.cursor = now
	pos := index & d.mask
	task := d.buckets[pos]
	d.buckets[pos] = nil
	return task
}

// 添加定时任务
func (d *Wheel) Add(task *Task) bool {
	// 任务已经过期
	if task.expire <= d.cursor {
		return false
	}

	// 已经过期
	expire := task.expire >> d.shift
	cursor := d.cursor >> d.shift
	if expire <= cursor {
		return false
	}

	// 超出了该定时器能够存储的范围
	if expire >= cursor+d.mask+1 {
		return true
	}

	// 插入
	task.next = d.buckets[expire&d.mask]
	d.buckets[expire&d.mask] = task
	return true
}
