package service

import (
	"net"
	"runtime/debug"
	"universal/common/config"
	"universal/common/pb"
	"universal/framework/basic/uerror"
	"universal/framework/basic/util"
	"universal/framework/cluster/internal/discovery"
	"universal/framework/cluster/internal/handler"
	"universal/framework/cluster/internal/nodes"
	"universal/framework/cluster/internal/router"
	"universal/framework/plog"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/proto"
)

const (
	ROOT_DIR = "server"
)

var (
	etcd    *discovery.EtcdClient // 服务发现
	natsCli *nats.Conn            // nats服务
)

func Close() {
	etcd.Close()
	natsCli.Close()
}

func Init(cfg *config.Config, typ pb.SERVICE, serverId uint32, expire int64) (err error) {
	// 建立etcd连接
	if etcd, err = discovery.NewEtcdClient(cfg.Etcd.Endpoints...); err != nil {
		return
	}
	// 建立nats连接
	if natsCli, err = nats.Connect(cfg.Nats.Endpoints); err != nil {
		return
	}
	// 订阅
	natsCli.Subscribe(nodes.GetSelfChannel(), func(msg *nats.Msg) {
		inner := &pb.RpcPacket{}
		proto.Unmarshal(msg.Data, inner)
		util.SafeRecover(func(err interface{}) {
			plog.Fatal("%v\nstack: %s", err, string(debug.Stack()))
		}, func() {
			handler.HandlePoint(inner.RpcHead, inner.RpcBody)
		})
	})
	natsCli.Subscribe(nodes.GetSelfTopicChannel(), func(msg *nats.Msg) {
		inner := &pb.RpcPacket{}
		proto.Unmarshal(msg.Data, inner)
		util.SafeRecover(func(err interface{}) {
			plog.Fatal("%v\nstack: %s", err, string(debug.Stack()))
		}, func() {
			handler.HandleTopic(inner.RpcHead, inner.RpcBody)
		})
	})
	// 初始化
	var host, port string
	if host, port, err = net.SplitHostPort(cfg.Server[serverId].Host); err != nil {
		return err
	}
	// 服务发现
	if err = etcd.Watch(ROOT_DIR, nodes.AddNotify, nodes.DeleteNotify); err != nil {
		return err
	}
	router.Init(expire)
	nodes.Init(&pb.ClusterInfo{
		Type:       typ,
		Ip:         host,
		Port:       cast.ToInt32(port),
		CreateTime: uint64(util.GetNowUnixMilli()),
		ClusterID:  util.GetCrc(cfg.Server[serverId].Host),
	}, ROOT_DIR, cfg.Stub)
	buf, err := proto.Marshal(nodes.GetSelf())
	if err != nil {
		return err
	}
	// 设置租赁保活
	return etcd.KeepAlive(nodes.GetSelfChannel(), string(buf), 30)
}

func Dispatcher(head *pb.RpcHead) (err error) {
	head.SendType = pb.SEND_POINT
	head.SrcClusterId = nodes.GetSelf().ClusterID
	if head.Id <= 0 {
		return uerror.NewUError(1, -1, "玩家ID为空")
	}
	// 加载路由表
	table := router.GetOrNew(head.Id)
	head.Route = table.Get()

	// 目的节点是否已经确定
	if head.ClusterId > 0 {
		dst := nodes.Get(head.ClusterId)
		if dst == nil || dst.Type != head.DestServerType {
			return uerror.NewUError(1, -1, "服务节点不存在: %s(%d)", head.DestServerType.String(), head.ClusterId)
		}
		// 更新路由
		table.Update(head.DestServerType, head.ClusterId)
		return
	}

	// 从路由表中加载
	clusterId := table.GetClusterID(head.DestServerType)
	if dst := nodes.Get(clusterId); dst != nil {
		head.ClusterId = clusterId
		return
	}

	// 重新路由
	switch head.RouteType {
	case 0: // 路由类型-玩家id
		if node := nodes.Random(head.DestServerType, head.Id, head.ActorName); node == nil {
			return uerror.NewUError(1, -1, "服务节点不存在: %s", head.DestServerType.String())
		} else {
			head.ClusterId = node.ClusterID
			// 更新路由
			table.Update(head.DestServerType, node.ClusterID)
		}
	case 1: // 路由类型-区服id
		if node := nodes.Random(head.DestServerType, uint64(head.RegionID), head.ActorName); node == nil {
			return uerror.NewUError(1, -1, "服务节点不存在: %s", head.DestServerType.String())
		} else {
			head.ClusterId = node.ClusterID
			// 更新路由
			table.Update(head.DestServerType, node.ClusterID)
		}
	}
	return
}

// 发送内部广播
func Broadcast(head *pb.RpcHead, data proto.Message) error {
	head.SendType = pb.SEND_BROAD_CAST
	head.SrcClusterId = nodes.GetSelf().GetClusterID()

	// 判断服务节点是否存在
	if rets := nodes.Gets(head.DestServerType); len(rets) <= 0 {
		return uerror.NewUError(1, -1, "服务节点不存在: %s", head.DestServerType.String())
	}

	// 内部服务转发
	buf, _ := proto.Marshal(data)
	inner := &pb.RpcPacket{RpcHead: head, RpcBody: buf}
	ret, _ := proto.Marshal(inner)
	return natsCli.Publish(nodes.GetTopicChannel(head.DestServerType), ret)
}

// 发送广播查询路由，路由存在就发送回报，不存在路由就丢弃
func Query(head *pb.RpcHead, data proto.Message) error {
	head.SendType = pb.SEND_BROAD_CAST
	head.SrcClusterId = nodes.GetSelf().GetClusterID()

	if head.Id <= 0 {
		return uerror.NewUError(1, -1, "玩家ID为空")
	}

	// 携带路由
	if rr := router.Get(head.Id); rr != nil {
		head.Route = rr.Get()
	}

	// 内部服务转发
	buf, _ := proto.Marshal(data)
	inner := &pb.RpcPacket{RpcHead: head, RpcBody: buf}
	ret, _ := proto.Marshal(inner)
	return natsCli.Publish(nodes.GetTopicChannel(head.DestServerType), ret)
}

// 发送内部单播
func Send(head *pb.RpcHead, data proto.Message) error {
	// 路由
	if err := Dispatcher(head); err != nil {
		return err
	}

	// 内部服务转发
	buf, _ := proto.Marshal(data)
	inner := &pb.RpcPacket{RpcHead: head, RpcBody: buf}
	ret, _ := proto.Marshal(inner)
	return natsCli.Publish(nodes.GetChannel(head.DestServerType, head.ClusterId), ret)
}
