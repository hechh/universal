package router

import (
	"fmt"
	"poker_server/common/pb"
	"poker_server/framework/domain"
	"poker_server/library/safe"
	"sync"
	"time"
)

type Table struct {
	mutex sync.RWMutex
	data  map[string]*Router
	exit  chan struct{}
	ttl   int64
}

func New(ttl int64) *Table {
	ret := &Table{
		exit: make(chan struct{}),
		data: make(map[string]*Router),
		ttl:  ttl,
	}
	safe.Go(ret.run)
	return ret
}

func getkey(routerType uint32, id uint64) string {
	return fmt.Sprintf("%d:%d", routerType, id)
}

func (d *Table) Get(routerType uint32, id uint64) domain.IRouter {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	if val, ok := d.data[getkey(routerType, id)]; ok {
		return val
	}
	return nil
}

func (d *Table) GetOrNew(routerType uint32, id uint64, nn *pb.Node) domain.IRouter {
	if rr := d.Get(routerType, id); rr != nil {
		return rr
	}

	// 创建路由信息
	val := &Router{updateTime: time.Now().Unix(), Router: &pb.Router{}}
	val.Set(nn.Type, nn.Id)

	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.data[getkey(routerType, id)] = val
	return val
}

func (r *Table) Close() {
	close(r.exit)
}

func (t *Table) Walk(f func(string, *Router) bool) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	for id, rr := range t.data {
		if !f(id, rr) {
			return
		}
	}
}

func (t *Table) Remove(ids ...string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	for _, id := range ids {
		delete(t.data, id)
	}
}

func (t *Table) run() {
	tt := time.NewTicker(time.Duration(t.ttl/2) * time.Second)
	defer tt.Stop()
	for {
		select {
		case <-tt.C:
			now := time.Now().Unix()
			dels := []string{}
			t.Walk(func(id string, rr *Router) bool {
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
