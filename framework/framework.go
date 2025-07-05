package framework

import (
	"universal/common/pb"
	"universal/framework/cluster"
	"universal/framework/rpc"
	"universal/library/util"

	"github.com/spf13/cast"
)

func NewGateHead(uid uint64, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	act := cast.ToString(util.Index[interface{}](srcs, 0, ""))
	actId := cast.ToUint64(util.Index[interface{}](srcs, 1, 0))
	return &pb.Head{
		Uid: uid,
		Src: rpc.NewNodeRouter(cluster.GetSelf().Type, act, actId),
		Dst: rpc.NewNodeRouter(pb.NodeType_NodeTypeGate, actorFunc, actorId),
	}
}

func NewDbHead(uid uint64, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	act := cast.ToString(util.Index[interface{}](srcs, 0, ""))
	actId := cast.ToUint64(util.Index[interface{}](srcs, 1, 0))
	return &pb.Head{
		Uid: uid,
		Src: rpc.NewNodeRouter(cluster.GetSelf().Type, act, actId),
		Dst: rpc.NewNodeRouter(pb.NodeType_NodeTypeDb, actorFunc, actorId),
	}
}

func NewGameHead(uid uint64, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	act := cast.ToString(util.Index[interface{}](srcs, 0, ""))
	actId := cast.ToUint64(util.Index[interface{}](srcs, 1, 0))
	return &pb.Head{
		Uid: uid,
		Src: rpc.NewNodeRouter(cluster.GetSelf().Type, act, actId),
		Dst: rpc.NewNodeRouter(pb.NodeType_NodeTypeGame, actorFunc, actorId),
	}
}

func NewMatchHead(uid uint64, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	act := cast.ToString(util.Index[interface{}](srcs, 0, ""))
	actId := cast.ToUint64(util.Index[interface{}](srcs, 1, 0))
	return &pb.Head{
		Uid: uid,
		Src: rpc.NewNodeRouter(cluster.GetSelf().Type, act, actId),
		Dst: rpc.NewNodeRouter(pb.NodeType_NodeTypeMatch, actorFunc, actorId),
	}
}

func NewRoomHead(uid uint64, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	act := cast.ToString(util.Index[interface{}](srcs, 0, ""))
	actId := cast.ToUint64(util.Index[interface{}](srcs, 1, 0))
	return &pb.Head{
		Uid: uid,
		Src: rpc.NewNodeRouter(cluster.GetSelf().Type, act, actId),
		Dst: rpc.NewNodeRouter(pb.NodeType_NodeTypeRoom, actorFunc, actorId),
	}
}

func NewBuildHead(uid uint64, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	act := cast.ToString(util.Index[interface{}](srcs, 0, ""))
	actId := cast.ToUint64(util.Index[interface{}](srcs, 1, 0))
	return &pb.Head{
		Uid: uid,
		Src: rpc.NewNodeRouter(cluster.GetSelf().Type, act, actId),
		Dst: rpc.NewNodeRouter(pb.NodeType_NodeTypeBuild, actorFunc, actorId),
	}
}

func NewGmHead(uid uint64, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	act := cast.ToString(util.Index[interface{}](srcs, 0, ""))
	actId := cast.ToUint64(util.Index[interface{}](srcs, 1, 0))
	return &pb.Head{
		Uid: uid,
		Src: rpc.NewNodeRouter(cluster.GetSelf().Type, act, actId),
		Dst: rpc.NewNodeRouter(pb.NodeType_NodeTypeGm, actorFunc, actorId),
	}
}
