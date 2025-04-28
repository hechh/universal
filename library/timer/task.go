package timer

import "universal/library/baselib/safe"

type Task struct {
	taskId *uint64
	task   func()
	ttl    int64
	expire int64
	times  int32
	next   *Task
}

func (t *Task) Do(now int64, cb func(interface{})) *Task {
	if *t.taskId > 0 {
		// 执行定时任务
		safe.SafeGo(cb, t.task)
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
