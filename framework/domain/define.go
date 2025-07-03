package domain

import (
	"universal/common/pb"

	"github.com/golang/protobuf/proto"
)

type IRspProto interface {
	proto.Message
	GetHead() *pb.RspHead
	SetHead(*pb.RspHead)
}

type IAsync interface {
	GetIdPointer() *uint64
	GetId() uint64
	SetId(uint64)
	Start()
	Stop()
	Push(f func())
}

type IActor interface {
	IAsync
	GetActorName() string
	Register(IActor, ...int)
	ParseFunc(interface{})
	SendMsg(*pb.Head, ...interface{}) error
	Send(*pb.Head, []byte) error
}

type INode interface {
	GetCount(pb.NodeType) int
	Get(pb.NodeType, int32) *pb.Node
	Del(pb.NodeType, int32) bool
	Add(*pb.Node) bool
	Random(pb.NodeType, uint64) *pb.Node
}

type IDiscovery interface {
	Register(*pb.Node, int64) error
	Watch(INode) error
	Close() error
}

type IRouter interface {
	GetData() *pb.Router
	SetData(*pb.Router) IRouter
	Get(pb.NodeType) int32
	Set(pb.NodeType, int32) IRouter
}

type ITable interface {
	GetOrNew(uint64, *pb.Node) IRouter
	Get(uint64) IRouter
	Close()
}

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

type IFrame interface {
	GetSize(*pb.Packet) int
	Decode([]byte, *pb.Packet) error
	Encode(*pb.Packet, []byte) error
}

type INet interface {
	Register(IFrame)
	SetReadExpire(int64)
	SetWriteExpire(int64)
	Write(*pb.Packet) error
	Read(*pb.Packet) error
	Close() error
}
