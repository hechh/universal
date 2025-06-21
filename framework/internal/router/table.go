package router

import (
	"sync/atomic"
	"time"
	"universal/common/pb"
)

type Router struct {
	updateTime int64
	build      uint32
	db         uint32
	game       uint32
	gate       uint32
	room       uint32
	match      uint32
	gm         uint32
}

func (r *Router) Get(nodeType pb.NodeType) uint32 {
	switch nodeType {
	case pb.NodeType_NodeTypeBuild:
		return atomic.LoadUint32(&r.build)
	case pb.NodeType_NodeTypeDb:
		return atomic.LoadUint32(&r.db)
	case pb.NodeType_NodeTypeGame:
		return atomic.LoadUint32(&r.game)
	case pb.NodeType_NodeTypeGate:
		return atomic.LoadUint32(&r.gate)
	case pb.NodeType_NodeTypeRoom:
		return atomic.LoadUint32(&r.room)
	case pb.NodeType_NodeTypeMatch:
		return atomic.LoadUint32(&r.match)
	case pb.NodeType_NodeTypeGm:
		return atomic.LoadUint32(&r.gm)
	}
	return 0
}

func (r *Router) Set(nodeType pb.NodeType, nodeId uint32) {
	switch nodeType {
	case pb.NodeType_NodeTypeBuild:
		if !atomic.CompareAndSwapUint32(&r.build, nodeId, nodeId) {
			atomic.StoreUint32(&r.build, nodeId)
			atomic.StoreInt64(&r.updateTime, time.Now().Unix())
		}
	case pb.NodeType_NodeTypeDb:
		if !atomic.CompareAndSwapUint32(&r.db, nodeId, nodeId) {
			atomic.StoreUint32(&r.db, nodeId)
			atomic.StoreInt64(&r.updateTime, time.Now().Unix())
		}
	case pb.NodeType_NodeTypeGame:
		if !atomic.CompareAndSwapUint32(&r.game, nodeId, nodeId) {
			atomic.StoreUint32(&r.game, nodeId)
			atomic.StoreInt64(&r.updateTime, time.Now().Unix())
		}
	case pb.NodeType_NodeTypeGate:
		if !atomic.CompareAndSwapUint32(&r.gate, nodeId, nodeId) {
			atomic.StoreUint32(&r.gate, nodeId)
			atomic.StoreInt64(&r.updateTime, time.Now().Unix())
		}
	case pb.NodeType_NodeTypeRoom:
		if !atomic.CompareAndSwapUint32(&r.room, nodeId, nodeId) {
			atomic.StoreUint32(&r.room, nodeId)
			atomic.StoreInt64(&r.updateTime, time.Now().Unix())
		}
	case pb.NodeType_NodeTypeMatch:
		if !atomic.CompareAndSwapUint32(&r.match, nodeId, nodeId) {
			atomic.StoreUint32(&r.match, nodeId)
			atomic.StoreInt64(&r.updateTime, time.Now().Unix())
		}
	case pb.NodeType_NodeTypeGm:
		if !atomic.CompareAndSwapUint32(&r.gm, nodeId, nodeId) {
			atomic.StoreUint32(&r.gm, nodeId)
			atomic.StoreInt64(&r.updateTime, time.Now().Unix())
		}
	}
}

/*
type Table struct {
	mutex   sync.RWMutex
	routers map[uint64]domain.IRouter
	exit    chan struct{}
}

func NewTable() *Table {
	return &Table{
		routers: make(map[uint64]domain.IRouter),
		exit:    make(chan struct{}),
	}
}

func (d *Table) Get(id uint64) domain.IRouter {
	d.mutex.RLock()
	rr, ok := d.routers[id]
	d.mutex.RUnlock()
	if ok {
		return rr
	}

	val := &Router{Router: &pb.Router{}, updateTime: time.Now().Unix()}
	d.mutex.Lock()
	d.routers[id] = val
	d.mutex.Unlock()
	return val
}

func (d *Table) Close() error {
	close(d.exit)
	return nil
}

func (d *Table) Expire(ttl int64) {
	safe.Go(func() {
		tt := time.NewTicker(time.Duration(ttl/2) * time.Second)
		defer tt.Stop()

		for {
			select {
			case <-tt.C:
				if keys := d.getExpires(ttl); len(keys) > 0 {
					d.mutex.Lock()
					for _, k := range keys {
						delete(d.routers, k)
					}
					d.mutex.Unlock()
				}
			case <-d.exit:
				return
			}
		}
	})
}

func (d *Table) getExpires(ttl int64) (keys []uint64) {
	now := time.Now().Unix()
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	for key, rr := range d.routers {
		if rr.GetUpdateTime()+ttl <= now {
			keys = append(keys, key)
		}
	}
	return
}
*/
