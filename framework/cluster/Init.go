package cluster

import (
	"universal/common/pb"
	"universal/common/yaml"
	"universal/framework/domain"
	"universal/framework/internal/bus"
	"universal/framework/internal/discovery"
	"universal/framework/internal/node"
	"universal/framework/internal/router"
	"universal/library/uerror"
)

var (
	self *pb.Node
	tab  domain.ITable
	cls  domain.INode
	dis  domain.IDiscovery
	buss domain.IBus
)

func Init(cfg *yaml.Config, nodeType pb.NodeType, nodeId int32) error {
	srvCfg := yaml.GetNodeConfig(cfg, nodeType, nodeId)
	if srvCfg == nil {
		return uerror.N(1, int32(pb.ErrorCode_NodeConfigNotFound), "%s(%d)", nodeType, nodeId)
	}

	self = yaml.GetNode(cfg, nodeType, nodeId)
	tab = router.NewTable(srvCfg.RotuerTTL)
	cls = node.NewNode()

	// 消息中间件
	if cli, err := bus.NewNats(cfg.Nats); err != nil {
		return uerror.E(1, int32(pb.ErrorCode_NatsConnectFailed), err)
	} else {
		buss = cli
	}

	// 服务注册与发现
	if cli, err := discovery.NewEtcd(cfg.Etcd); err != nil {
		return uerror.E(1, int32(pb.ErrorCode_EtcdConnectFailed), err)
	} else {
		dis = cli
	}
	if err := dis.Watch(cls); err != nil {
		return uerror.E(1, int32(pb.ErrorCode_EtcdWatchFailed), err)
	}
	if err := dis.Register(self, srvCfg.DiscoveryTTL); err != nil {
		return uerror.E(1, int32(pb.ErrorCode_EtcdRegisterFailed), err)
	}
	return nil
}

func Close() {
	tab.Close()
	dis.Close()
	buss.Close()
}

func SetBroadcastHandler(f func(*pb.Head, []byte)) error {
	return buss.SetBroadcastHandler(self, f)
}

func SetSendHandler(f func(*pb.Head, []byte)) error {
	return buss.SetSendHandler(self, func(head *pb.Head, body []byte) {
		if head.Src != nil && head.Src.Router != nil {
			tabSrc := tab.GetOrNew(head.Src.ActorId)
			tabSrc.SetData(head.Src.Router).Set(self.Type, self.Id)
		}

		tab.GetOrNew(head.Dst.ActorId).SetData(head.Dst.Router).Set(self.Type, self.Id)

		f(head, body)
	})
}

func SetReplyHandler(f func(*pb.Head, []byte)) error {
	return buss.SetReplyHandler(self, func(head *pb.Head, body []byte) {
		if head.Src != nil && head.Src.Router != nil {
			tabSrc := tab.GetOrNew(head.Src.ActorId)
			tabSrc.SetData(head.Src.Router).Set(self.Type, self.Id)
		}

		tab.GetOrNew(head.Dst.ActorId).SetData(head.Dst.Router).Set(self.Type, self.Id)

		f(head, body)
	})
}
