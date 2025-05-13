package router

import (
	"universal/common/pb"
)

type Router struct {
	*pb.Router
	updateTime int64
}

func (d *Router) Get(nodeType pb.NodeType) int32 {
	switch nodeType {
	case pb.NodeType_Gate:
		return d.Gate
	case pb.NodeType_Game:
		return d.Game
	case pb.NodeType_Gm:
		return d.Gm
	case pb.NodeType_Db:
		return d.Db
	case pb.NodeType_Room:
		return d.Room
	case pb.NodeType_Match:
		return d.Match
	}
	return d.Gate
}

func (d *Router) Set(nodeType pb.NodeType, nodeId int32) {
	switch nodeType {
	case pb.NodeType_Gm:
		d.Gm = nodeId
	case pb.NodeType_Gate:
		d.Gate = nodeId
	case pb.NodeType_Game:
		d.Game = nodeId
	case pb.NodeType_Db:
		d.Db = nodeId
	case pb.NodeType_Room:
		d.Room = nodeId
	case pb.NodeType_Match:
		d.Match = nodeId
	}
}
