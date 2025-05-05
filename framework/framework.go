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
	"universal/library/encode"
	"universal/library/mlog"

	"google.golang.org/protobuf/proto"
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
func (f *Framework) Init(cfg *config.Config, nodeType define.NodeType, appid uint32) (err error) {
	// 节点配置
	nodeName := define.NodeType_name[uint32(nodeType)]
	nodeCfg, ok := cfg.Cluster[nodeName]
	if !ok || nodeCfg == nil {
		return uerror.New(1, -1, "服务节点配置不存在：%s", define.NodeType_name[uint32(nodeType)])
	}

	// 路由表初始化
	f.router = router.NewRouter()
	f.router.Expire(nodeCfg.RouterTTL)

	// 集群初始化
	f.cls = cluster.NewCluster(&cluster.Node{
		Name: nodeName,
		Addr: nodeCfg.Nodes[appid],
		Type: uint32(nodeType),
		Id:   appid,
	})

	// 服务发现与注册
	if err := f.initDiscovery(cfg); err != nil {
		return err
	}
	// 消息中间件
	if err = f.initNetwork(cfg); err != nil {
		return err
	}
	return
}

func (f *Framework) initDiscovery(cfg *config.Config) error {
	dis, err := discovery.Init(cfg, discovery.WithTopic("universal/discovery"), discovery.WithNode(cluster.NewNode))
	if err != nil {
		return err
	}
	f.dis = dis
	if err := dis.KeepAlive(f.cls.GetSelf(), 15); err != nil {
		return err
	}
	return dis.Watch(f.cls)
}

func (f *Framework) initNetwork(cfg *config.Config) error {
	net, err := network.Init(cfg, network.WithTopic("universal/network"), network.WithPacket(packet.NewPacket), network.WithHeader(packet.NewHeader))
	if err != nil {
		return err
	}
	f.net = net
	return net.Read(f.cls.GetSelf(), func(head define.IHeader, body []byte) {
		act, ok := f.actors[head.GetActorName()]
		if !ok {
			mlog.Error("Actor不存在: %v", head)
			return
		}
		if err := act.SendFrom(packet.NewContext(head, f.router), body); err != nil {
			mlog.Error("Actor发送失败: %v", err)
		}
	})
}

// 服务内调用
func (f *Framework) Send(head define.IHeader, args ...interface{}) error {
	act, ok := f.actors[head.GetActorName()]
	if !ok {
		return uerror.New(1, -1, "Actor不存在: %v", head)
	}
	return act.Send(packet.NewContext(head, f.router), args...)
}

// 夸服务调用
func (f *Framework) SendTo(head define.IHeader, args ...interface{}) error {
	if len(args) == 1 {
		if msg, ok := args[0].(proto.Message); ok {
			buf, _ := proto.Marshal(msg)
			return f.net.Send(head, buf)
		}
	}
	return f.net.Send(head, encode.Encode(args...))
}

// 夸服务路由
func (f *Framework) dispatcher(uid, routeId uint64, nodeType uint32) define.IHeader {
	//node := f.cls.GetSelf()

	return &packet.Header{}
}
