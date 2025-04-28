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
	NodeTypeNone  NodeType = 0 // 未知
	NodeTypeGate  NodeType = 1 // 网关
	NodeTypeDb    NodeType = 2 // 数据服务
	NodeTypeLogin NodeType = 3 // 登录服务
	NodeTypeGame  NodeType = 4 // 游戏服务
	NodeTypeTool  NodeType = 5 // 工具服务
	NodeTypeRank  NodeType = 6 // 排行服务
)

var (
	NodeType_name = map[int32]string{
		int32(NodeTypeNone):  "node",
		int32(NodeTypeGate):  "gate",
		int32(NodeTypeDb):    "db",
		int32(NodeTypeLogin): "login",
		int32(NodeTypeGame):  "game",
		int32(NodeTypeTool):  "tool",
		int32(NodeTypeRank):  "rank",
	}
)
