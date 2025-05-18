package router

import (
	"poker_server/common/pb"
	"sync/atomic"
	"time"
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
	case pb.NodeType_Room:
		return atomic.LoadInt32(&d.Room)
	case pb.NodeType_Match:
		return atomic.LoadInt32(&d.Match)
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
	case pb.NodeType_Room:
		if !atomic.CompareAndSwapInt32(&d.Room, nodeId, nodeId) {
			atomic.StoreInt32(&d.Room, nodeId)
			flag = true
		}
	case pb.NodeType_Match:
		if !atomic.CompareAndSwapInt32(&d.Match, nodeId, nodeId) {
			atomic.StoreInt32(&d.Match, nodeId)
			flag = true
		}
	}
	if flag {
		atomic.StoreInt64(&d.updateTime, time.Now().Unix())
	}
}
