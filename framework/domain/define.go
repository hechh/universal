package domain

import "universal/common/pb"

// -----------------actor核心接口-------------------
// Actor接口定义
type IActor interface {
	GetActorName() string                // 获取Actor名称
	Register(IActor, ...int)             // 注册Actor
	ParseFunc(interface{})               // 解析方法列表
	Send(*pb.Head, ...interface{}) error // 发送消息
	SendRpc(*pb.Head, []byte) error      // 发送远程调用
}

// Actor管理接口
type IActorMgr interface {
	Register(IActor)                     // 注册Actor
	Get(string) IActor                   // 获取Actor
	Send(*pb.Head, ...interface{}) error // 发送消息
	SendRpc(*pb.Head, []byte) error      // 发送远程调用
}

// 路由管理接口
type IRouterMgr interface {
	Get(uint64) *pb.Router  // 获取路由信息
	Set(uint64, *pb.Router) // 设置路由信息
	Close()                 // 关闭路由管理
}

// -----------------------服务发现接口-------------------
// 服务集群接口
type ICluster interface {
	Get(pb.NodeType, int32) *pb.Node     // 获取节点
	Del(pb.NodeType, int32)              // 删除节点
	Add(*pb.Node)                        // 添加节点
	List(pb.NodeType) []*pb.Node         // 获取节点列表
	Random(pb.NodeType, uint64) *pb.Node // 随机节点
}

// --------------------服务注册与发现接口---------------------
type IDiscovery interface {
	Register(node *pb.Node, ttl int64) error // 注册服务
	Watch(ICluster) error                    // 监听服务
	Close() error                            // 关闭服务发现
}

// ------------------------消息总线接口-----------------------
type IBus interface {
	Listen(*pb.Node, IActorMgr) error // 监听消息
	Send(*pb.Head, []byte) error      // 发送消息
	Broadcast(*pb.Head, []byte) error // 广播消息
	Close() error                     // 关闭网络
	//	Request(head *pb.Head, data []byte) error   // 请求消息
}
