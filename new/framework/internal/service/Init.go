package service

import (
	"poker_server/common/pb"
	"poker_server/common/yaml"
	"poker_server/framework/domain"
	"poker_server/framework/internal/core/actor"
	"poker_server/framework/internal/core/bus"
	"poker_server/framework/internal/core/cluster"
	"poker_server/framework/internal/core/discovery"
	"poker_server/framework/internal/core/router"
	"poker_server/framework/library/mlog"
)

var (
	self          *pb.Node               // 本节点信息
	clusterObj    = cluster.New()        // 服务集群对象
	tableObj      = router.New()         // 路由管理对象
	disObj        domain.IDiscovery      // 服务发现对象
	busObj        domain.IBus            // 消息总线对象
	broadcastFunc func(*pb.Head, []byte) // 广播消息回调函数
	sendFunc      func(*pb.Head, []byte) // 单播消息回调函数
	replyFunc     func(*pb.Head, []byte) // rpc消息回调函数
)

func init() {
	broadcastFunc = defaultBroadcastHandler
	sendFunc = defaultSendHandler
	replyFunc = defaultReplyHandler
}

func GetSelf() *pb.Node {
	return self
}

func SetSelf(node *pb.Node) {
	self = node
}

func GetCluster() domain.ICluster {
	return clusterObj
}

func GetTable() domain.ITable {
	return tableObj
}

func RegisterBroadcastHandler(f func(*pb.Head, []byte)) {
	broadcastFunc = f
}

func RegisterSendHandler(f func(*pb.Head, []byte)) {
	sendFunc = f
}

func RegisterReplyHandler(f func(*pb.Head, []byte)) {
	replyFunc = f
}

// 初始化框架
func Init(cfg *yaml.Config) (err error) {
	nodeCfg := cfg.Cluster[self.Name]
	// 初始化Etcd服务
	if disObj, err = discovery.NewEtcd(cfg.Etcd); err != nil {
		return
	}
	if err := disObj.Watch(clusterObj); err != nil {
		return err
	}
	if err := disObj.Register(self, nodeCfg.DicoveryExpire); err != nil {
		return err
	}

	// 初始化消息总线
	if busObj, err = bus.NewNats(cfg.Nats, tableObj); err != nil {
		return err
	}
	if err = busObj.SetBroadcastHandler(self, broadcastFunc); err != nil {
		return err
	}
	if err = busObj.SetSendHandler(self, sendFunc); err != nil {
		return err
	}
	if err = busObj.SetReplyHandler(self, replyFunc); err != nil {
		return err
	}

	// 初始化actor管理
	actor.Init(busObj)
	tableObj.SetExpire(nodeCfg.RouterExpire)

	mlog.Infof("框架服务初始化成功: %v", self)
	return nil
}

// 默认内网消息处理器
func defaultSendHandler(head *pb.Head, buf []byte) {
	head.SendType = pb.SendType_POINT
	if err := actor.Send(head, buf); err != nil {
		mlog.Errorf("跨服务调用错误: %v", err)
	}
}

func defaultReplyHandler(head *pb.Head, buf []byte) {
	head.SendType = pb.SendType_RPC
	if err := actor.Send(head, buf); err != nil {
		mlog.Errorf("跨服务调用错误: %v", err)
	}
}

// 默认内网广播消息处理器
func defaultBroadcastHandler(head *pb.Head, buf []byte) {
	head.SendType = pb.SendType_BROADCAST
	if err := actor.Send(head, buf); err != nil {
		mlog.Errorf("跨服务调用错误: %v", err)
	}
}
