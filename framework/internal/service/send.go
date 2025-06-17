package service

import (
	"universal/common/pb"
	"universal/library/uerror"
)

func Send(head *pb.Head, args ...interface{}) error {
	if head.Src == nil {
		return uerror.N(1, int32(pb.ErrorCode_ParamInvalid), "Src为空: %v", head)
	}
	srcRouter := tab.Get(head.Src.ActorId)
	srcRouter.Set(node.Type, node.Id)
	head.Src.NodeType = node.Type
	head.Src.NodeId = node.Id
	head.Src.Router = srcRouter.GetData()

	if head.Dst == nil ||
		head.Dst.NodeType <= pb.NodeType_NodeTypeBegin ||
		head.Dst.NodeType >= pb.NodeType_NodeTypeEnd {
		return uerror.N(1, int32(pb.ErrorCode_ParamInvalid), "Dst为空: %v", head)
	}

	if cls.GetCount(head.Dst.NodeType) <= 0 {
		return uerror.N(1, int32(pb.ErrorCode_NodeNotFound), "服务节点不存在：%v", head)
	}

	return nil
}
