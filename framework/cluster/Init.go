package cluster

import (
	"universal/common/pb"
	"universal/common/yaml"
	"universal/framework/domain"
	"universal/framework/internal/bus"
	"universal/framework/internal/discovery"
	"universal/framework/internal/node"
	"universal/framework/internal/router"
	"universal/library/mlog"
	"universal/library/safe"
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
	safe.Catch(mlog.Fatalf)
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
		UpdateRouter(head)
		f(head, body)
	})
}

func SetReplyHandler(f func(*pb.Head, []byte)) error {
	return buss.SetReplyHandler(self, func(head *pb.Head, body []byte) {
		UpdateRouter(head)
		f(head, body)
	})
}

func UpdateRouter(head *pb.Head) {
	if head.Src != nil && head.Src.Router != nil {
		tabSrc := tab.GetOrNew(head.Src.ActorId)
		tabSrc.SetData(head.Src.Router).Set(self.Type, self.Id)
	}
	tab.GetOrNew(head.Dst.ActorId).SetData(head.Dst.Router).Set(self.Type, self.Id)
}

func Dispatcher(head *pb.Head) error {
	if head.Dst == nil || head.Dst.ActorId <= 0 {
		return uerror.N(1, int32(pb.ErrorCode_DstNodeRouterIsNil), "%v", head)
	}
	if head.Dst.NodeType >= pb.NodeType_NodeTypeEnd || head.Dst.NodeType <= pb.NodeType_NodeTypeBegin {
		return uerror.N(1, int32(pb.ErrorCode_NodeTypeNotSupported), "%v", head.Dst)
	}
	// 业务层直接指定具体节点
	dstTab := tab.GetOrNew(head.Dst.ActorId).Set(self.Type, self.Id)
	if head.Dst.NodeId > 0 {
		if cls.Get(head.Dst.NodeType, head.Dst.NodeId) == nil {
			return uerror.N(1, int32(pb.ErrorCode_NodeNotFound), "%v", head.Dst)
		}
		head.Dst.Router = dstTab.GetData()
		return nil
	}
	// 优先从路由中选择
	if nodeId := dstTab.Get(head.Dst.NodeType); nodeId > 0 {
		if nn := cls.Get(head.Dst.NodeType, nodeId); nn != nil {
			dstTab.Set(nn.Type, nn.Id)
			head.Dst.Router = dstTab.GetData()
			return nil
		}
	}
	//从集群中随机获取一个节点
	if nn := cls.Random(head.Dst.NodeType, head.Dst.ActorId); nn != nil {
		dstTab.Set(nn.Type, nn.Id)
		head.Dst.NodeId = nn.Id
		return nil
	}
	return uerror.N(1, int32(pb.ErrorCode_NodeNotFound), "%v", head.Dst)
}
