package router

import (
	"sync"
	"sync/atomic"
	"time"
	"universal/framework/define"
	"universal/library/baselib/safe"
	"universal/library/mlog"
)

type index struct {
	Id       uint64
	NodeType int32
}

type RouteInfo struct {
	updateTime int64
	node       atomic.Pointer[define.INode] // 节点
}

type Router struct {
	mutex   *sync.RWMutex        // 互斥锁
	routers map[index]*RouteInfo // 路由信息
}

func NewRouter() *Router {
	return &Router{
		mutex:   new(sync.RWMutex),
		routers: make(map[index]*RouteInfo),
	}
}

func (r *Router) get(id uint64, nodeType int32) *RouteInfo {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if val, ok := r.routers[index{id, nodeType}]; ok {
		return val
	}
	return nil
}
func (r *Router) Get(id uint64, nodeType int32) define.INode {
	if val := r.get(id, nodeType); val != nil {
		return *val.node.Load()
	}
	return nil
}

func (r *Router) Update(id uint64, node define.INode) {
	if val := r.get(id, node.GetType()); val != nil {
		atomic.StoreInt64(&val.updateTime, time.Now().Unix())
		val.node.Store(&node)
		return
	}

	item := &RouteInfo{updateTime: time.Now().Unix()}
	item.node.Store(&node)
	// 更新路由
	r.mutex.Lock()
	r.routers[index{id, node.GetType()}] = item
	r.mutex.Unlock()
}

func (r *Router) Expire(ttl int64) {
	tt := time.NewTicker(time.Duration(ttl/2) * time.Second)
	defer tt.Stop()
	// 定时清理
	safe.SafeGo(mlog.Fatal, func() {
		for {
			<-tt.C
			now := time.Now().Unix()
			// 获取过期节点
			kks := []index{}
			r.mutex.RLock()
			for kk, vv := range r.routers {
				if now >= vv.updateTime+ttl {
					kks = append(kks, kk)
				}
			}
			r.mutex.RUnlock()
			// 删除节点
			r.mutex.Lock()
			for _, val := range kks {
				delete(r.routers, val)
			}
			r.mutex.Unlock()
		}
	})
}
