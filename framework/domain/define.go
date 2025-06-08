package domain

import (
	"universal/common/pb"

	"github.com/golang/protobuf/proto"
)

// 路由接口
type IRouter interface {
	Get(pb.NodeType) int32        // 获取路由信息
	Set(pb.NodeType, int32)       // 设置路由信息
	GetData() *pb.Router          // 获取路由信息
	SetData(*pb.Router)           // 设置路由信息
	IsExpire(now, ttl int64) bool // 是否过期
}

// 路由表接口
type ITable interface {
	Get(pb.NodeType, string, uint64) IRouter // 获取路由信息
	SetExpire(int64)                         // 设置路由过期时间
	Close()                                  // 关闭路由管理
}

// 服务集群接口
type ICluster interface {
	GetCount(pb.NodeType) int            // 获取节点数量
	Get(pb.NodeType, int32) *pb.Node     // 获取节点
	Del(pb.NodeType, int32)              // 删除节点
	Add(*pb.Node)                        // 添加节点
	List(pb.NodeType) []*pb.Node         // 获取节点列表
	Random(pb.NodeType, uint64) *pb.Node // 随机节点
}

// 服务注册与发现接口
type IDiscovery interface {
	Register(node *pb.Node, ttl int64) error // 注册服务
	Watch(ICluster) error                    // 监听服务
	Close() error                            // 关闭服务发现
}

// 消息总线接口
type IBus interface {
	SetBroadcastHandler(*pb.Node, func(*pb.Head, []byte)) error // 监听消息
	SetSendHandler(*pb.Node, func(*pb.Head, []byte)) error      // 监听消息
	SetReplyHandler(*pb.Node, func(*pb.Head, []byte)) error     // 监听消息
	Broadcast(*pb.Head, []byte) error                           // 广播消息
	Send(*pb.Head, []byte) error                                // 发送消息
	Request(*pb.Head, []byte, proto.Message) error              // 请求应答消息
	Response(*pb.Head, []byte) error                            // 响应消息
	Close()                                                     // 关闭网络
}
