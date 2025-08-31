package mlog

import (
	"bytes"
	"sync"
	"sync/atomic"
	"time"
)

var (
	pools = sync.Pool{
		New: func() interface{} {
			return &meta{buf: bytes.NewBuffer(nil)}
		},
	}
)

func get() *meta {
	mm := pools.Get().(*meta)
	mm.buf.Reset()
	mm.tt = time.Now()
	atomic.StoreInt32(&mm.reference, 0)
	return mm
}

func put(m *meta) {
	if m.Done() == 0 {
		pools.Put(m)
	}
}

type meta struct {
	reference int32
	tt        time.Time
	buf       *bytes.Buffer
}

func (m *meta) Add() {
	atomic.AddInt32(&m.reference, 1)
}

func (m *meta) Done() int32 {
	return atomic.AddInt32(&m.reference, -1)
}
