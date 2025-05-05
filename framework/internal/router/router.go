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
	timestamp int64             // 更新时间
	table     *define.RouteInfo // 节点
}

type Router struct {
	mutex   *sync.RWMutex         // 互斥锁
	exit    chan struct{}         // 退出信号
	routers map[uint64]*RouteInfo // 路由信息
}

func NewRouter() *Router {
	return &Router{
		mutex:   new(sync.RWMutex),
		exit:    make(chan struct{}),
		routers: make(map[uint64]*RouteInfo),
	}
}

func (r *Router) Get(id uint64) *define.RouteInfo {
	r.mutex.RLock()
	val, ok := r.routers[id]
	r.mutex.RUnlock()
	if ok {
		return val.table
	}
	return nil
}

func (r *Router) Update(id uint64, tab *define.RouteInfo) {
	if val := r.Get(id); val != nil {
		atomic.StoreUint32(&val.Gate, tab.Gate)
		atomic.StoreUint32(&val.Db, tab.Db)
		atomic.StoreUint32(&val.Game, tab.Game)
		atomic.StoreUint32(&val.Tool, tab.Tool)
		atomic.StoreUint32(&val.Rank, tab.Rank)
		return
	}
	// 更新路由
	r.mutex.Lock()
	r.routers[id] = &RouteInfo{timestamp: time.Now().Unix(), table: tab}
	r.mutex.Unlock()
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
