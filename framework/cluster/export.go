package cluster

import (
	"universal/common/pb"
	"universal/framework/basic"
	"universal/framework/cluster/domain"
	"universal/framework/cluster/internal/service"

	"google.golang.org/protobuf/proto"
)

func InitCluster(node *pb.ClusterNode, natsUrl string, etcds []string) error {
	return service.InitCluster(node, natsUrl, etcds)
}

func Subscribe(h domain.ClusterFunc) {
	service.GetCluster().Subscribe(h)
}

// 转发到game服务集群
func SendToGame(head *pb.PacketHead, params ...interface{}) error {
	node := service.GetCluster().GetDiscovery().GetSelf()
	head.SrcClusterID = node.ClusterID
	head.SrcClusterType = node.ClusterType

	head.DstClusterType = pb.ClusterType_GAME
	head.SendType = pb.SendType_POINT
	pac, err := basic.ToReqPacket(head, params...)
	if err != nil {
		return err
	}
	return service.GetCluster().Send(pac)
}

// 转发到db服务集群
func SendToDb(head *pb.PacketHead, params ...interface{}) error {
	node := service.GetCluster().GetDiscovery().GetSelf()
	head.SrcClusterID = node.ClusterID
	head.SrcClusterType = node.ClusterType

	head.DstClusterType = pb.ClusterType_DB
	head.SendType = pb.SendType_POINT
	pac, err := basic.ToReqPacket(head, params...)
	if err != nil {
		return err
	}
	return service.GetCluster().Send(pac)
}

// 转发到db服务集群
func SendToGate(head *pb.PacketHead, params ...interface{}) error {
	node := service.GetCluster().GetDiscovery().GetSelf()
	head.SrcClusterID = node.ClusterID
	head.SrcClusterType = node.ClusterType

	head.DstClusterType = pb.ClusterType_GATE
	head.SendType = pb.SendType_POINT
	pac, err := basic.ToReqPacket(head, params...)
	if err != nil {
		return err
	}
	return service.GetCluster().Send(pac)
}

// 主动发送到客户端
func SendToClient(head *pb.PacketHead, rsp proto.Message, err error) error {
	node := service.GetCluster().GetDiscovery().GetSelf()
	head.SrcClusterID = node.ClusterID
	head.SrcClusterType = node.ClusterType

	head.DstClusterType = pb.ClusterType_GATE
	head.DstClusterID = 0

	head.ActorName = ""
	head.FuncName = ""
	return service.GetCluster().Send(basic.ToRspPacket(head, err, rsp))
}
