package domain

import (
	"universal/common/pb"
	"universal/framework/basic"

	"google.golang.org/protobuf/proto"
)

// 对外接口定义
type ApiFunc func(*basic.Context, proto.Message, proto.Message) error

type IApi interface {
	Call(*basic.Context, proto.Message, proto.Message) error
}

type IPacket interface {
	Call(*basic.Context, []byte) (*pb.Packet, error)
}
