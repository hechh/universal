package cluster

import (
	"poker_server/common/pb"
	"poker_server/common/yaml"
	"poker_server/framework/internal/bus"
	"poker_server/framework/internal/discovery"
	"poker_server/framework/internal/method"
	"poker_server/framework/internal/node"
	"poker_server/framework/internal/router"

	"github.com/golang/protobuf/proto"
)

var (
	obj *Cluster
)

func Init(nn *pb.Node, srvCfg *yaml.NodeConfig, cfg *yaml.Config) error {
	// 服务注册与发现
	cls := node.New(nn)
	tab := router.New(srvCfg.RouterTTL)
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

	// 消息中间件
	buss, err := bus.NewNats(cfg.Nats)
	if err != nil {
		return err
	}
	obj = New(cls, tab, buss, dis)
	method.Init(obj.SendResponse)
	return nil
}

func Close() {
	obj.Close()
}

func GetSelf() *pb.Node {
	return obj.GetSelf()
}

func GetAppName() string {
	return obj.GetSelf().Name
}

func GetAppAddr() string {
	return obj.GetSelf().Addr
}

func SendResponse(head *pb.Head, rsp proto.Message) error {
	return obj.SendResponse(head, rsp)
}

func SetBroadcastHandler(f func(*pb.Head, []byte)) error {
	return obj.SetBroadcastHandler(f)
}

func SetSendHandler(f func(*pb.Head, []byte)) error {
	return obj.SetSendHandler(f)
}

func SetReplyHandler(f func(*pb.Head, []byte)) error {
	return obj.SetReplyHandler(f)
}

func Broadcast(head *pb.Head, args ...interface{}) error {
	return obj.Broadcast(head, args...)
}

func Send(head *pb.Head, args ...interface{}) error {
	return obj.Send(head, args...)
}

func SendToClient(head *pb.Head, msg proto.Message, uids ...uint64) error {
	return obj.SendToClient(head, msg, uids...)
}

func Request(head *pb.Head, msg interface{}, rsp proto.Message) error {
	return obj.Request(head, msg, rsp)
}
