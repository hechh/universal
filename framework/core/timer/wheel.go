package timer

import (
	"sort"
	"sync/atomic"
	"time"
	"unsafe"
)

// 定时任务
type Task struct {
	next   *Task  // 下一个
	expire int64  // 过期时间
	handle func() // 定时任务
}

// 定时任务列表
type TaskList struct {
	head   *Task
	tail   *Task
	expire int64
}

func NewTaskList() *TaskList {
	node := new(Task)
	return &TaskList{head: node, tail: node}
}

func (d *TaskList) SetExpire(val int64) {
	atomic.StoreInt64(&d.expire, val)
}

func (d *TaskList) GetExpire() int64 {
	return atomic.LoadInt64(&d.expire)
}

// 并发安全
func (d *TaskList) Push(item *Task) {
	newNode := unsafe.Pointer(item)
	// 将tail指针指向新添加大的item元素
	prevNode := (*Task)(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&d.tail)), newNode))
	// 将item接入链表中
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&prevNode.next)), newNode)
}

// 单协程安全
func (d *TaskList) Pop() *Task {
	item := (*Task)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&d.head.next))))
	if item == nil {
		return nil
	}
	d.head.next = nil
	d.head = item
	return item
}

// 延时任务队列
type DelayTaskList struct {
	list      []*TaskList
	handleCh  chan *TaskList
	notifyCh  chan *TaskList
	refreshCh chan struct{}
}

func (d *DelayTaskList) Push(tasks *TaskList) {
	d.notifyCh <- tasks
}

func (d *DelayTaskList) Pop() (val *TaskList) {
	if ll := len(d.list); ll > 0 {
		// 排序
		sort.Slice(d.list, func(i, j int) bool {
			return d.list[i].GetExpire() > d.list[j].GetExpire()
		})
		val = d.list[ll-1]
		d.list = d.list[:ll-1]
	}
	return
}

func (d *DelayTaskList) run() {
	for {
		select {
		case item := <-d.handleCh:
			// 等待到期
			<-time.After(time.Duration(item.GetExpire()))

			//执行任务
			for task := item.Pop(); task != nil; task = item.Pop() {
				task.handle()
			}

			// 更新下一个最短时间
			select {
			case d.refreshCh <- struct{}{}:
			default:
			}
		case item := <-d.notifyCh:
			// 添加任务
			d.list = append(d.list, item)

			select {
			case d.refreshCh <- struct{}{}:
			default:
			}
		case <-d.refreshCh:
			d.handleCh <- d.Pop()
		}
	}
}

// 任务轮(用户对定时任务分类、并插入相应任务队列)
type TaskWheel struct {
	tick      int64       // 时间刻度
	wheelSize int64       // 时间轮盘大小
	interval  int64       // 总时长
	cursor    int64       // 任务轮游标
	buckets   []*TaskList // 任务集合
}

func NewTaskWheel(tick time.Duration, wheelSize int64) *TaskWheel {
	return &TaskWheel{
		tick:      int64(tick),
		wheelSize: wheelSize,
		interval:  int64(tick) * wheelSize,
	}
}

func (d *TaskWheel) run() {
	tt := time.NewTicker(time.Duration(d.tick))
	for {
		select {
		case <-tt.C:
			// 更新游标
			index := atomic.AddInt64(&d.cursor, 1)
			// 执行定时任务
			taskList := d.buckets[index%d.wheelSize]
			for item := taskList.Pop(); item != nil; item = taskList.Pop() {
				item.handle()
			}
		}
	}
}

func (d *TaskWheel) Push(ttl int64, f func()) {
	// 取整
	index := ttl/d.tick + atomic.LoadInt64(&d.cursor)
	// 放置定时器
	d.buckets[index].Push(&Task{handle: f, expire: ttl})
}
