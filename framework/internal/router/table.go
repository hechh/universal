package router

import (
	"sync"
	"time"
	"universal/common/pb"
	"universal/framework/define"
	"universal/library/safe"
)

type Table struct {
	mutex sync.RWMutex
	data  map[uint64]*Router
	exit  chan struct{}
	ttl   int64
}

func NewTable(ttl int64) *Table {
	ret := &Table{
		data: make(map[uint64]*Router),
		exit: make(chan struct{}),
		ttl:  ttl,
	}
	safe.Go(ret.run)
	return ret
}

func (t *Table) Get(id uint64) define.IRouter {
	t.mutex.RLock()
	defer t.mutex.RLock()
	if val, ok := t.data[id]; ok {
		return val
	}
	return nil
}

func (t *Table) GetOrNew(id uint64, self *pb.Node) define.IRouter {
	if rr := t.Get(id); rr != nil {
		return rr
	}

	// 创建路由信息
	rr := &Router{Router: &pb.Router{}, updateTime: time.Now().Unix()}
	rr.Set(self.Type, self.Id)

	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.data[id] = rr
	return rr
}

func (t *Table) Walk(f func(uint64, *Router) bool) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	for id, rr := range t.data {
		if !f(id, rr) {
			return
		}
	}
}

func (t *Table) Remove(ids ...uint64) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	for _, id := range ids {
		delete(t.data, id)
	}
}

func (t *Table) Close() {
	close(t.exit)
}

func (t *Table) run() {
	tt := time.NewTicker(time.Duration(t.ttl/2) * time.Second)
	defer tt.Stop()

	for {
		select {
		case <-tt.C:
			now := time.Now().Unix()
			dels := []uint64{}
			t.Walk(func(id uint64, rr *Router) bool {
				if t.ttl >= now-rr.GetUpdateTime() {
					dels = append(dels, id)
				}
				return true
			})
			t.Remove(dels...)
		case <-t.exit:
			return
		}
	}
}
