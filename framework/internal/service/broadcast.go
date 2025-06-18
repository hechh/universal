package service

import (
	"universal/common/pb"
	"universal/library/uerror"
)

func Broadcast(head *pb.Head, args ...interface{}) error {
	if head.Dst == nil {
		return uerror.N(1, int32(pb.ErrorCode_NodeRouterIsNil), "%v", head)
	}
	if head.Dst.NodeType <= pb.NodeType_NodeTypeBegin || head.Dst.NodeType >= pb.NodeType_NodeTypeEnd {
		return uerror.N(1, int32(pb.ErrorCode_NodeTypeNotSupport), "%v", head)
	}
	if cls.GetCount(head.Dst.NodeType) <= 0 {
		return uerror.N(1, int32(pb.ErrorCode_NodeNotFound), "%v", head)
	}
	if head.Src != nil && head.Src.ActorId > 0 {
		srcRouter := tab.Get(head.Src.ActorId)
		srcRouter.Set(head.Src.NodeType, head.Src.NodeId)
		head.Src.NodeType = node.Type
		head.Src.NodeId = node.Id
		head.Src.Router = srcRouter.GetData()
	}
	buf, err := marshal(args...)
	if err != nil {
		return uerror.E(1, int32(pb.ErrorCode_MarshalFailed), err)
	}
	return bus.Broadcast(head, buf)
}
