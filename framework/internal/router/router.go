package router

import (
	"sync/atomic"
	"time"
	"universal/common/pb"
	"universal/framework/domain"
)

type Router struct {
	*pb.Router
	updateTime int64
}

func (r *Router) GetUpdateTime() int64 {
	return atomic.LoadInt64(&r.updateTime)
}

func (r *Router) SetUpdateTime(now int64) domain.IRouter {
	atomic.StoreInt64(&r.updateTime, now)
	return r
}

func (r *Router) GetData() *pb.Router {
	return r.Router
}

func (r *Router) SetData(data *pb.Router) domain.IRouter {
	r.Router = data
	return r.SetUpdateTime(time.Now().Unix())
}

func (r *Router) Get(nodeType pb.NodeType) int32 {
	switch nodeType {
	case pb.NodeType_NodeTypeBuild:
		return atomic.LoadInt32(&r.Build)
	case pb.NodeType_NodeTypeDb:
		return atomic.LoadInt32(&r.Db)
	case pb.NodeType_NodeTypeGame:
		return atomic.LoadInt32(&r.Game)
	case pb.NodeType_NodeTypeGate:
		return atomic.LoadInt32(&r.Gate)
	case pb.NodeType_NodeTypeRoom:
		return atomic.LoadInt32(&r.Room)
	case pb.NodeType_NodeTypeMatch:
		return atomic.LoadInt32(&r.Match)
	case pb.NodeType_NodeTypeGm:
		return atomic.LoadInt32(&r.Gm)
	}
	return 0
}

func (r *Router) Set(nodeType pb.NodeType, nodeId int32) domain.IRouter {
	switch nodeType {
	case pb.NodeType_NodeTypeBuild:
		atomic.StoreInt32(&r.Build, nodeId)
	case pb.NodeType_NodeTypeDb:
		atomic.StoreInt32(&r.Db, nodeId)
	case pb.NodeType_NodeTypeGame:
		atomic.StoreInt32(&r.Game, nodeId)
	case pb.NodeType_NodeTypeGate:
		atomic.StoreInt32(&r.Gate, nodeId)
	case pb.NodeType_NodeTypeRoom:
		atomic.StoreInt32(&r.Room, nodeId)
	case pb.NodeType_NodeTypeMatch:
		atomic.StoreInt32(&r.Match, nodeId)
	case pb.NodeType_NodeTypeGm:
		atomic.StoreInt32(&r.Gm, nodeId)
	}
	return r.SetUpdateTime(time.Now().Unix())
}
