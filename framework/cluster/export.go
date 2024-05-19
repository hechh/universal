package cluster

import (
	"universal/common/pb"
	"universal/framework/cluster/internal/nodes"
	"universal/framework/cluster/internal/router"
	"universal/framework/cluster/internal/service"
)

func Close() {
	service.Close()
}

// 初始化连接
func Init(etcds []string) error {
	return service.Init(etcds)
}

// 设置路由过期时间
func SetRouterClearExpire(expire int64) {
	router.SetClearExpire(expire)
}

// 获取本地节点
func GetSelfServerNode() *pb.ServerNode {
	return service.GetSelfServerNode()
}

// 服务发现
func Discovery(typ pb.ServerType, addr string) error {
	return service.Discovery(typ, addr)
}

// 随机路由一个服务节点
func RandomNode(head *pb.PacketHead) *pb.ServerNode {
	return nodes.Random(head)
}

// 获取服务节点信息
func GetNode(srvType pb.ServerType, srvID uint32) *pb.ServerNode {
	return nodes.Get(srvType, srvID)
}

// 玩家消息
func GetPlayerChannel(typ pb.ServerType, serverID uint32, uid uint64) string {
	return service.GetPlayerChannel(typ, serverID, uid)
}

// 服务节点消息
func GetNodeChannel(typ pb.ServerType, serverID uint32) string {
	return service.GetNodeChannel(typ, serverID)
}

// 所有节点消息
func GetClusterChannel(typ pb.ServerType) string {
	return service.GetClusterChannel(typ)
}

// 获取channel的key
func GetHeadChannel(head *pb.PacketHead) (str string, err error) {
	return service.GetHeadChannel(head)
}

// 路由
func Dispatcher(head *pb.PacketHead) error {
	return service.Dispatcher(head)
}
