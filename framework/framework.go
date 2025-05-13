package framework

import (
	"universal/common/pb"
	"universal/framework/config"
	"universal/framework/domain"
	"universal/framework/internal/actor"
	"universal/framework/internal/cluster"
	"universal/framework/internal/discovery"
	"universal/framework/internal/network"
	"universal/framework/internal/packet"
	"universal/framework/internal/router"
)

type Actor struct{ actor.Actor }

type ActorGroup struct{ actor.ActorGroup }

type Framework struct {
	self     *pb.Node
	routeMgr domain.IRouterMgr
	actMgr   domain.IActorMgr
	cls      domain.ICluster
	dis      domain.IDiscovery
	net      domain.INetwork
	newNode  func() *pb.Node
	newHead  func() *pb.Head
	newRoute func() domain.IRouter
	newPack  func() domain.IPacket
}

func (f *Framework) Init(node *pb.Node, cfg *config.Config) (err error) {
	f.newNode = cluster.NewNode
	f.newHead = packet.NewHeader
	f.newRoute = router.NewRouter
	f.newPack = packet.NewPacket
	f.self = node
	f.cls = cluster.NewCluster()
	f.actMgr = actor.NewActorMgr()

	// 初始化路由管理
	clsCfg := cfg.Cluster[node.GetName()]
	f.routeMgr = router.NewRouterMgr(f.newRoute, clsCfg.RouteTTL)

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
		network.WithRouteMgr(f.routeMgr),
	); err != nil {
		return err
	}
	if err := f.net.Receive(f.self, f.actMgr); err != nil {
		return err
	}
	return nil
}
