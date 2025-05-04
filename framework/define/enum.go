package define

// 路由类型
type RouteType int32

const (
	RouteTypeNone   RouteType = 0 // 未知
	RouteTypeRandom RouteType = 1 // 随机路由
	RouteTypeHash   RouteType = 2 // 哈希路由
)

// 节点类型
type NodeType int32

const (
	NodeTypeBegin NodeType = 0 // 未知
	NodeTypeGate  NodeType = 1 // 网关
	NodeTypeDb    NodeType = 2 // 数据服务
	NodeTypeLogin NodeType = 3 // 登录服务
	NodeTypeGame  NodeType = 4 // 游戏服务
	NodeTypeTool  NodeType = 5 // 工具服务
	NodeTypeRank  NodeType = 6 // 排行服务
	NodeTypeMax   NodeType = 7 // 最大节点类型
)

var (
	NodeType_name = map[uint32]string{
		uint32(NodeTypeGate):  "gate",
		uint32(NodeTypeDb):    "db",
		uint32(NodeTypeLogin): "login",
		uint32(NodeTypeGame):  "game",
		uint32(NodeTypeTool):  "tool",
		uint32(NodeTypeRank):  "rank",
	}
)
