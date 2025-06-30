package router

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

func (r *Router) SetUpdateTime(now int64) {
	atomic.StoreInt64(&r.updateTime, now)
}

func (r *Router) GetData() *pb.Router {
	return r.Router
}

func (r *Router) SetData(data *pb.Router) {
	r.Router = data
}

func (r *Router) Get(nodeType pb.NodeType) uint32 {
	switch nodeType {
	case pb.NodeType_NodeTypeBuild:
		return atomic.LoadUint32(&r.Build)
	case pb.NodeType_NodeTypeDb:
		return atomic.LoadUint32(&r.Db)
	case pb.NodeType_NodeTypeGame:
		return atomic.LoadUint32(&r.Game)
	case pb.NodeType_NodeTypeGate:
		return atomic.LoadUint32(&r.Gate)
	case pb.NodeType_NodeTypeRoom:
		return atomic.LoadUint32(&r.Room)
	case pb.NodeType_NodeTypeMatch:
		return atomic.LoadUint32(&r.Match)
	case pb.NodeType_NodeTypeGm:
		return atomic.LoadUint32(&r.Gm)
	}
	return 0
}

func (r *Router) Set(nodeType pb.NodeType, nodeId uint32) {
	switch nodeType {
	case pb.NodeType_NodeTypeBuild:
		atomic.StoreUint32(&r.Build, nodeId)
	case pb.NodeType_NodeTypeDb:
		atomic.StoreUint32(&r.Db, nodeId)
	case pb.NodeType_NodeTypeGame:
		atomic.StoreUint32(&r.Game, nodeId)
	case pb.NodeType_NodeTypeGate:
		atomic.StoreUint32(&r.Gate, nodeId)
	case pb.NodeType_NodeTypeRoom:
		atomic.StoreUint32(&r.Room, nodeId)
	case pb.NodeType_NodeTypeMatch:
		atomic.StoreUint32(&r.Match, nodeId)
	case pb.NodeType_NodeTypeGm:
		atomic.StoreUint32(&r.Gm, nodeId)
	}
	r.SetUpdateTime(time.Now().Unix())
}
