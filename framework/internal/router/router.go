package router

import (
	"sync"
	"sync/atomic"
	"time"
	"universal/framework/define"
	"universal/library/baselib/safe"
	"universal/library/mlog"
)

type RouteInfo struct {
	define.ITable
	timestamp int64 // 更新时间
}

type Router struct {
	mutex    *sync.RWMutex         // 互斥锁
	exit     chan struct{}         // 退出信号
	newTable func() define.ITable  // 创建路由表
	routers  map[uint64]*RouteInfo // 路由信息
}

func NewRouter(f func() define.ITable) *Router {
	return &Router{
		mutex:    new(sync.RWMutex),
		exit:     make(chan struct{}),
		newTable: f,
		routers:  make(map[uint64]*RouteInfo),
	}
}

func (r *Router) Get(id uint64) define.ITable {
	r.mutex.RLock()
	val, ok := r.routers[id]
	r.mutex.RUnlock()
	if !ok {
		r.mutex.Lock()
		val = &RouteInfo{timestamp: time.Now().Unix(), ITable: r.newTable()}
		r.routers[id] = val
		r.mutex.Unlock()
	}
	return val.ITable
}

func (r *Router) Update(id uint64, tab define.ITable) {
	val := r.Get(id)
	for i := uint32(define.NodeTypeBegin) + 1; i < uint32(define.NodeTypeMax); i++ {
		if tab.Get(i) > 0 && val.Get(i) != tab.Get(i) {
			val.Set(i, tab.Get(i))
		}
	}
}

func (r *Router) Close() error {
	r.exit <- struct{}{}
	return nil
}

func (r *Router) Expire(ttl int64) {
	safe.SafeGo(mlog.Fatal, func() {
		tt := time.NewTicker(time.Duration(ttl) * time.Second)
		defer tt.Stop()
		// 定时清理
		for {
			select {
			case <-r.exit:
				return
			case <-tt.C:
				now := time.Now().Unix()
				// 获取过期节点
				ids := []uint64{}
				r.mutex.RLock()
				for kk, vv := range r.routers {
					if now >= atomic.LoadInt64(&vv.timestamp)+ttl {
						ids = append(ids, kk)
					}
				}
				r.mutex.RUnlock()
				// 删除节点
				r.mutex.Lock()
				for _, val := range ids {
					delete(r.routers, val)
				}
				r.mutex.Unlock()
			}
		}
	})
}
