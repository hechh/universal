package define

import (
	"context"
)

// 服务节点
type IServer interface {
	GetServerType() int // 服务类型
	GetServerId() int   // 服务id
	GetAddress() string // 服务地址
}

// 服务发现
type IDiscovery interface {
	KeepAlive(context.Context, IServer, int64) error           // 心跳
	Put(context.Context, IServer) error                        // 注册服务
	Delete(context.Context, IServer) error                     // 删除服务
	Watch(context.Context, func(IServer), func(IServer)) error // 服务发现
}

// 路由表
type IRouter interface {
	Get(uint64) IServer
	Update(uint64, IServer) error
}

// 服务节点
type INode interface {
	Insert(IServer) error  // 添加节点
	Delete(IServer) error  // 删除节点
	Random(uint32) IServer // 随机节点
}
