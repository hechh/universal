package cluster

import (
	"universal/common/pb"
	"universal/common/yaml"
	"universal/framework/internal/bus"
	"universal/framework/internal/discovery"
	"universal/framework/internal/method"
	"universal/framework/internal/node"
	"universal/framework/internal/router"
	"universal/library/pprof"

	"github.com/golang/protobuf/proto"
)

var (
	clusterObj *Cluster
)

func Init(cfg *yaml.Config, srvCfg *yaml.NodeConfig, nn *pb.Node) error {
	cls := node.NewNode(nn)
	tab := router.NewTable(srvCfg.RouterTTL)

	dis, err := discovery.NewEtcd(cfg.Etcd)
	if err != nil {
		return err
	}
	if err := dis.Watch(cls); err != nil {
		return err
	}
	if err := dis.Register(cls, srvCfg.DiscoveryTTL); err != nil {
		return err
	}

	buss, err := bus.NewNats(cfg.Nats)
	if err != nil {
		return err
	}
	clusterObj = NewCluster(cls, tab, dis, buss)

	method.Init(clusterObj.SendResponse)
	pprof.Init("localhost", srvCfg.Port+10000)
	return nil
}

func Close() {
	clusterObj.Close()
}

func GetSelf() *pb.Node {
	return clusterObj.GetSelf()
}

func SetBroadcastHandler(f func(*pb.Head, []byte)) error {
	return clusterObj.SetBroadcastHandler(f)
}

func SetSendHandler(f func(*pb.Head, []byte)) error {
	return clusterObj.SetSendHandler(f)
}

func SetReplyHandler(f func(*pb.Head, []byte)) error {
	return clusterObj.SetReplyHandler(f)
}

func Broadcast(head *pb.Head, args ...interface{}) error {
	return clusterObj.Broadcast(head, args...)
}

func Send(head *pb.Head, args ...interface{}) error {
	return clusterObj.Send(head, args...)
}

func SendCmd(head *pb.Head, args ...interface{}) error {
	return clusterObj.SendCmd(head, args...)
}

func SendToClient(head *pb.Head, msg proto.Message, uids ...uint64) error {
	return clusterObj.SendToClient(head, msg, uids...)
}

func Request(head *pb.Head, msg interface{}, rsp proto.Message) error {
	return clusterObj.Request(head, msg, rsp)
}
