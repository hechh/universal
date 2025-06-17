package service

import (
	"sync/atomic"
	"universal/common/pb"
	"universal/library/encode"
	"universal/library/uerror"

	"github.com/golang/protobuf/proto"
)

func Send(head *pb.Head, args ...interface{}) error {
	if head.Src == nil {
		return uerror.N(1, int32(pb.ErrorCode_ParamInvalid), "Src为空: %v", head)
	}
	if head.Src.ActorId <= 0 {
		return uerror.N(1, int32(pb.ErrorCode_ParamInvalid), "Src.ActorId为空: %v", head)
	}
	if head.Dst == nil {
		return uerror.N(1, int32(pb.ErrorCode_ParamInvalid), "Dst为空: %v", head)
	}
	if head.Dst.NodeType <= pb.NodeType_NodeTypeBegin || head.Dst.NodeType >= pb.NodeType_NodeTypeEnd {
		return uerror.N(1, int32(pb.ErrorCode_ParamInvalid), "Dst服务类型不支持: %v", head)
	}
	if cls.GetCount(head.Dst.NodeType) <= 0 {
		return uerror.N(1, int32(pb.ErrorCode_NodeNotFound), "服务节点不存在：%v", head)
	}
	if err := dispatcher(head); err != nil {
		return err
	}
	if head.Dst.NodeType == node.Type && head.Dst.NodeId == node.Id {
		return uerror.N(1, int32(pb.ErrorCode_ParamInvalid), "不能调用当前服务节点: %v", head)
	}
	atomic.AddUint32(&head.Reference, 1)
	buf, err := marshal(args...)
	if err != nil {
		return err
	}
	return bus.Send(head, buf)
}

func marshal(args ...interface{}) ([]byte, error) {
	if len(args) == 1 {
		switch vv := args[0].(type) {
		case []byte:
			return vv, nil
		case proto.Message:
			if buf, err := proto.Marshal(vv); err != nil {
				return nil, err
			} else {
				return buf, nil
			}
		}
	}
	return encode.Encode(args...)
}

func dispatcher(head *pb.Head) error {
	srcRouter := tab.Get(head.Src.ActorId)
	srcRouter.Set(node.Type, node.Id)

	dstRouter := tab.Get(head.Dst.ActorId)
	dstRouter.Set(node.Type, node.Id)

	if head.Dst.NodeId > 0 {
		if cls.Get(head.Dst.NodeType, head.Dst.NodeId) == nil {
			return uerror.N(1, int32(pb.ErrorCode_NodeNotFound), "未找到服务节点: %v", head)
		}
	} else {
		head.Dst.NodeId = dstRouter.Get(head.Dst.NodeType)
		if cls.Get(head.Dst.NodeType, head.Dst.NodeId) == nil {
			nn := cls.Random(head.Dst.NodeType, head.Dst.ActorId)
			head.Dst.NodeId = nn.Id
		}
	}

	srcRouter.Set(head.Dst.NodeType, head.Dst.NodeId)
	dstRouter.Set(head.Dst.NodeType, head.Dst.NodeId)

	head.Src.Router = srcRouter.GetData()
	head.Dst.Router = dstRouter.GetData()
	return nil
}
