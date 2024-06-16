package timer

import (
	"time"
	"universal/framework/common/fbasic"
	"universal/framework/common/queue"
)

// 1秒=1000毫秒, 1毫秒=1000微妙, 1微妙=1000纳秒
const (
	INTERVAL_TIME = 100 * time.Millisecond // 10毫秒定时器
	TIME_NEAR     = 60
)

type Timer struct {
	time    uint64               // 触发次数
	short   *queue.Queue         // 二级任务队列
	long    *queue.Queue         // 二级任务队列
	near    [TIME_NEAR]*TaskList // 临近时间节点的定时器
	list    *NodeList            // 超长时长定时起
	shortCh chan struct{}        // 更新
	longCh  chan struct{}        // 更新
}

func (d *Timer) RegisterTimer(expire int64, f func()) {
	task := NewTask(f, expire)
	if expire > TIME_NEAR {
		d.short.Push(task)
	} else {
		d.long.Push(task)
	}
}

func (d *Timer) refreshShort() {
	now := fbasic.GetNow()
	// 移动任务
	for i := 1; i < TIME_NEAR; i++ {
		tasks := d.near[i].Remove()
		index := tasks.GetExpire(now)
		d.near[index].AddTask(tasks)
	}
	// 将任务队列中的数据更新到临近出发器
	for task := d.short.Pop(); task != nil; task = d.short.Pop() {
		tt := task.(*Task)
		d.near[tt.GetExpire(now)].AddTask(tt)
	}
}

func (d *Timer) executeShort() {
	for i := 0; i < TIME_NEAR; i++ {
		if d.near[i].GetExpire(fbasic.GetNow()) <= 0 {
			task := d.near[i].Remove()
			task.Handle()
		}
	}
}

func (d *Timer) run() {
	tt := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-tt.C:
		case <-d.shortCh:
			d.refreshShort()
		}
	}
}
