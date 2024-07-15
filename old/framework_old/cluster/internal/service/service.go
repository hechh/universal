package service

import (
	"net"
	"universal/common/pb"
	"universal/framework/cluster/domain"
	"universal/framework/cluster/internal/discovery/etcd"
	"universal/framework/cluster/internal/nodes"
	"universal/framework/cluster/internal/router"
	"universal/framework/common/fbasic"
	"universal/framework/common/uerror"

	"github.com/spf13/cast"
	"google.golang.org/protobuf/proto"
)

var (
	dis      domain.IDiscovery // 服务发现
	selfNode *pb.ServerNode
)

func Close() {
	dis.Close()
}

func GetSelfServerNode() *pb.ServerNode {
	return selfNode
}

func GetDiscovery() domain.IDiscovery {
	return dis
}

// 初始化
func Init(ends []string) (err error) {
	// 初始化etcd
	etcd, err := etcd.NewEtcdClient(ends...)
	if err != nil {
		return err
	}
	// 设置服务发现
	dis = etcd
	return
}

// 服务发现
func Discovery(typ pb.ServerType, addr string) error {
	ip, port, err := net.SplitHostPort(addr)
	if err != nil {
		return uerror.NewUError(1, -1, err)
	}
	// 构建自身节点
	selfNode = &pb.ServerNode{
		ServerType: typ,
		ServerID:   fbasic.GetCrc32(addr),
		Ip:         ip,
		Port:       cast.ToInt32(port),
	}
	buf, err := proto.Marshal(selfNode)
	if err != nil {
		return uerror.NewUError(1, -1, err)
	}
	// 注册自身服务（保活，服务下线会自动删除）
	dis.KeepAlive(GetNodeChannel(selfNode.ServerType, selfNode.ServerID), string(buf), 10)
	// 设置监听 + 发现其他服务
	if err := dis.Watch(ROOT_DIR, watchServerNode); err != nil {
		return err
	}
	return nil
}

func watchServerNode(action int, key string, value string) {
	vv := &pb.ServerNode{}
	proto.Unmarshal(fbasic.StrToBytes(value), vv)
	switch action {
	case domain.DELETE:
		serverType, serverID, _ := ParseChannel(key)
		nodes.Delete(serverType, serverID)
	default:
		// 添加服务节点
		nodes.Add(vv)
	}
}

// 对玩家路由
func Dispatcher(head *pb.PacketHead) error {
	// 从路由表中更新
	table := router.GetRouteList(head.UID)
	if item := table.GetRouteInfo(int32(head.DstServerType)); item != nil {
		head.DstServerID = item.ServerID
	}
	// 判断服务节点是否存在
	if node := nodes.Get(head.DstServerType, head.DstServerID); node != nil {
		return nil
	}
	// 重新路由服务节点
	node := nodes.Random(head)
	if node == nil {
		return uerror.NewUErrorf(1, -1, "%s not found", head.DstServerType.String())
	}
	// 更新路由表
	head.DstServerID = node.ServerID
	table.UpdateRouteInfo(head, node)
	return nil
}
