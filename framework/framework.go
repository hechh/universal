package framework

import (
	"universal/common/pb"
	"universal/common/yaml"
	"universal/framework/cluster"
	"universal/framework/internal/funcs"
	"universal/library/pprof"
	"universal/library/util"

	"github.com/spf13/cast"
)

func Init(cfg *yaml.Config, srvCfg *yaml.NodeConfig, nn *pb.Node) error {
	if err := cluster.Init(cfg, srvCfg, nn); err != nil {
		return err
	}
	funcs.Init(nn, cluster.SendResponse)
	pprof.Init("localhost", srvCfg.Port+10000)
	return nil
}

func CopyToGate(head *pb.Head, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return CopyHead(head, pb.NodeType_Gate, actorFunc, actorId, srcs...)
}

func CopyToGame(head *pb.Head, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return CopyHead(head, pb.NodeType_Game, actorFunc, actorId, srcs...)
}

func CopyToDb(head *pb.Head, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return CopyHead(head, pb.NodeType_Db, actorFunc, actorId, srcs...)
}

func CopyToGm(head *pb.Head, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return CopyHead(head, pb.NodeType_Gm, actorFunc, actorId, srcs...)
}

func CopyToRoom(head *pb.Head, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return CopyHead(head, pb.NodeType_Room, actorFunc, actorId, srcs...)
}

func CopyToMatch(head *pb.Head, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return CopyHead(head, pb.NodeType_Match, actorFunc, actorId, srcs...)
}

func CopyToBuild(head *pb.Head, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return CopyHead(head, pb.NodeType_Build, actorFunc, actorId, srcs...)
}

func CopyHead(head *pb.Head, nt pb.NodeType, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	act := cast.ToString(util.Index[interface{}](srcs, 0, ""))
	actId := cast.ToUint64(util.Index[interface{}](srcs, 1, 0))
	return &pb.Head{
		Src: cluster.NewNodeRouter(act, actId),
		Dst: &pb.NodeRouter{
			NodeType:  nt,
			ActorFunc: util.GetCrc32(actorFunc),
			ActorId:   actorId,
		},
		Uid:   head.Uid,
		Seq:   head.Seq,
		Cmd:   head.Cmd,
		Reply: head.Reply,
	}
}

func NewToGate(uid uint64, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return NewHead(uid, pb.NodeType_Gate, actorFunc, actorId, srcs...)
}

func NewToDb(uid uint64, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return NewHead(uid, pb.NodeType_Db, actorFunc, actorId, srcs...)
}

func NewToGame(uid uint64, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return NewHead(uid, pb.NodeType_Game, actorFunc, actorId, srcs...)
}

func NewToMatch(uid uint64, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return NewHead(uid, pb.NodeType_Match, actorFunc, actorId, srcs...)
}

func NewToRoom(uid uint64, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return NewHead(uid, pb.NodeType_Room, actorFunc, actorId, srcs...)
}

func NewToBuild(uid uint64, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return NewHead(uid, pb.NodeType_Build, actorFunc, actorId, srcs...)
}

func NewToGm(uid uint64, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return NewHead(uid, pb.NodeType_Gm, actorFunc, actorId, srcs...)
}

func NewHead(uid uint64, nt pb.NodeType, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	srcFunc := cast.ToString(util.Index[interface{}](srcs, 0, ""))
	srcId := cast.ToUint64(util.Index[interface{}](srcs, 1, 0))
	return &pb.Head{
		Uid: uid,
		Src: cluster.NewNodeRouter(srcFunc, srcId),
		Dst: &pb.NodeRouter{
			NodeType:  nt,
			ActorFunc: util.GetCrc32(actorFunc),
			ActorId:   actorId,
		},
	}
}
