package service

import (
	"universal/common/pb"
	"universal/library/uerror"

	"github.com/golang/protobuf/proto"
)

func Request(head *pb.Head, msg interface{}, rsp proto.Message) error {
	if head.Src == nil {
		return uerror.N(1, int32(pb.ErrorCode_NodeRouterIsNil), "%v", head)
	}
	if head.Src.ActorId <= 0 {
		return uerror.N(1, int32(pb.ErrorCode_ActorIdIsZero), "%v", head)
	}
	if head.Dst == nil {
		return uerror.N(1, int32(pb.ErrorCode_NodeRouterIsNil), "%v", head)
	}
	if head.Dst.NodeType <= pb.NodeType_NodeTypeBegin || head.Dst.NodeType >= pb.NodeType_NodeTypeEnd {
		return uerror.N(1, int32(pb.ErrorCode_NodeTypeNotSupport), "%v", head)
	}
	if cls.GetCount(head.Dst.NodeType) <= 0 {
		return uerror.N(1, int32(pb.ErrorCode_NodeNotFound), "%v", head)
	}
	if err := Dispatcher(head); err != nil {
		return err
	}
	if head.Dst.NodeType == node.Type && head.Dst.NodeId == node.Id {
		return uerror.N(1, int32(pb.ErrorCode_SendToSelfWrong), "%v", head)
	}
	buf, err := marshal(msg)
	if err != nil {
		return uerror.N(1, int32(pb.ErrorCode_MarshalFailed), "%v|%v", head, err)
	}
	return bus.Request(head, buf, rsp)
}

func Response(head *pb.Head, msg interface{}) error {
	if len(head.Reply) <= 0 {
		return nil
	}
	head.SendType = pb.SendType_POINT
	buf, err := marshal(msg)
	if err != nil {
		return uerror.N(1, int32(pb.ErrorCode_MarshalFailed), "%v|%v", head, err)
	}
	return bus.Response(head, buf)
}
