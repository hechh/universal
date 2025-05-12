package domain

// 节点类型
type NodeType int32

const (
	NodeTypeBegin NodeType = 0 // 未知
	NodeTypeGate  NodeType = 1 // 网关
	NodeTypeGame  NodeType = 2 // 游戏服务
	NodeTypeDb    NodeType = 3 // 数据服务
	NodeTypeRoom  NodeType = 4 // 房间服务
	NodeTypeMatch NodeType = 5 // 匹配服务
	NodeTypeMax   NodeType = 6 // 最大节点类型
)

var (
	NodeType_name = map[int32]string{
		int32(NodeTypeGate):  "gate",
		int32(NodeTypeGame):  "game",
		int32(NodeTypeDb):    "db",
		int32(NodeTypeRoom):  "room",
		int32(NodeTypeMatch): "match",
	}
)

type IBytes interface {
	GetSize() int          // 获取协议头大小
	WriteTo([]byte) error  // 写入数据
	ReadFrom([]byte) error // 读取数据
}

// 内网协议包接口
type IPacket interface {
	IBytes
	GetHead() IHead           // 获取协议头
	GetRoute() IRouter        // 获取路由信息
	GetBody() []byte          // 获取协议体
	SetHead(IHead) IPacket    // 设置协议头
	SetRoute(IRouter) IPacket // 设置路由信息
	SetBody([]byte) IPacket   // 设置协议体
}

// 协议头接口
type IHead interface {
	IBytes
	GetSrcNodeType() int32      // 消息来源类型
	GetSrcNodeId() int32        // 消息来源ID
	GetDstNodeType() int32      // 消息目的类型
	GetDstNodeId() int32        // 消息目的ID
	GetRouteId() uint64         // 消息路由ID
	GetUid() uint64             // 消息用户ID
	GetActorName() string       // 消息Actor名称
	GetFuncName() string        // 消息函数名称
	SetSrcNodeType(int32) IHead // 设置消息来源类型
	SetSrcNodeId(int32) IHead   // 设置消息来源ID
	SetDstNodeType(int32) IHead // 设置消息目的类型
	SetDstNodeId(int32) IHead   // 设置消息目的ID
	SetRouteId(uint64) IHead    // 设置消息路由ID
	SetUid(uint64) IHead        // 设置消息用户ID
	SetActorName(string) IHead  // 设置消息Actor名称
	SetFuncName(string) IHead   // 设置消息函数名称
}

// 路由信息接口
type IRouter interface {
	IBytes
	Get(nodeType int32) int32   // 获取路由信息
	Set(nodeType, nodeId int32) // 设置路由信息
}

// 路由管理接口
type IRouterMgr interface {
	Get(uint64) IRouter  // 获取路由信息
	Set(uint64, IRouter) // 设置路由信息
	Close()              // 关闭路由管理
}

// Actor接口定义
type IActor interface {
	GetActorName() string             // 获取Actor名称
	Register(self IActor)             // 注册Actor
	ParseFunc(interface{})            // 解析方法列表
	Send(IHead, ...interface{}) error // 发送消息
	SendRpc(IHead, []byte) error      // 发送远程调用
}

// Actor管理接口
type IActorMgr interface {
	Register(IActor)                  // 注册Actor
	Get(string) IActor                // 获取Actor
	Send(IHead, ...interface{}) error // 发送消息
	SendRpc(IHead, []byte) error      // 发送远程调用
}

// 服务节点接口
type INode interface {
	IBytes
	GetName() string      // 获取节点名称
	GetType() int32       // 获取节点类型
	GetId() int32         // 获取节点ID
	GetAddr() string      // 获取节点地址
	SetName(string) INode // 设置节点名称
	SetType(int32) INode  // 设置节点类型
	SetId(int32) INode    // 设置节点ID
	SetAddr(string) INode // 设置节点地址
	String() string       // 获取节点字符串
}

// 服务集群接口
type ICluster interface {
	Get(nodeType, nodeId int32) INode           // 获取节点
	Del(NodeType, nodeId int32)                 // 删除节点
	Add(node INode)                             // 添加节点
	Random(nodeType int32, hashId uint64) INode // 随机节点
}

// 服务发现接口
type IDiscovery interface {
	Register(self INode, ttl int64) error // 注册服务
	Watch(cls ICluster) error             // 监听服务
	Close() error                         // 关闭服务发现
}

// 网络接口
type INetwork interface {
	Receive(node INode, mgr IActorMgr) error // 监听消息
	Send(head IHead, data []byte) error      // 发送消息
	Broadcast(head IHead, data []byte) error // 广播消息
	Close() error                            // 关闭网络
}
