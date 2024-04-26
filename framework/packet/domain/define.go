package domain

import (
	"universal/common/pb"
	"universal/framework/fbasic"

	"google.golang.org/protobuf/proto"
)

// 对外接口定义
type ApiFunc func(*fbasic.Context, proto.Message, proto.Message) error

type IApi interface {
	Call(*fbasic.Context, proto.Message, proto.Message) error
}

type IPacket interface {
	Call(*fbasic.Context, []byte) (*pb.Packet, error)
}
