package router

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"universal/common/pb"
	"universal/framework/domain"
)

type Table struct {
	ttl   int64
	exit  chan struct{}
	mutex sync.RWMutex
	infos map[string]domain.IRouter
}

func New() *Table {
	ret := &Table{
		ttl:   256 * 24 * 60 * 60,
		exit:  make(chan struct{}),
		infos: make(map[string]domain.IRouter),
	}
	return ret
}

func (d *Table) get(key string) domain.IRouter {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	if router, ok := d.infos[key]; ok {
		return router
	}
	return nil
}

func (d *Table) Get(nodeType pb.NodeType, actorName string, actorId uint64) domain.IRouter {
	// 读取路由信息
	key := fmt.Sprintf("%d:%s:%d", nodeType, actorName, actorId)
	if info := d.get(key); info != nil {
		return info
	}

	// 创建路由信息
	val := &Router{updateTime: time.Now().Unix(), Router: &pb.Router{}}
	d.mutex.Lock()
	d.infos[key] = val
	d.mutex.Unlock()
	return val
}

func (d *Table) SetExpire(ttl int64) {
	atomic.StoreInt64(&d.ttl, ttl)
}

func (r *Table) Close() {
	r.exit <- struct{}{}
}

func (r *Table) run() {
	tt := time.NewTicker(5 * 60 * time.Second)
	defer tt.Stop()

	for {
		select {
		case <-tt.C:
			now := time.Now().Unix()
			if rets := r.getExpires(now); len(rets) > 0 {
				r.mutex.Lock()
				for _, val := range rets {
					delete(r.infos, val)
				}
				r.mutex.Unlock()
			}
		case <-r.exit:
			return
		}
	}
}

func (r *Table) getExpires(now int64) (rets []string) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for routeId, val := range r.infos {
		if val.IsExpire(now, r.ttl) {
			rets = append(rets, routeId)
		}
	}
	return
}
