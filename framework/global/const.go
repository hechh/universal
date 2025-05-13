package global

// 节点类型
type NodeType int32

const (
	NodeTypeBegin NodeType = 0 // 未知
	NodeTypeGate  NodeType = 1 // 网关
	NodeTypeGame  NodeType = 2 // 游戏服务
	NodeTypeGm    NodeType = 3 // 游戏服务
	NodeTypeDb    NodeType = 4 // 数据服务
	NodeTypeRoom  NodeType = 5 // 房间服务
	NodeTypeMatch NodeType = 6 // 匹配服务
	NodeTypeMax   NodeType = 7 // 最大节点类型
)

var (
	NodeType_name = map[int32]string{
		int32(NodeTypeGate):  "gate",
		int32(NodeTypeGame):  "game",
		int32(NodeTypeGm):    "gm",
		int32(NodeTypeDb):    "db",
		int32(NodeTypeRoom):  "room",
		int32(NodeTypeMatch): "match",
	}
)
