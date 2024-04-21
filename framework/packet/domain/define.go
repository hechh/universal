package domain

import (
	"universal/common/pb"
	"universal/framework/basic"

	"google.golang.org/protobuf/proto"
)

type IPacket interface {
	Call(*basic.Context, *pb.Packet) *pb.Packet
}

type CmdFunc func(*basic.Context, proto.Message, proto.Message) error
