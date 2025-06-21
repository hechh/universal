package router

/*
import (
	"sync/atomic"
	"time"
	"universal/common/pb"
)

type Router struct {
	*pb.Router
	updateTime int64
}

func (r *Router) GetUpdateTime() int64 {
	return atomic.LoadInt64(&r.updateTime)
}

func (r *Router) GetData() *pb.Router {
	return r.Router
}

func (r *Router) SetData(rr *pb.Router) {
	r.Set(pb.NodeType_NodeTypeBuilder, rr.Builder)
	r.Set(pb.NodeType_NodeTypeDb, rr.Db)
	r.Set(pb.NodeType_NodeTypeGame, rr.Game)
	r.Set(pb.NodeType_NodeTypeGm, rr.Gm)
	r.Set(pb.NodeType_NodeTypeMatch, rr.Match)
	r.Set(pb.NodeType_NodeTypeRoom, rr.Room)
	r.Set(pb.NodeType_NodeTypeGate, rr.Gate)
}

func (r *Router) Get(nodeType pb.NodeType) int32 {
	switch nodeType {
	case pb.NodeType_NodeTypeGame:
		return atomic.LoadInt32(&r.Game)
	case pb.NodeType_NodeTypeRoom:
		return atomic.LoadInt32(&r.Room)
	case pb.NodeType_NodeTypeMatch:
		return atomic.LoadInt32(&r.Match)
	case pb.NodeType_NodeTypeDb:
		return atomic.LoadInt32(&r.Db)
	case pb.NodeType_NodeTypeBuilder:
		return atomic.LoadInt32(&r.Builder)
	case pb.NodeType_NodeTypeGm:
		return atomic.LoadInt32(&r.Gm)
	default:
		return atomic.LoadInt32(&r.Gate)
	}
}

func (r *Router) Set(nodeType pb.NodeType, nodeId int32) {
	if nodeId <= 0 {
		return
	}
	switch nodeType {
	case pb.NodeType_NodeTypeGame:
		if !atomic.CompareAndSwapInt32(&r.Game, nodeId, nodeId) {
			atomic.StoreInt32(&r.Game, nodeId)
			atomic.StoreInt64(&r.updateTime, time.Now().Unix())
		}
	case pb.NodeType_NodeTypeRoom:
		if !atomic.CompareAndSwapInt32(&r.Room, nodeId, nodeId) {
			atomic.StoreInt32(&r.Room, nodeId)
			atomic.StoreInt64(&r.updateTime, time.Now().Unix())
		}
	case pb.NodeType_NodeTypeMatch:
		if !atomic.CompareAndSwapInt32(&r.Match, nodeId, nodeId) {
			atomic.StoreInt32(&r.Match, nodeId)
			atomic.StoreInt64(&r.updateTime, time.Now().Unix())
		}
	case pb.NodeType_NodeTypeDb:
		if !atomic.CompareAndSwapInt32(&r.Db, nodeId, nodeId) {
			atomic.StoreInt32(&r.Db, nodeId)
			atomic.StoreInt64(&r.updateTime, time.Now().Unix())
		}
	case pb.NodeType_NodeTypeBuilder:
		if !atomic.CompareAndSwapInt32(&r.Builder, nodeId, nodeId) {
			atomic.StoreInt32(&r.Builder, nodeId)
			atomic.StoreInt64(&r.updateTime, time.Now().Unix())
		}
	case pb.NodeType_NodeTypeGm:
		if !atomic.CompareAndSwapInt32(&r.Gm, nodeId, nodeId) {
			atomic.StoreInt32(&r.Gm, nodeId)
			atomic.StoreInt64(&r.updateTime, time.Now().Unix())
		}
	default:
		if !atomic.CompareAndSwapInt32(&r.Gate, nodeId, nodeId) {
			atomic.StoreInt32(&r.Gate, nodeId)
			atomic.StoreInt64(&r.updateTime, time.Now().Unix())
		}
	}
}
*/
