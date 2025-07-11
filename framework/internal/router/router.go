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
	return &pb.Router{
		Build: atomic.LoadInt32(&r.Build),
		Room:  atomic.LoadInt32(&r.Room),
		Match: atomic.LoadInt32(&r.Match),
		Db:    atomic.LoadInt32(&r.Db),
		Game:  atomic.LoadInt32(&r.Game),
		Gate:  atomic.LoadInt32(&r.Gate),
		Gm:    atomic.LoadInt32(&r.Gm),
	}
}

func (r *Router) SetData(data *pb.Router) domain.IRouter {
	r.Set(pb.NodeType_Build, data.Build)
	r.Set(pb.NodeType_Room, data.Room)
	r.Set(pb.NodeType_Match, data.Match)
	r.Set(pb.NodeType_Db, data.Db)
	r.Set(pb.NodeType_Game, data.Game)
	r.Set(pb.NodeType_Gate, data.Gate)
	r.Set(pb.NodeType_Gm, data.Gm)
	return r
}

func (r *Router) Get(nodeType pb.NodeType) int32 {
	switch nodeType {
	case pb.NodeType_Build:
		return atomic.LoadInt32(&r.Build)
	case pb.NodeType_Db:
		return atomic.LoadInt32(&r.Db)
	case pb.NodeType_Game:
		return atomic.LoadInt32(&r.Game)
	case pb.NodeType_Gate:
		return atomic.LoadInt32(&r.Gate)
	case pb.NodeType_Room:
		return atomic.LoadInt32(&r.Room)
	case pb.NodeType_Match:
		return atomic.LoadInt32(&r.Match)
	case pb.NodeType_Gm:
		return atomic.LoadInt32(&r.Gm)
	}
	return 0
}

func (r *Router) Set(nodeType pb.NodeType, nodeId int32) domain.IRouter {
	if nodeId > 0 {
		switch nodeType {
		case pb.NodeType_Build:
			atomic.StoreInt32(&r.Build, nodeId)
			r.SetUpdateTime(time.Now().Unix())
		case pb.NodeType_Db:
			atomic.StoreInt32(&r.Db, nodeId)
			r.SetUpdateTime(time.Now().Unix())
		case pb.NodeType_Game:
			atomic.StoreInt32(&r.Game, nodeId)
			r.SetUpdateTime(time.Now().Unix())
		case pb.NodeType_Gate:
			atomic.StoreInt32(&r.Gate, nodeId)
			r.SetUpdateTime(time.Now().Unix())
		case pb.NodeType_Room:
			atomic.StoreInt32(&r.Room, nodeId)
			r.SetUpdateTime(time.Now().Unix())
		case pb.NodeType_Match:
			atomic.StoreInt32(&r.Match, nodeId)
			r.SetUpdateTime(time.Now().Unix())
		case pb.NodeType_Gm:
			atomic.StoreInt32(&r.Gm, nodeId)
			r.SetUpdateTime(time.Now().Unix())
		}
	}
	return r
}
