package define

import (
	"universal/common/pb"

	"github.com/golang/protobuf/proto"
)

// 支持处理的类型
type ZeroFunc func() error
type OneFunc func(proto.Message) error
type TwoFunc func(proto.Message, proto.Message) error
type HeadZeroFunc func(*pb.Head) error
type HeadOneFunc func(*pb.Head, proto.Message) error
type HeadTwoFunc func(*pb.Head, proto.Message, proto.Message) error

// 业务层Rsp必须实现的接口
type IRspProto interface {
	proto.Message
	GetHead() *pb.RspHead
	SetHead(*pb.RspHead)
}

// 应答函数
type SendRspFunc func(*pb.Head, IRspProto) error

// pb对象工厂类
type IFactory interface {
	New(string) proto.Message
}
