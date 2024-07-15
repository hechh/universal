package base

import (
	"sync/atomic"
)

var (
	uuidGenerator uint64 // 全局唯一分配器
)

func AssignUUID() uint64 {
	return atomic.AddUint64(&uuidGenerator, 1)
}
