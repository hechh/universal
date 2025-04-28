package framework

import (
	"strings"
	"universal/framework/config"
	"universal/framework/define"
	"universal/framework/internal/cluster"
	"universal/framework/internal/discovery"
)

type Framework struct {
	router define.IRouter    // 路由表
	cls    define.ICluster   // 服务集群
	dis    define.IDiscovery // 服务发现
	net    define.INetwork   // 消息中间件
}

// Init 初始化框架
func (f *Framework) Init(appid int32, srv *config.ServerConfig, cfg *config.Config) (err error) {
	// 节点配置
	node := &cluster.Node{Name: srv.NodeName, Addr: srv.Nodes[appid], Type: srv.NodeType, Id: appid}

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
