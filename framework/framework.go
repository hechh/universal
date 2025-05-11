package framework

import (
	"universal/framework/config"
	"universal/framework/domain"
	"universal/framework/internal/actor"
	"universal/framework/internal/cluster"
	"universal/framework/internal/discovery"
	"universal/framework/internal/network"
	"universal/framework/internal/packet"
	"universal/framework/internal/route"
)

type Actor struct{ actor.Actor }
type ActorMgr struct{ actor.ActorMgr }

type Framework struct {
	self     domain.INode
	rmgr     domain.IRouteMgr
	cls      domain.ICluster
	dis      domain.IDiscovery
	net      domain.INetwork
	actors   map[string]domain.IActor
	newNode  func() domain.INode
	newHead  func() domain.IHead
	newRoute func() domain.IRoute
	newPack  func() domain.IPacket
}

func (f *Framework) Init(node domain.INode, cfg *config.Config) (err error) {
	f.newNode = cluster.NewNode
	f.newHead = packet.NewHeader
	f.newRoute = route.NewRoute
	f.newPack = packet.NewPacket
	f.actors = make(map[string]domain.IActor)
	f.self = node
	f.cls = cluster.NewCluster()

	// 初始化路由管理
	clsCfg := cfg.Cluster[node.GetName()]
	f.rmgr = route.NewRouterMgr(f.newRoute, clsCfg.RouteTTL)

	// 服务注册与发现
	if f.dis, err = discovery.Init(cfg,
		discovery.WithTopic("universl/discovery"),
		discovery.WithNode(f.newNode),
	); err != nil {
		return err
	}
	if err := f.dis.Register(f.self, 15); err != nil {
		return err
	}
	if err := f.dis.Watch(f.cls); err != nil {
		return err
	}

	// 初始化网络
	if f.net, err = network.Init(cfg,
		network.WithTopic("universl/network"),
		network.WithPacket(f.newPack),
		network.WithHead(f.newHead),
		network.WithRoute(f.newRoute),
		network.WithRouteMgr(f.rmgr),
	); err != nil {
		return err
	}

	return nil
}
