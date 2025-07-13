package router

import (
	"sync/atomic"
	"time"
	"universal/common/pb"
	"universal/framework/define"
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
	r.SetUpdateTime(time.Now().Unix())
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

func (r *Router) SetData(data *pb.Router) define.IRouter {
	if data != nil {
		r.Set(pb.NodeType_Build, data.Build)
		r.Set(pb.NodeType_Room, data.Room)
		r.Set(pb.NodeType_Match, data.Match)
		r.Set(pb.NodeType_Db, data.Db)
		r.Set(pb.NodeType_Game, data.Game)
		r.Set(pb.NodeType_Gate, data.Gate)
		r.Set(pb.NodeType_Gm, data.Gm)
	}
	return r
}

func (r *Router) Get(nodeType pb.NodeType) (ret int32) {
	r.SetUpdateTime(time.Now().Unix())
	switch nodeType {
	case pb.NodeType_Build:
		ret = atomic.LoadInt32(&r.Build)
	case pb.NodeType_Db:
		ret = atomic.LoadInt32(&r.Db)
	case pb.NodeType_Game:
		ret = atomic.LoadInt32(&r.Game)
	case pb.NodeType_Gate:
		ret = atomic.LoadInt32(&r.Gate)
	case pb.NodeType_Room:
		ret = atomic.LoadInt32(&r.Room)
	case pb.NodeType_Match:
		ret = atomic.LoadInt32(&r.Match)
	case pb.NodeType_Gm:
		ret = atomic.LoadInt32(&r.Gm)
	}
	return
}

func (r *Router) Set(nodeType pb.NodeType, nodeId int32) define.IRouter {
	if nodeId > 0 {
		r.SetUpdateTime(time.Now().Unix())
		switch nodeType {
		case pb.NodeType_Build:
			atomic.StoreInt32(&r.Build, nodeId)
		case pb.NodeType_Db:
			atomic.StoreInt32(&r.Db, nodeId)
		case pb.NodeType_Game:
			atomic.StoreInt32(&r.Game, nodeId)
		case pb.NodeType_Gate:
			atomic.StoreInt32(&r.Gate, nodeId)
		case pb.NodeType_Room:
			atomic.StoreInt32(&r.Room, nodeId)
		case pb.NodeType_Match:
			atomic.StoreInt32(&r.Match, nodeId)
		case pb.NodeType_Gm:
			atomic.StoreInt32(&r.Gm, nodeId)
		}
	}
	return r
}
