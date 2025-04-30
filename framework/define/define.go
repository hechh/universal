package define

type ParseNodeFunc func([]byte) INode

// 服务节点
type INode interface {
	GetName() string // 服务名称
	GetType() int32  // 服务类型
	GetId() int32    // 服务id
	GetAddr() string // 服务地址
	ToBytes() []byte // 转换为字节数组
}

// 路由表
type IRouter interface {
	Get(id uint64, nodeType int32) INode // 获取路由信息
	Update(id uint64, node INode)        // 更新路由信息
	Expire(ttl int64)                    // 设置存活时间
}

// 服务集群
type ICluster interface {
	GetSelf() INode                           // 获取当前节点
	Get(nodeType int32, nodeId int32) INode   // 获取节点
	Put(node INode) error                     // 添加节点
	Del(nodeType int32, nodeId int32) error   // 删除节点
	Random(nodeType int32, seed uint64) INode // 随机节点
}

// 服务发现
type IDiscovery interface {
	Get() ([]INode, error)                // 获取服务列表
	Put(srv INode) error                  // 注册服务
	Del(srv INode) error                  // 删除服务
	Watch(cluster ICluster) error         // 服务发现
	KeepAlive(srv INode, ttl int64) error // 心跳
	Close() error                         // 删除
}

// 转发消息
type INetwork interface {
	Read(node INode, listen func(IHeader, []byte)) error   // 接收消息
	Send(node INode, head IHeader, body []byte) error      // 发送消息
	Broadcast(node INode, head IHeader, body []byte) error // 广播消息
}

type ParsePacketFunc func([]byte) (IPacket, error)
type NewPacketFunc func(IHeader, []byte) IPacket

// 内网消息协议
type IPacket interface {
	GetHeader() IHeader
	GetBody() []byte
	ToBytes() []byte
}

type IHeader interface {
	GetSendType() int32               // 获取发送类型
	GetSrcType() int32                // 获取源服务类型
	GetSrcId() int32                  // 获取源服务id
	GetDstType() int32                // 获取目的服务类型
	GetDstId() int32                  // 获取目的服务id
	GetCmd() uint32                   // 获取命令
	GetUid() uint64                   // 获取用户id
	GetRouteId(nodeType int32) uint64 // 获取路由id
}

// cmd 请求接口
type IProto interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

type IContext interface {
}
type HandleFunc func(IContext, IProto, IProto) error
