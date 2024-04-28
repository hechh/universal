package fbasic

import (
	"bytes"
	"sync"
)

var (
	objectPool = &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(nil)
		},
	}
)

func GetBuffer() *bytes.Buffer {
	if obj, ok := objectPool.Get().(*bytes.Buffer); ok && obj != nil {
		return obj
	}
	return bytes.NewBuffer(nil)
}

func PutBuffer(obj *bytes.Buffer) {
	if obj != nil {
		obj.Reset()
		objectPool.Put(obj)
	}
}
