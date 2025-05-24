package router

import (
	"sync"
	"sync/atomic"
	"time"
	"universal/common/pb"
	"universal/framework/library/async"
	"universal/framework/library/mlog"
)

type RouterMgr struct {
	expire  int64
	mutex   sync.RWMutex
	exit    chan struct{}
	routers map[uint64]*Router
}

func NewRouterMgr(ttl int64) *RouterMgr {
	mgr := &RouterMgr{
		expire:  ttl,
		exit:    make(chan struct{}),
		routers: make(map[uint64]*Router),
	}
	async.SafeGo(mlog.Fatal, mgr.run)
	return mgr
}

func (d *RouterMgr) Get(routeId uint64) *pb.Router {
	d.mutex.RLock()
	router, ok := d.routers[routeId]
	d.mutex.RUnlock()
	if ok {
		return router.Router
	}
	// 创建新的路由节点
	val := &Router{
		updateTime: time.Now().Unix(),
		Router:     &pb.Router{},
	}
	d.mutex.Lock()
	d.routers[routeId] = val
	d.mutex.Unlock()
	return val.Router
}

func (d *RouterMgr) Set(routeId uint64, info *pb.Router) {
	route := &Router{Router: d.Get(routeId)}

	iin := &Router{Router: info}

	now := time.Now().Unix()

	for i := pb.NodeType_Begin + 1; i < pb.NodeType_End; i++ {
		if iin.Get(i) > 0 && route.Get(i) != iin.Get(i) {
			route.Set(i, iin.Get(i))
			atomic.StoreInt64(&route.updateTime, now)
		}
	}
}

func (r *RouterMgr) Close() {
	r.exit <- struct{}{}
}

func (r *RouterMgr) run() {
	tt := time.NewTicker(time.Duration(r.expire) * time.Second)
	defer tt.Stop()

	// 定时清理
	for {
		select {
		case <-tt.C:
			// 获取过期节点
			rets := []uint64{}
			now := time.Now().Unix()
			r.mutex.RLock()
			for routeId, val := range r.routers {
				if now >= atomic.LoadInt64(&val.updateTime) {
					rets = append(rets, routeId)
				}
			}
			r.mutex.RUnlock()

			// 获取过期节点
			if len(rets) > 0 {
				// 删除节点
				r.mutex.Lock()
				for _, val := range rets {
					delete(r.routers, val)
				}
				r.mutex.Unlock()
			}
		case <-r.exit:
			return
		}
	}
}
