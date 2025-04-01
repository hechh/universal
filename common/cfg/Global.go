package cfg

type QualityType uint32

const (
	QualityType_None  QualityType = 0 // 品质类型-空
	QualityType_White QualityType = 1 // 品质类型-白
	QualityType_Blue  QualityType = 2 // 品质类型-蓝
)

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

type PropertyType uint32

const (
	PropertyType_Empty   PropertyType = 0 // 道具类型-空
	PropertyType_Coin    PropertyType = 1 // 道具类型-金币
	PropertyType_Diamond PropertyType = 2 // 道具类型-钻石
	PropertyType_Sword   PropertyType = 3 // 装备类型-剑
)
