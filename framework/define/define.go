package define

type RouteInfo struct {
	Gate uint32
	Db   uint32
	Game uint32
	Tool uint32
	Rank uint32
}

// 路由表
type IRouter interface {
	Get(uint64) *RouteInfo     // 获取路由信息
	Update(uint64, *RouteInfo) // 更新路由信息
	Expire(int64)              // 设置存活时间
	Close() error              // 关闭路由
}

// 服务节点
type INode interface {
	GetName() string      // 服务名称
	GetType() uint32      // 服务类型
	GetId() uint32        // 服务id
	GetAddr() string      // 服务地址
	ToBytes() []byte      // 转换为字节数组
	Parse([]byte) INode   // 解析 bytes
	SetName(string) INode // 设置名字
	SetAddr(string) INode // 设置地址
	SetType(uint32) INode // 设置类型
	SetId(uint32) INode   // 设置 ID
}

// 服务集群
type ICluster interface {
	GetSelf() INode                            // 获取当前节点
	Put(node INode)                            // 添加节点
	Del(nodeType uint32, nodeId uint32)        // 删除节点
	Get(nodeType uint32, nodeId uint32) INode  // 获取节点
	Random(nodeType uint32, seed uint64) INode // 随机节点
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
	Read(node INode, listen func(IHeader, []byte)) error // 接收消息
	Send(head IHeader, body []byte) error                // 发送消息
	Broadcast(head IHeader, body []byte) error           // 发送广播
	Close() error                                        // 关闭
}

// 内网消息协议
type IPacket interface {
	GetHeader() IHeader        // 获取包头
	GetBody() []byte           // 获取包体
	ToBytes() []byte           // 转成 bytes
	Parse([]byte) IPacket      // 解析 bytes
	SetHeader(IHeader) IPacket // 设置包头
	SetBody([]byte) IPacket    // 设置包体
}

type IHeader interface {
	GetSrcNodeType() uint32      // 获取源服务类型
	GetSrcNodeId() uint32        // 获取源服务id
	GetDstNodeType() uint32      // 获取目的服务类型
	GetDstNodeId() uint32        // 获取目的服务id
	GetCmd() uint32              // 获取命令
	GetUid() uint64              // 获取用户id
	GetRouteId() uint64          // 获取路由id
	GetActorName() string        // 获取服务名称
	GetFuncName() string         // 获取函数名称
	GetTable() *RouteInfo        // 获取路由表
	GetSize() int                // 获取头部大小
	ToBytes([]byte) []byte       // 转换为字节数组
	Parse([]byte) IHeader        // 解析头部
	SetCmd(uint32) IHeader       // 设置命令
	SetUid(uint64) IHeader       // 设置用户id
	SetRouteId(uint64) IHeader   // 设置路由id
	SetSrcNode(INode) IHeader    // 设置源服务
	SetDstNode(INode) IHeader    // 设置目的服务
	SetActorName(string) IHeader // 设置服务名称
	SetFuncName(string) IHeader  // 设置函数名称
	SetTable(*RouteInfo) IHeader // 设置路由表

}

type IContext interface {
	IHeader
	GetRouter() IRouter
}

type IActor interface {
	GetName() string                              // 获取名称
	Register(IActor, interface{}) error           // 注册方法
	Send(ctx IContext, args ...interface{}) error // 发送请求
	SendFrom(ctx IContext, buf []byte) error      // 发送请求
}
