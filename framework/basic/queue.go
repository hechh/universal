package basic

import (
	"sync/atomic"
	"unsafe"
)

type node struct {
	next  *node
	value interface{}
}

type Queue struct {
	head  *node
	tail  *node
	count int64
}

func NewQueue() *Queue {
	node := new(node)
	return &Queue{head: node, tail: node}
}

// 多协程安全
func (d *Queue) Push(val interface{}) {
	addNode := new(node)
	addNode.value = val

	// 将新增节点插入链表
	prevNode := (*node)(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&d.tail)), unsafe.Pointer(addNode)))
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&prevNode.next)), unsafe.Pointer(addNode))
	atomic.AddInt64(&d.count, 1)
}

func (d *Queue) Pop() interface{} {
	// 读取一个节点
	node := (*node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&d.head.next))))
	if node == nil {
		return nil
	}

	atomic.AddInt64(&d.count, -1)
	d.head.next = nil
	d.head = node
	return node.value
}
