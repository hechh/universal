package define

import (
	"context"
)

// 服务节点
type IServer interface {
	GetServerType() int32 // 服务类型
	GetServerId() int32   // 服务id
	GetAddress() string   // 服务地址
}

// 服务发现
type IDiscovery interface {
	Put(ctx context.Context, srv IServer) error                            // 注册服务
	Delete(ctx context.Context, srv IServer) error                         // 删除服务
	Watch(ctx context.Context, add func(IServer), del func(IServer)) error // 服务发现
	KeepAlive(ctx context.Context, srv IServer, ttl int64) error           // 心跳
}

// 路由表
type IRouter interface {
	Get(uint64) IServer
	Update(uint64, IServer) error
}

// 服务节点
type INode interface {
	GetSelf() IServer                              // 获取当前节点
	Get(serverType int32, serverId int32) IServer  // 获取节点
	Insert(IServer) error                          // 添加节点
	Delete(serverType int32, serverId int32) error // 删除节点
	Random(serverType int32, seed uint64) IServer  // 随机节点
}
