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

func (r *Router) SetUpdateTime() {
	atomic.StoreInt64(&r.updateTime, time.Now().Unix())
}

func (r *Router) Get(nodeType pb.NodeType) (ret int32) {
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

func (r *Router) GetData() *pb.Router {
	return &pb.Router{
		Build: r.Get(pb.NodeType_Build),
		Room:  r.Get(pb.NodeType_Room),
		Match: r.Get(pb.NodeType_Match),
		Db:    r.Get(pb.NodeType_Db),
		Game:  r.Get(pb.NodeType_Game),
		Gate:  r.Get(pb.NodeType_Gate),
		Gm:    r.Get(pb.NodeType_Gm),
	}
}

func (r *Router) Set(nodeType pb.NodeType, nodeId int32) define.IRouter {
	if nodeId > 0 {
		r.SetUpdateTime()
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

func (r *Router) SetData(data *pb.Router) define.IRouter {
	if data != nil {
		r.SetUpdateTime()
		atomic.StoreInt32(&r.Build, data.Build)
		atomic.StoreInt32(&r.Room, data.Room)
		atomic.StoreInt32(&r.Match, data.Match)
		atomic.StoreInt32(&r.Db, data.Db)
		atomic.StoreInt32(&r.Game, data.Game)
		atomic.StoreInt32(&r.Gate, data.Gate)
		atomic.StoreInt32(&r.Gm, data.Gm)
	}
	return r
}
