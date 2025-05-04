package framework

import (
	"fmt"
	"universal/framework/config"
	"universal/framework/define"
	"universal/framework/internal/actor"
	"universal/framework/internal/cluster"
	"universal/framework/internal/discovery"
	"universal/framework/internal/network"
	"universal/framework/internal/packet"
	"universal/framework/internal/router"
	"universal/library/baselib/uerror"
)

type Actor struct{ actor.Actor }
type ActorGroup struct{ actor.ActorGroup }

type Framework struct {
	router define.IRouter           // 路由表
	cls    define.ICluster          // 服务集群
	dis    define.IDiscovery        // 服务发现
	net    define.INetwork          // 消息中间件
	actors map[string]define.IActor // Actor列表
}

func (f *Framework) RegisterActor(st define.IActor) {
	name := st.GetName()
	if _, ok := f.actors[name]; ok {
		panic(fmt.Sprintf("Actor已存在: %s", name))
	}
	f.actors[name] = st
}

func (f *Framework) Close() error {
	if err := f.router.Close(); err != nil {
		return err
	}
	if err := f.dis.Close(); err != nil {
		return err
	}
	if err := f.net.Close(); err != nil {
		return err
	}
	return nil
}

// Init 初始化框架
func (f *Framework) Init(cfg *config.Config, nodeType define.NodeType, appid int32) (err error) {
	// 节点配置
	nodeName := define.NodeType_name[int32(nodeType)]
	nodeCfg, ok := cfg.Cluster[nodeName]
	if !ok || nodeCfg == nil {
		return uerror.New(1, -1, "服务节点配置不存在：%s", define.NodeType_name[int32(nodeType)])
	}

	// 路由表初始化
	f.router = router.NewRouter()
	f.router.Expire(nodeCfg.RouterTTL)

	// 集群初始化
	f.cls = cluster.NewCluster(&cluster.Node{
		Name: nodeName,
		Addr: nodeCfg.Nodes[appid],
		Type: int32(nodeType),
		Id:   appid,
	})

	// 服务发现与注册
	f.dis, err = discovery.Init(cfg, discovery.WithPath("universal/discovery"), discovery.WithParse(cluster.NewNode))
	if err != nil {
		return err
	}
	// 注册服务
	if err := f.dis.KeepAlive(f.cls.GetSelf(), 15); err != nil {
		return err
	}
	// 监听服务
	if err := f.dis.Watch(f.cls); err != nil {
		return err
	}

	// 消息中间件
	f.net, err = network.Init(cfg, network.WithTopic("universal/network"), network.WithParse(packet.ParsePacket), network.WithNew(packet.NewPacket))
	if err != nil {
		return err
	}
	return
}
