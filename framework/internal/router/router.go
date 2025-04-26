package router

import (
	"sync"
	"time"
	"universal/framework/define"
	"universal/library/baselib/safe"
	"universal/library/mlog"
)

type RouteInfo struct {
	uid        uint64
	updateTime int64
	table      define.IRouter
}

type RouterMgr struct {
	newFunc define.NewRouterFunc  // 创建路由表
	mutex   *sync.RWMutex         // 互斥锁
	routers map[uint64]*RouteInfo // 路由表
}

func NewRouterMgr(f define.NewRouterFunc) *RouterMgr {
	return &RouterMgr{
		newFunc: f,
		mutex:   new(sync.RWMutex),
		routers: make(map[uint64]*RouteInfo),
	}
}

func (r *RouterMgr) Get(id uint64, nodeType int32) int32 {
	if val := r.get(id); val != nil {
		return val.table.Get(nodeType)
	}
	return 0
}

func (r *RouterMgr) Update(id uint64, nodeType, nodeId int32) {
	item := r.getOrNew(id)
	item.updateTime = time.Now().Unix()
	item.table.Update(nodeType, nodeId)
}

func (r *RouterMgr) get(uid uint64) *RouteInfo {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if val, ok := r.routers[uid]; ok {
		return val
	}
	return nil
}

func (r *RouterMgr) getOrNew(uid uint64) *RouteInfo {
	if val := r.get(uid); val != nil {
		return val
	}
	item := &RouteInfo{uid: uid, table: r.newFunc()}
	r.mutex.Lock()
	r.routers[uid] = item
	r.mutex.Unlock()
	return item
}

func (r *RouterMgr) getExpires(ttl int64) (rets []*RouteInfo) {
	now := time.Now().Unix()
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for _, vv := range r.routers {
		if vv.updateTime+ttl <= now {
			rets = append(rets, vv)
		}
	}
	return
}

func (r *RouterMgr) Expire(ttl int64) {
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
				delete(r.routers, val.uid)
			}
			r.mutex.Unlock()
		}
	})
}
