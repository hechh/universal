package router

import (
	"sync"
	"sync/atomic"
	"time"
	"universal/framework/domain"
	"universal/library/baselib/safe"
	"universal/library/mlog"
)

type RouteInfo struct {
	domain.IRoute
	updateTime int64
}

type RouterMgr struct {
	mutex    sync.RWMutex
	exit     chan struct{}
	newRoute func() domain.IRoute
	routes   map[uint64]*RouteInfo
}

func (d *RouterMgr) Get(routeId uint64) domain.IRoute {
	d.mutex.RLock()
	route, ok := d.routes[routeId]
	d.mutex.RUnlock()
	if ok {
		return route.IRoute
	}
	// 创建新的路由节点
	val := &RouteInfo{
		updateTime: time.Now().Unix(),
		IRoute:     d.newRoute(),
	}
	d.mutex.Lock()
	d.routes[routeId] = val
	d.mutex.Unlock()
	return val
}

func (d *RouterMgr) Set(routeId uint64, info domain.IRoute) {
	route := d.Get(routeId).(*RouteInfo)
	now := time.Now().Unix()
	for i := int32(domain.NodeTypeBegin) + 1; i < int32(domain.NodeTypeMax); i++ {
		if info.Get(i) > 0 && route.Get(i) != info.Get(i) {
			route.Set(i, info.Get(i))
			atomic.StoreInt64(&route.updateTime, now)
		}
	}
}

func (d *RouterMgr) getExpireRouteIds() (rets []uint64) {
	now := time.Now().Unix()
	d.mutex.RLock()
	for routeId, val := range d.routes {
		if now >= atomic.LoadInt64(&val.updateTime) {
			rets = append(rets, routeId)
		}
	}
	d.mutex.RUnlock()
	return
}

func (r *RouterMgr) Expire(ttl int64) {
	safe.SafeGo(mlog.Fatal, func() {
		tt := time.NewTicker(time.Duration(ttl) * time.Second)
		defer tt.Stop()
		// 定时清理
		for {
			select {
			case <-tt.C:
				// 获取过期节点
				rets := r.getExpireRouteIds()
				if len(rets) > 0 {
					// 删除节点
					r.mutex.Lock()
					for _, val := range rets {
						delete(r.routes, val)
					}
					r.mutex.Unlock()
				}
			case <-r.exit:
				return
			}
		}
	})
}

func (r *RouterMgr) Close() {
	r.exit <- struct{}{}
}
