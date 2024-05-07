package cluster

import (
	"universal/common/pb"
	"universal/framework/cluster/internal/nodes"
	"universal/framework/cluster/internal/service"
)

// 初始化连接
func Init(etcds []string, types ...pb.ClusterType) error {
	return service.Init(etcds, types...)
}

// 获取本地节点
func GetLocalClusterNode() *pb.ClusterNode {
	return service.GetLocalClusterNode()
}

// 服务发现
func Discovery(typ pb.ClusterType, addr string) error {
	return service.Discovery(typ, addr)
}

// 随机路由一个服务节点
func RandomNode(head *pb.PacketHead) *pb.ClusterNode {
	return nodes.Random(head)
}

// 获取服务节点信息
func GetNode(clusterType pb.ClusterType, clusterID uint32) *pb.ClusterNode {
	return nodes.Get(clusterType, clusterID)
}
