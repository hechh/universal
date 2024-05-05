package fbasic

import (
	"bytes"
	"sync"
)

var (
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(nil)
		},
	}
	bytePool = &sync.Pool{
		New: func() interface{} {
			return make([]byte, 1024*1024)
		},
	}
)

func GetBuffer() *bytes.Buffer {
	if obj, ok := bufferPool.Get().(*bytes.Buffer); ok && obj != nil {
		return obj
	}
	return bytes.NewBuffer(nil)
}

func PutBuffer(obj *bytes.Buffer) {
	if obj != nil {
		obj.Reset()
		bufferPool.Put(obj)
	}
}

func GetBytes() []byte {
	obj, _ := bytePool.Get().([]byte)
	return obj
}

func PutBytes(bs []byte) {
	bytePool.Put(bs)
}
