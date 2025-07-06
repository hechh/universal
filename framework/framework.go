package framework

import (
	"universal/common/pb"
	"universal/framework/cluster"
	"universal/framework/rpc"
	"universal/library/util"

	"github.com/spf13/cast"
)

func SwapToGate(head *pb.Head, actorFunc string, actorId uint64) *pb.Head {
	head.Dst, head.Src = head.Src, head.Dst
	head.Dst = rpc.NewNodeRouter(pb.NodeType_NodeTypeGate, actorFunc, actorId)
	return head
}
func SwapToGame(head *pb.Head, actorFunc string, actorId uint64) *pb.Head {
	head.Dst, head.Src = head.Src, head.Dst
	head.Dst = rpc.NewNodeRouter(pb.NodeType_NodeTypeGame, actorFunc, actorId)
	return head
}
func SwapToDb(head *pb.Head, actorFunc string, actorId uint64) *pb.Head {
	head.Dst, head.Src = head.Src, head.Dst
	head.Dst = rpc.NewNodeRouter(pb.NodeType_NodeTypeDb, actorFunc, actorId)
	return head
}
func SwapToGm(head *pb.Head, actorFunc string, actorId uint64) *pb.Head {
	head.Dst, head.Src = head.Src, head.Dst
	head.Dst = rpc.NewNodeRouter(pb.NodeType_NodeTypeGm, actorFunc, actorId)
	return head
}
func SwapToRoom(head *pb.Head, actorFunc string, actorId uint64) *pb.Head {
	head.Dst, head.Src = head.Src, head.Dst
	head.Dst = rpc.NewNodeRouter(pb.NodeType_NodeTypeRoom, actorFunc, actorId)
	return head
}
func SwapToMatch(head *pb.Head, actorFunc string, actorId uint64) *pb.Head {
	head.Dst, head.Src = head.Src, head.Dst
	head.Dst = rpc.NewNodeRouter(pb.NodeType_NodeTypeMatch, actorFunc, actorId)
	return head
}
func SwapToBuild(head *pb.Head, actorFunc string, actorId uint64) *pb.Head {
	head.Dst, head.Src = head.Src, head.Dst
	head.Dst = rpc.NewNodeRouter(pb.NodeType_NodeTypeBuild, actorFunc, actorId)
	return head
}

func SendToGate(head *pb.Head, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return sendTo(head, pb.NodeType_NodeTypeGate, actorFunc, actorId, srcs...)
}
func SendToGame(head *pb.Head, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return sendTo(head, pb.NodeType_NodeTypeGame, actorFunc, actorId, srcs...)
}
func SendToDb(head *pb.Head, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return sendTo(head, pb.NodeType_NodeTypeDb, actorFunc, actorId, srcs...)
}
func SendToGm(head *pb.Head, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return sendTo(head, pb.NodeType_NodeTypeGm, actorFunc, actorId, srcs...)
}
func SendToRoom(head *pb.Head, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return sendTo(head, pb.NodeType_NodeTypeRoom, actorFunc, actorId, srcs...)
}
func SendToMatch(head *pb.Head, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return sendTo(head, pb.NodeType_NodeTypeMatch, actorFunc, actorId, srcs...)
}
func SendToBuild(head *pb.Head, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return sendTo(head, pb.NodeType_NodeTypeBuild, actorFunc, actorId, srcs...)
}

func sendTo(head *pb.Head, nt pb.NodeType, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	act := cast.ToString(util.Index[interface{}](srcs, 0, ""))
	actId := cast.ToUint64(util.Index[interface{}](srcs, 1, 0))
	src := rpc.NewNodeRouter(cluster.GetSelf().Type, act, actId)
	return &pb.Head{
		Src:       util.Or[*pb.NodeRouter](src != nil, src, head.Src),
		Dst:       rpc.NewNodeRouter(nt, actorFunc, actorId),
		Uid:       head.Uid,
		Seq:       head.Seq,
		Cmd:       head.Cmd,
		Reference: head.Reference,
		Reply:     head.Reply,
		SendType:  head.SendType,
	}
}

func NewToGate(uid uint64, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return newTo(uid, pb.NodeType_NodeTypeGate, actorFunc, actorId, srcs...)
}
func NewToDb(uid uint64, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return newTo(uid, pb.NodeType_NodeTypeDb, actorFunc, actorId, srcs...)
}
func NewToGame(uid uint64, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return newTo(uid, pb.NodeType_NodeTypeGame, actorFunc, actorId, srcs...)
}
func NewToMatch(uid uint64, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return newTo(uid, pb.NodeType_NodeTypeMatch, actorFunc, actorId, srcs...)
}
func NewToRoom(uid uint64, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return newTo(uid, pb.NodeType_NodeTypeRoom, actorFunc, actorId, srcs...)
}
func NewToBuild(uid uint64, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return newTo(uid, pb.NodeType_NodeTypeBuild, actorFunc, actorId, srcs...)
}
func NewToGm(uid uint64, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	return newTo(uid, pb.NodeType_NodeTypeGm, actorFunc, actorId, srcs...)
}

func newTo(uid uint64, nt pb.NodeType, actorFunc string, actorId uint64, srcs ...interface{}) *pb.Head {
	act := cast.ToString(util.Index[interface{}](srcs, 0, ""))
	actId := cast.ToUint64(util.Index[interface{}](srcs, 1, 0))
	return &pb.Head{
		Uid: uid,
		Src: rpc.NewNodeRouter(cluster.GetSelf().Type, act, actId),
		Dst: rpc.NewNodeRouter(nt, actorFunc, actorId),
	}
}
