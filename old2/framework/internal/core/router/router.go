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

func (d *Router) GetData() *pb.Router {
	return d.Router
}

func (d *Router) IsExpire(now, ttl int64) bool {
	return now >= atomic.LoadInt64(&d.updateTime)+ttl
}

func (d *Router) Get(nodeType pb.NodeType) int32 {
	switch nodeType {
	case pb.NodeType_Gate:
		return atomic.LoadInt32(&d.Gate)
	case pb.NodeType_Db:
		return atomic.LoadInt32(&d.Db)
	case pb.NodeType_Game:
		return atomic.LoadInt32(&d.Game)
	}
	return d.Gate
}

func (d *Router) Set(nodeType pb.NodeType, nodeId int32) {
	flag := false
	switch nodeType {
	case pb.NodeType_Gate:
		if !atomic.CompareAndSwapInt32(&d.Gate, nodeId, nodeId) {
			atomic.StoreInt32(&d.Gate, nodeId)
			flag = true
		}
	case pb.NodeType_Db:
		if !atomic.CompareAndSwapInt32(&d.Db, nodeId, nodeId) {
			atomic.StoreInt32(&d.Db, nodeId)
			flag = true
		}
	case pb.NodeType_Game:
		if !atomic.CompareAndSwapInt32(&d.Game, nodeId, nodeId) {
			atomic.StoreInt32(&d.Game, nodeId)
			flag = true
		}
	}
	if flag {
		atomic.StoreInt64(&d.updateTime, time.Now().Unix())
	}
}

func (d *Router) Update(info *pb.Router) {
	if info == nil {
		return
	}
	// 更新路由信息
	if info.Gate > 0 {
		d.Set(pb.NodeType_Gate, info.Gate)
	}
	if info.Db > 0 {
		d.Set(pb.NodeType_Db, info.Db)
	}
	if info.Game > 0 {
		d.Set(pb.NodeType_Game, info.Game)
	}
}
