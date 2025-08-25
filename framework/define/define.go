package define

import (
	"universal/common/pb"

	"github.com/golang/protobuf/proto"
)

// 业务层Rsp必须实现的接口
type IRspProto interface {
	proto.Message
	GetHead() *pb.RspHead
	SetHead(*pb.RspHead)
}

// 应答函数
type SendRspFunc func(*pb.Head, IRspProto) error

type IHandler interface {
	Call(SendRspFunc, *pb.Head, ...interface{}) func()
	Rpc(SendRspFunc, *pb.Head, []byte) func()
}
