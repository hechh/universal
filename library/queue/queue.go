package queue

import (
	"sync/atomic"
	"unsafe"
)

type Node[T any] struct {
	next  *Node[T]
	value T
}

type Queue[T any] struct {
	head  *Node[T]
	tail  *Node[T]
	count int32
}

func New[T any]() *Queue[T] {
	nn := new(Node[T])
	return &Queue[T]{head: nn, tail: nn}
}

func (d *Queue[T]) Size() int32 {
	return atomic.LoadInt32(&d.count)
}

func (d *Queue[T]) Push(val T) {
	addNode := new(Node[T])
	addNode.value = val
	prevNode := (*Node[T])(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&d.tail)), unsafe.Pointer(addNode)))
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&prevNode.next)), unsafe.Pointer(addNode))
	atomic.AddInt32(&d.count, 1)
}

func (d *Queue[T]) Pop() (ret T) {
	if node := (*Node[T])(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&d.head.next)))); node != nil {
		atomic.AddInt32(&d.count, -1)
		ret = node.value
		d.head.next = nil
		d.head = node
	}
	return
}
