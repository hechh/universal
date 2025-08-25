package define

import (
	"time"
	"universal/common/pb"

	"github.com/golang/protobuf/proto"
)

type IAsync interface {
	GetIdPointer() *uint64
	GetId() uint64
	SetId(uint64)
	Start()
	Stop()
}

type IActor interface {
	IAsync
	GetActorName() string
	Register(IActor, ...int)
	ParseFunc(interface{})
	SendMsg(*pb.Head, ...interface{}) error
	Send(*pb.Head, []byte) error
	RegisterTimer(*uint64, *pb.Head, time.Duration, int32) error
}

// 业务层Rsp必须实现的接口
type IRspProto interface {
	proto.Message
	GetHead() *pb.RspHead
	SetHead(*pb.RspHead)
}

// 应答函数
type SendRspFunc func(*pb.Head, IRspProto) error

type IHandler interface {
	Call(SendRspFunc, interface{}, *pb.Head, ...interface{}) func()
	Rpc(SendRspFunc, interface{}, *pb.Head, []byte) func()
}
