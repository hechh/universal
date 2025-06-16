package domain

import (
	"time"
	"universal/common/pb"

	"github.com/golang/protobuf/proto"
)

type IActor interface {
	GetId() uint64                                      // 获取actor id
	SetId(uint64)                                       // 设置 actor id
	Start()                                             // 开始协程
	Stop()                                              // 停止协程
	GetActorName() string                               // 获取 actor名称
	Register(IActor)                                    // 注册 iactor
	ParseFunc(interface{})                              // 解析成员函数
	RegisterTimer(*pb.Head, time.Duration, int32) error // 注册定时器
	SendMsg(*pb.Head, ...interface{}) error             // 发送消息
	Send(*pb.Head, []byte) error                        // 发送消息
}

type ICluster interface {
	GetCount(pb.NodeType) int            // 获取节点数量
	Get(pb.NodeType, int32) *pb.Node     // 获取节点
	Del(pb.NodeType, int32) bool         // 删除节点
	Add(*pb.Node) bool                   // 添加节点
	Random(pb.NodeType, uint64) *pb.Node // 随机节点
}

type IDiscovery interface {
	Register(*pb.Node, int64) error // 注册服务
	Watch(ICluster) error           // 监听服务
	Close() error                   // 关闭服务
}

type IRouter interface {
	Get(pb.NodeType) int32
	Set(pb.NodeType, int32)
	GetData() *pb.Node
	SetData(*pb.Node)
}

type ITable interface {
	Get(uint64) IRouter
}

type IBus interface {
	SetBroadcastHandler(*pb.Node, func(*pb.Packet)) error
	SetSendHandler(*pb.Node, func(*pb.Packet)) error
	SetReplyHandler(*pb.Node, func(*pb.Packet)) error
	Broadcast(*pb.Head, []byte) error
	Send(*pb.Head, []byte) error
	Request(*pb.Head, []byte, proto.Message) error
	Response(*pb.Head, []byte) error
	Close() error
}

type IRspProto interface {
	proto.Message
	GetHead() *pb.RspHead
	SetHead(*pb.RspHead)
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
