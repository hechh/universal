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
	NodeType uint32
}

type RouteInfo struct {
	updateTime int64
	node       define.INode // 节点
}

type Router struct {
	mutex   *sync.RWMutex        // 互斥锁
	exit    chan struct{}        // 退出信号
	routers map[index]*RouteInfo // 路由信息
}

func NewRouter() *Router {
	return &Router{
		mutex:   new(sync.RWMutex),
		exit:    make(chan struct{}),
		routers: make(map[index]*RouteInfo),
	}
}

func (r *Router) get(id uint64, nodeType uint32) *RouteInfo {
	r.mutex.RLock()
	val, ok := r.routers[index{id, nodeType}]
	r.mutex.RUnlock()
	if ok {
		return val
	}
	return nil
}
func (r *Router) Get(id uint64, nodeType uint32) define.INode {
	if val := r.get(id, nodeType); val != nil {
		return val.node
	}
	return nil
}

func (r *Router) Update(id uint64, node define.INode) {
	if val := r.get(id, node.GetType()); val != nil {
		atomic.StoreInt64(&val.updateTime, time.Now().Unix())
		val.node = node
		return
	}

	// 更新路由
	item := &RouteInfo{updateTime: time.Now().Unix(), node: node}
	r.mutex.Lock()
	r.routers[index{id, node.GetType()}] = item
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
		}
	})
}
