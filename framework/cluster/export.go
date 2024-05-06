package cluster

import (
	"universal/common/pb"
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

// 对玩家路由
func Dispatcher(head *pb.PacketHead) error {
	return service.Dispatcher(head)
}
