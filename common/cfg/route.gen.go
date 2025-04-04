type RouteType uint32

const (
	RouteType_None     RouteType = 0 // 路由类型-空
	RouteType_Uid      RouteType = 1 // 路由类型-UID
	RouteType_Random   RouteType = 2 // 路由类型-随机
	RouteType_RegionId RouteType = 3 // 路由类型-地区
)

type ServerType uint32

const (
	ServerType_Gate ServerType = 0 // 服务类型-Gate
	ServerType_Game ServerType = 1 // 服务类型-Game
	ServerType_Db   ServerType = 2 // 服务类型-Db
	ServerType_Gm   ServerType = 3 // 服务类型-Gm
)

type Reward struct {
	PropType PropertyType // 道具类型
	Quality  QualityType  // 品质
	Star     uint32       // 星级
	Add      int64        // 道具数量
}
type RouteConfig struct {
	ID         uint32     // 唯一ID
	Request    string     // 请求名称
	Response   string     // 应答协议
	FuncName   string     // 函数名称
	ServerType ServerType // 服务类型
	RouteType  RouteType  // 路由类型
	Rewards    []*Reward  // 奖励
}
