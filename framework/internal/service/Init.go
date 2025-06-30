package service

import (
	"universal/common/pb"
	"universal/common/yaml"
	"universal/framework/domain"
	buss "universal/framework/internal/bus"
	"universal/framework/internal/cluster"
	"universal/framework/internal/discovery"
	"universal/framework/internal/router"
	"universal/library/mlog"
	"universal/library/uerror"
)

var (
	tab  domain.ITable = router.NewTable()
	cls  domain.INode  = cluster.NewCluster()
	node *pb.Node
	dis  domain.IDiscovery
	bus  domain.IBus
)

func Init(nn *pb.Node, cfg *yaml.Config) error {
	node = nn
	etcdcli, err := discovery.NewEtcd(cfg.Etcd)
	if err != nil {
		mlog.Errorf("Etcd客户端连接失败: %v", err)
		return err
	}
	dis = etcdcli

	if err := dis.Watch(cls); err != nil {
		mlog.Errorf("Etcd服务监听失败：%v", err)
		return err
	}

	if err := dis.Register(node, cfg.Common.DiscoveryExpire); err != nil {
		mlog.Errorf("Etcd注册服务失败：%v", err)
		return err
	}

	natsCli, err := buss.NewNats(cfg.Nats)
	if err != nil {
		mlog.Errorf("Nats客户端连接失败")
		return err
	}
	bus = natsCli
	return nil
}

func SetBroadcastHandler(f func(*pb.Head, []byte)) error {
	return bus.SetBroadcastHandler(node, f)
}

func SetSendHandler(f func(*pb.Head, []byte)) error {
	return bus.SetSendHandler(node, f)
}

func SetReplyHandler(f func(*pb.Head, []byte)) error {
	return bus.SetReplyHandler(node, f)
}

func Close() error {
	tab.Close()
	dis.Close()
	bus.Close()
	return nil
}

func Dispatcher(head *pb.Head) error {
	srcRouter := tab.Get(head.Src.ActorId)
	dstRouter := tab.Get(head.Dst.ActorId)
	if head.Dst.NodeId > 0 {
		if cls.Get(head.Dst.NodeType, head.Dst.NodeId) == nil {
			return uerror.N(1, int32(pb.ErrorCode_NodeNotFound), "%v", head)
		}
	} else {
		head.Dst.NodeId = dstRouter.Get(head.Dst.NodeType)
		if cls.Get(head.Dst.NodeType, head.Dst.NodeId) == nil {
			nn := cls.Random(head.Dst.NodeType, head.Dst.ActorId)
			head.Dst.NodeId = nn.Id
		}
	}
	head.Src.NodeType = node.Type
	head.Src.NodeId = node.Id
	srcRouter.Set(head.Src.NodeType, head.Src.NodeId)
	srcRouter.Set(head.Dst.NodeType, head.Dst.NodeId)
	dstRouter.Set(head.Src.NodeType, head.Src.NodeId)
	dstRouter.Set(head.Dst.NodeType, head.Dst.NodeId)
	head.Src.Router = srcRouter.GetData()
	head.Dst.Router = dstRouter.GetData()
	return nil
}
