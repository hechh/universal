package cfg

type RouteConfig struct {
	ID         uint32     // 唯一ID
	Request    string     // 请求名称
	Response   string     // 应答协议
	FuncName   string     // 函数名称
	ServerType ServerType // 服务类型
	RouteType  RouteType  // 路由类型
	Rewards    []*Reward  // 奖励
}
