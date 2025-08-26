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

// 处理器接口
type IHandler interface {
	Call(SendRspFunc, IActor, *pb.Head, ...interface{}) func()
	Rpc(SendRspFunc, IActor, *pb.Head, []byte) func()
}

// 数据帧
type IFrame interface {
	GetSize(*pb.Packet) int
	Decode([]byte, *pb.Packet) error
	Encode(*pb.Packet, []byte) error
}

// 网络接口
type INet interface {
	SetFrame(IFrame)
	SetReadExpire(int64)
	SetWriteExpire(int64)
	Write(*pb.Packet) error
	Read(*pb.Packet) error
	Close() error
}

// 节点管理接口
type INode interface {
	GetSelf() *pb.Node
	GetCount(pb.NodeType) int
	Get(pb.NodeType, int32) *pb.Node
	Del(pb.NodeType, int32) bool
	Add(*pb.Node) bool
	Random(pb.NodeType, uint64) *pb.Node
}

// 服务注册与服务发现接口
type IDiscovery interface {
	Register(INode, int64) error
	Watch(INode) error
	Close() error
}

// 路由接口
type IRouter interface {
	GetData() *pb.Router
	SetData(*pb.Router) IRouter
	Get(pb.NodeType) int32
	Set(pb.NodeType, int32) IRouter
}

type ITable interface {
	GetOrNew(pb.RouterType, uint64, *pb.Node) IRouter
	Get(pb.RouterType, uint64) IRouter
	Close()
}

// 消息接口
type IBus interface {
	SetBroadcastHandler(*pb.Node, func(*pb.Head, []byte)) error
	SetSendHandler(*pb.Node, func(*pb.Head, []byte)) error
	SetReplyHandler(*pb.Node, func(*pb.Head, []byte)) error
	Broadcast(*pb.Head, []byte) error
	Send(*pb.Head, []byte) error
	Request(*pb.Head, []byte, proto.Message) error
	Response(*pb.Head, []byte) error
	Close()
}
