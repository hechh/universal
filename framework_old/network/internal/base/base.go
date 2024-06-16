package base

import (
	"universal/common/pb"
	"universal/framework/common/fbasic"
	"universal/framework/common/uerror"

	"google.golang.org/protobuf/proto"
)

func ReqToPacket(head *pb.PacketHead, req proto.Message, params ...interface{}) (*pb.Packet, error) {
	// 设置参数
	if len(params) > 0 {
		switch vv := req.(type) {
		case *pb.ActorRequest:
			vv.Buff = fbasic.EncodeAny(params...)
		}
	}
	// 封装
	buf, err := proto.Marshal(req)
	if err != nil {
		return nil, uerror.NewUError(1, -1, err)
	}
	return &pb.Packet{Head: head, Buff: buf}, nil
}

func RspToPacket(head *pb.PacketHead, rsp proto.Message, params ...interface{}) (*pb.Packet, error) {
	// 设置参数
	if len(params) > 0 {
		switch vv := rsp.(type) {
		case *pb.ActorResponse:
			vv.Buff = fbasic.EncodeAny(params...)
		}
	}
	// 序列化
	buf, err := proto.Marshal(rsp)
	if err != nil {
		return nil, uerror.NewUError(1, -1, err)
	}
	return &pb.Packet{Head: head, Buff: buf}, nil
}

func NewActorRequest(actorName, funcName string, params ...interface{}) proto.Message {
	return &pb.ActorRequest{
		ActorName: actorName,
		FuncName:  funcName,
		Buff:      fbasic.EncodeAny(params...),
	}
}
