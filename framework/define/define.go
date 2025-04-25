package define

import (
	"context"
)

type IHeader interface {
	GetSendType() uint32
	GetSrcType() uint32
	GetSrcId() uint32
	GetDstType() uint32
	GetDstId() uint32
	GetCmd() uint32
	GetUid() uint64
}

// 服务节点
type INode interface {
	Unique() string     // 唯一标识
	GetName() string    // 服务名称
	GetType() int32     // 服务类型
	GetId() int32       // 服务id
	GetAddress() string // 服务地址
}

// 服务集群
type ICluster interface {
	Get(nodeType int32, nodeId int32) INode    // 获取节点
	Insert(INode) error                        // 添加节点
	Delete(nodeType int32, nodeId int32) error // 删除节点
	Random(nodeType int32, seed uint64) INode  // 随机节点
}

// 服务发现
type IDiscovery interface {
	Put(ctx context.Context, srv INode) error                  // 注册服务
	Delete(ctx context.Context, srv INode) error               // 删除服务
	Watch(ctx context.Context, cluster ICluster) error         // 服务发现
	KeepAlive(ctx context.Context, srv INode, ttl int64) error // 心跳
}

// 路由表
type IRouter interface {
	Get(uint64) ICluster
	Update(uint64, ICluster) error
}

// 内网消息协议
type IPacket interface {
	GetHeader() IHeader
	GetBody() []byte
}

// 转发消息
type ISend interface {
	Send(ctx context.Context, pack IPacket) error      // 发送消息
	Broadcast(ctx context.Context, pack IPacket) error // 广播消息
}

type Header struct {
	Sequence       uint32
	Cmd            uint32
	Uid            uint64
	SendType       uint32
	SrcClusterType uint32
	SrcClusterId   uint32
	DstClusterType uint32
	DstClusterId   uint32
}

type Packet struct {
	Header *Header
	Body   []byte
}
