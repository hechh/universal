package timer

import (
	"fmt"
	"universal/framework/common/fbasic"
)

type Task struct {
	timestamp int64  // 任务开启时间点(秒)
	expire    int64  // 定时时长(秒)
	task      func() // 任务
	isOnce    bool   // 是否一次性任务
	next      *Task
}

func NewTask(f func(), expire int64) *Task {
	return &Task{
		timestamp: fbasic.GetNow(),
		expire:    expire,
		task:      f,
	}
}

func (d *Task) GetExpire(val int64) int64 {
	if diff := val - d.timestamp; diff < d.expire {
		return d.expire - diff
	}
	return 0
}

func (d *Task) Handle() {

}

func (d *Task) ToString(now int64) string {
	return fmt.Sprintf("isOnce: %v, expire: %d", d.isOnce, d.GetExpire(now))
}

type TaskList struct {
	head Task
	tail *Task
	prev *TaskList
	next *TaskList
}

func NewTaskList(t *Task) *TaskList {
	ret := new(TaskList)
	ret.tail = &ret.head
	ret.AddTask(t)
	return ret
}

func (d *TaskList) AddTask(node *Task) {
	for ; node != nil; node = node.next {
		d.tail.next = node
		d.tail = node
	}
}

func (d *TaskList) GetExpire(val int64) int64 {
	if diff := val - d.tail.timestamp; diff < d.tail.expire {
		return d.tail.expire - diff
	}
	return 0
}

func (d *TaskList) Remove() (ret *Task) {
	ret = d.head.next
	d.head.next = nil
	d.tail = &d.head
	return
}

func (d *TaskList) Print() {
	now := fbasic.GetNow()
	for item := d.head.next; item != nil; item = item.next {
		fmt.Println(item.ToString(now))
	}
}

// ---------列表----------
type NodeList struct {
	head *TaskList
	tail *TaskList
}

func (d *NodeList) Print() {
	for item := d.head; item != nil; item = item.next {
		item.Print()
	}
}

func (d *NodeList) Insert(node *Task) {
	if d.head == nil {
		taskList := NewTaskList(node)
		d.head = taskList
		d.tail = taskList
		return
	}
	now := fbasic.GetNow()
	expire := node.GetExpire(now)
	// 搜索插入位置
	item := d.head
	for ; item != nil && item.GetExpire(now) < expire; item = item.next {
	}
	// 插入表尾
	if item == nil {
		taskList := NewTaskList(node)
		taskList.prev = d.tail
		d.tail.next = taskList
		d.tail = taskList
		return
	}
	// 相等
	if item.GetExpire(now) == expire {
		item.AddTask(node)
		return
	}
	// 插入表头
	if item == d.head {
		taskList := NewTaskList(node)
		taskList.next = d.head
		d.head.prev = taskList
		d.head = taskList
		return
	}
	// 中间插入
	taskList := NewTaskList(node)
	taskList.next = item
	taskList.prev = item.prev
	if item.prev.next != nil {
		item.prev.next = taskList
	}
	item.prev = taskList
	return
}

func (d *NodeList) Pop(now, val int64) (ret []*TaskList) {
	for item := d.head; item != nil && item.GetExpire(now) <= val; {
		if item.next == nil {
			d.head = nil
			d.tail = nil
			ret = append(ret, item)
			break
		}

		d.head = item.next
		d.head.prev = nil
		item.next = nil
		ret = append(ret, item)
		item = d.head
	}
	return
}
