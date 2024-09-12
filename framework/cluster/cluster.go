package cluster

import (
	"net"
	"runtime/debug"
	"universal/common/global"
	"universal/common/pb"
	"universal/framework/basic"
	"universal/framework/cluster/domain"
	"universal/framework/cluster/internal/discovery"
	"universal/framework/cluster/internal/handler"
	"universal/framework/cluster/internal/nodes"
	"universal/framework/cluster/internal/router"
	"universal/framework/plog"
	"universal/framework/uerror"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/proto"
)

var (
	etcd    *discovery.EtcdClient // 服务发现
	natsCli *nats.Conn            // nats服务
)

func Close() {
	etcd.Close()
	natsCli.Close()
}

func Init(cfg *global.Config, typ pb.SERVER, serverId uint32, expire int64) (err error) {
	// 建立etcd连接
	if etcd, err = discovery.NewEtcdClient(cfg.Etcd.Endpoints...); err != nil {
		return
	}
	// 建立nats连接
	if natsCli, err = nats.Connect(cfg.Nats.Endpoints); err != nil {
		return
	}
	// 订阅
	fatal := func(err interface{}) {
		plog.Fatal("%v\nstack: %s", err, string(debug.Stack()))
	}
	natsCli.Subscribe(nodes.GetSelfChannel(), func(msg *nats.Msg) {
		inner := &pb.Packet{}
		proto.Unmarshal(msg.Data, inner)
		basic.SafeRecover(fatal, func() { handler.HandlePoint(inner.Head, inner.Body) })
	})
	natsCli.Subscribe(nodes.GetSelfTopicChannel(), func(msg *nats.Msg) {
		inner := &pb.Packet{}
		proto.Unmarshal(msg.Data, inner)
		basic.SafeRecover(fatal, func() { handler.HandleTopic(inner.Head, inner.Body) })
	})
	// 初始化
	var host, port string
	if host, port, err = net.SplitHostPort(cfg.Server[serverId].Host); err != nil {
		return err
	}
	router.Init(expire)
	nodes.Init(domain.ROOT_DIR, &pb.ServerInfo{
		Type:       typ,
		Ip:         host,
		Port:       cast.ToInt32(port),
		CreateTime: uint64(basic.GetNowUnixMilli()),
		ServerID:   basic.GetCrc(cfg.Server[serverId].Host),
	})
	// 服务发现
	if err = etcd.Watch(domain.ROOT_DIR, nodes.AddNotify, nodes.DeleteNotify); err != nil {
		return err
	}
	buf, err := proto.Marshal(nodes.GetSelf())
	if err != nil {
		return err
	}
	// 设置租赁保活
	return etcd.KeepAlive(nodes.GetSelfChannel(), string(buf), 15)
}

func Dispatcher(head *pb.Head) (err error) {
	head.SendType = pb.SEND_Point
	head.SrcServerID = nodes.GetSelf().ServerID
	head.SrcServerType = nodes.GetSelf().Type
	if head.UID <= 0 {
		return uerror.NewUError(1, -1, "玩家ID为空")
	}
	// 加载路由表
	table := router.GetOrNew(head.UID)
	// 目的节点是否已经确定
	if head.DstServerID > 0 {
		dst := nodes.Get(head.DstServerID)
		if dst == nil || dst.Type != head.DstServerType {
			return uerror.NewUError(1, -1, "服务节点不存在: %s(%d)", head.DstServerType.String(), head.DstServerID)
		}
		// 更新路由
		table.Update(head.DstServerType, head.DstServerID)
		head.Table = table.Get()
		return
	}

	// 从路由表中加载
	clusterId := table.GetServerID(head.DstServerType)
	if dst := nodes.Get(clusterId); dst != nil {
		head.DstServerID = clusterId
		head.Table = table.Get()
		return
	}

	// 重新路由
	switch head.RouteType {
	case 0: // 路由类型-玩家id
		if node := nodes.Random(head.DstServerType, head.UID); node == nil {
			return uerror.NewUError(1, -1, "服务节点不存在: %s", head.DstServerType.String())
		} else {
			head.DstServerID = node.ServerID
			// 更新路由
			table.Update(head.DstServerType, node.ServerID)
			head.Table = table.Get()
		}
	case 1: // 路由类型-区服id
		if node := nodes.Random(head.DstServerType, uint64(head.RegionID)); node == nil {
			return uerror.NewUError(1, -1, "服务节点不存在: %s", head.DstServerType.String())
		} else {
			head.DstServerID = node.ServerID
			// 更新路由
			table.Update(head.DstServerType, node.ServerID)
			head.Table = table.Get()
		}
	}
	return
}

// 发送内部广播
func Broadcast(head *pb.Head, data proto.Message) error {
	head.SendType = pb.SEND_Broadcast
	head.SrcServerID = nodes.GetSelf().GetServerID()

	// 判断服务节点是否存在
	if rets := nodes.Gets(head.DstServerType); len(rets) <= 0 {
		return uerror.NewUError(1, -1, "服务节点不存在: %s", head.DstServerType.String())
	}

	// 内部服务转发
	buf, _ := proto.Marshal(data)
	inner := &pb.Packet{Head: head, Body: buf}
	ret, _ := proto.Marshal(inner)
	return natsCli.Publish(nodes.GetTopicChannel(head.DstServerType), ret)
}

// 发送广播查询路由，路由存在就发送回报，不存在路由就丢弃
func Query(head *pb.Head, data proto.Message) error {
	head.SendType = pb.SEND_Broadcast
	head.SrcServerID = nodes.GetSelf().GetServerID()

	if head.UID <= 0 {
		return uerror.NewUError(1, -1, "玩家ID为空")
	}

	// 携带路由
	if rr := router.Get(head.UID); rr != nil {
		head.Table = rr.Get()
	}

	// 内部服务转发
	buf, _ := proto.Marshal(data)
	inner := &pb.Packet{Head: head, Body: buf}
	ret, _ := proto.Marshal(inner)
	return natsCli.Publish(nodes.GetTopicChannel(head.DstServerType), ret)
}

// 发送内部单播
func Send(head *pb.Head, data proto.Message) error {
	// 路由
	if err := Dispatcher(head); err != nil {
		return err
	}

	// 内部服务转发
	buf, _ := proto.Marshal(data)
	inner := &pb.Packet{Head: head, Body: buf}
	ret, _ := proto.Marshal(inner)
	return natsCli.Publish(nodes.GetChannel(head.DstServerType, head.DstServerID), ret)
}
