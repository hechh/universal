package router

import (
	"sync"
	"time"
	"universal/framework/define"
	"universal/library/baselib/safe"
	"universal/library/mlog"
)

type Table struct {
	newFunc define.NewRouterFunc      // 创建路由表
	mutex   *sync.RWMutex             // 互斥锁
	routers map[uint64]define.IRouter // 路由表
}

func NewTable(f define.NewRouterFunc) *Table {
	return &Table{
		newFunc: f,
		mutex:   new(sync.RWMutex),
		routers: make(map[uint64]define.IRouter),
	}
}

func (r *Table) get(id uint64) define.IRouter {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if val, ok := r.routers[id]; ok {
		return val
	}
	return nil
}

func (r *Table) Get(id uint64) define.IRouter {
	if val := r.get(id); val != nil {
		return val
	}
	item := r.newFunc(id)
	r.mutex.Lock()
	r.routers[id] = item
	r.mutex.Unlock()
	return item
}

func (r *Table) getExpires(ttl int64) (rets []define.IRouter) {
	now := time.Now().Unix()
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for _, vv := range r.routers {
		if vv.IsExpire(now, ttl) {
			rets = append(rets, vv)
		}
	}
	return
}

func (r *Table) Expire(ttl int64) {
	tt := time.NewTicker(time.Duration(ttl/2) * time.Second)
	defer tt.Stop()
	// 定时清理
	safe.SafeGo(mlog.Fatal, func() {
		for {
			<-tt.C
			items := r.getExpires(ttl)
			// 删除节点
			r.mutex.Lock()
			for _, val := range items {
				delete(r.routers, val.GetId())
			}
			r.mutex.Unlock()
		}
	})
}
