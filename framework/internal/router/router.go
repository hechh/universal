package router

import "sync"

type RouteInfo struct {
	Gate int32 // 网关id
	Db   int32 // 数据库id
	Game int32 // 游戏服务
}

type Router struct {
	mutex   *sync.RWMutex // 互斥锁
	players map[uint64]*RouteInfo
}

func NewRouter(types map[int32]int32) *Router {
	return &Router{
		mutex:   new(sync.RWMutex),
		players: make(map[uint64]*RouteInfo),
	}
}

func (r *Router) get(uid uint64) *RouteInfo {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if val, ok := r.players[uid]; ok {
		return val
	}
	return nil
}

func (r *Router) getOrNew(uid uint64) *RouteInfo {
	if val := r.get(uid); val != nil {
		return val
	}
	item := &RouteInfo{}
	r.mutex.Lock()
	r.players[uid] = item
	r.mutex.Unlock()
	return item
}

func (r *Router) Get(id uint64, nodeType int32) int32 {
	if val := r.get(id); val != nil {
		switch nodeType {

		}
	}
	return 0
}
