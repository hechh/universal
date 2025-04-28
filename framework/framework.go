package framework

import (
	"strings"
	"universal/framework/config"
	"universal/framework/define"
	"universal/framework/internal/cluster"
	"universal/framework/internal/discovery"
	"universal/framework/internal/router"
	"universal/library/baselib/uerror"
)

var (
	routes = map[int32]int32{
		int32(define.NodeTypeGate):  int32(define.RouteTypeHash),
		int32(define.NodeTypeDb):    int32(define.RouteTypeRandom),
		int32(define.NodeTypeLogin): int32(define.RouteTypeHash),
		int32(define.NodeTypeGame):  int32(define.RouteTypeHash),
		int32(define.NodeTypeTool):  int32(define.RouteTypeHash),
		int32(define.NodeTypeRank):  int32(define.RouteTypeHash),
	}
)

type Framework struct {
	router define.IRouter    // 路由表
	cls    define.ICluster   // 服务集群
	dis    define.IDiscovery // 服务发现
	net    define.INetwork   // 消息中间件
}

// Init 初始化框架
func (f *Framework) Init(cfg *config.Config, nodeType define.NodeType, appid int32) (err error) {
	f.router = router.NewRouter()
	// 节点配置
	nodeName := define.NodeType_name[int32(nodeType)]
	nodeCfg, ok := cfg.Cluster[nodeName]
	if !ok || nodeCfg == nil {
		return uerror.New(1, -1, "服务节点配置不存在：%s", define.NodeType_name[int32(nodeType)])
	}
	f.cls = cluster.NewCluster(&cluster.Node{
		Name: nodeName,
		Addr: nodeCfg.Nodes[appid],
		Type: int32(nodeType),
		Id:   appid,
	}, routes)

	// 服务发现与注册
	switch strings.ToLower(cfg.Middle.Discovery) {
	case "etcd":
		f.dis, err = discovery.NewEtcd(cfg.Etcd.Endpoints, discovery.WithPath(cfg.Middle.Path), discovery.WithParse(cluster.NewNode))
		if err != nil {
			return
		}
	case "consul":
		f.dis, err = discovery.NewConsul(cfg.Consul.Endpoints, discovery.WithPath(cfg.Middle.Path), discovery.WithParse(cluster.NewNode))
		if err != nil {
			return
		}
	}

	// 消息中间件
	switch cfg.Middle.Network {
	case "nats":
	}
	return
}
