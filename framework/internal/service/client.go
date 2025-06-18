package service

import (
	"universal/common/pb"
	"universal/library/mlog"
	"universal/library/uerror"

	"github.com/golang/protobuf/proto"
)

func SendToClient(head *pb.Head, msg proto.Message) error {
	if head.Src == nil {
		return uerror.N(1, int32(pb.ErrorCode_NodeRouterIsNil), "%v", head)
	}
	if head.Src.ActorId <= 0 {
		return uerror.N(1, int32(pb.ErrorCode_ActorIdIsZero), "%v", head)
	}
	dstRouter := tab.Get(head.Uid)
	if cls.Get(pb.NodeType_NodeTypeGate, dstRouter.Get(pb.NodeType_NodeTypeGate)) == nil {
		return uerror.N(1, int32(pb.ErrorCode_NodeNotFound), "%v", head)
	}
	if head.Cmd%2 == 0 {
		if _, ok := pb.CMD_name[int32(head.Cmd)]; ok {
			head.Cmd++
			head.Seq++
		}
	}
	if head.Dst == nil {
		head.Dst = &pb.NodeRouter{}
	}
	srcRouter := tab.Get(head.Src.ActorId)
	head.Src.NodeType = node.Type
	head.Src.NodeId = node.Id
	head.Src.Router = srcRouter.GetData()

	srcRouter.Set(head.Dst.NodeType, head.Dst.NodeId)
	dstRouter.Set(head.Src.NodeType, head.Src.NodeId)

	head.Dst.NodeType = pb.NodeType_NodeTypeGate
	head.Dst.NodeId = dstRouter.Get(pb.NodeType_NodeTypeGate)
	head.Dst.ActorName = "Player"
	head.Dst.FuncName = "SendToClient"
	head.Dst.ActorId = head.Uid
	head.Dst.Router = dstRouter.GetData()

	buf, err := marshal(msg)
	if err != nil {
		return uerror.N(1, int32(pb.ErrorCode_MarshalFailed), "%v|%v", head, err)
	}
	return bus.Send(head, buf)
}

func NotifyToClient(head *pb.Head, msg proto.Message, uids ...uint64) error {
	if head.Src == nil {
		return uerror.N(1, int32(pb.ErrorCode_NodeRouterIsNil), "%v", head)
	}
	if head.Src.ActorId <= 0 {
		return uerror.N(1, int32(pb.ErrorCode_ActorIdIsZero), "%v", head)
	}
	buf, err := marshal(msg)
	if err != nil {
		return uerror.N(1, int32(pb.ErrorCode_MarshalFailed), "%v|%v", head, err)
	}
	if head.Cmd%2 == 0 {
		if _, ok := pb.CMD_name[int32(head.Cmd)]; ok {
			head.Cmd++
			head.Seq++
		}
	}
	srcRouter := tab.Get(head.Src.ActorId)
	head.Src.NodeType = node.Type
	head.Src.NodeId = node.Id
	head.Src.Router = srcRouter.GetData()

	for _, uid := range uids {
		dstRouter := tab.Get(uid)
		dstNodeId := dstRouter.Get(pb.NodeType_NodeTypeGate)
		if cls.Get(pb.NodeType_NodeTypeGate, dstNodeId) == nil {
			mlog.Errorf("[%s] 节点不存在: %v", pb.ErrorCode_NodeNotFound, head)
			continue
		}
		if head.Dst == nil {
			head.Dst = &pb.NodeRouter{}
		}
		head.Uid = uid
		head.Dst.NodeType = pb.NodeType_NodeTypeGate
		head.Dst.NodeId = dstNodeId
		head.Dst.ActorName = "Player"
		head.Dst.FuncName = "SendToClient"
		head.Dst.ActorId = uid
		head.Dst.Router = dstRouter.GetData()
		srcRouter.Set(head.Dst.NodeType, head.Dst.NodeId)
		dstRouter.Set(head.Src.NodeType, head.Src.NodeId)
		if err := bus.Send(head, buf); err != nil {
			mlog.Errorf("通知客户端失败：head:%v, error:%v", head, err)
		}
	}
	return nil
}
