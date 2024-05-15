package service

import (
	"net"
	"universal/common/pb"
	"universal/framework/cluster/domain"
	"universal/framework/cluster/internal/discovery/etcd"
	"universal/framework/cluster/internal/nodes"
	"universal/framework/fbasic"

	"github.com/spf13/cast"
	"google.golang.org/protobuf/proto"
)

var (
	dis      domain.IDiscovery // 服务发现
	selfNode *pb.ClusterNode   // 服务自身节点
)

func GetLocalClusterNode() *pb.ClusterNode {
	return selfNode
}

func GetDiscovery() domain.IDiscovery {
	return dis
}

func Stop() {
	dis.Close()
}

// 初始化
func Init(ends []string, types ...pb.ClusterType) (err error) {
	// 初始化etcd
	etcd, err := etcd.NewEtcdClient(ends...)
	if err != nil {
		return err
	}
	// 初始化节点类型
	nodes.Init(types...)
	// 设置服务发现
	dis = etcd
	return
}

func watchClusterNode(action int, key string, value string) {
	vv := &pb.ClusterNode{}
	if err := proto.Unmarshal(fbasic.StringToBytes(value), vv); err != nil {
		panic(err)
	}
	switch action {
	case domain.ActionTypeDel:
		clusterType, clusterID, _ := fbasic.ParseChannel(key)
		nodes.Delete(clusterType, clusterID)
	default:
		// 添加服务节点
		nodes.Add(vv)
	}
}

// 服务发现
func Discovery(typ pb.ClusterType, addr string) error {
	ip, port, err := net.SplitHostPort(addr)
	if err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_SocketAddr, err)
	}
	// 构建自身节点
	selfNode = &pb.ClusterNode{
		ClusterType: typ,
		ClusterID:   fbasic.GetCrc32(addr),
		Ip:          ip,
		Port:        cast.ToInt32(port),
	}

	// 注册自身服务（保活，服务下线会自动删除）
	buf, err := proto.Marshal(selfNode)
	if err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_ProtoMarshal, err)
	}
	dis.KeepAlive(fbasic.GetNodeChannel(selfNode.ClusterType, selfNode.ClusterID), string(buf), 10)

	// 设置监听 + 发现其他服务
	if err := dis.Watch(fbasic.GetRootDir(), watchClusterNode); err != nil {
		return err
	}
	return nil
}
