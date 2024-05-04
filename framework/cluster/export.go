package cluster

import (
	"universal/common/pb"
	"universal/framework/cluster/internal/service"
	"universal/framework/fbasic"

	"google.golang.org/protobuf/proto"
)

// 初始化连接
func Init(natsUrl string, etcds []string) error {
	return service.Init(natsUrl, etcds)
}

// 服务发现
func Discovery(typ pb.ClusterType, addr string) error {
	return service.Discovery(typ, addr)
}

// 订阅消息
func Subscribe(h func(*pb.Packet)) error {
	return service.Subscribe(h)
}

// 夸服务发送消息
func Publish(pac *pb.Packet) error {
	return service.Publish(pac)
}

// 转发到game服务集群
func SendToGame(head *pb.PacketHead, params ...interface{}) error {
	node := service.GetLocalClusterNode()
	head.SrcClusterID = node.ClusterID
	head.SrcClusterType = node.ClusterType

	head.DstClusterType = pb.ClusterType_GAME
	head.SendType = pb.SendType_POINT
	pac, err := fbasic.ReqToPacket(head, params...)
	if err != nil {
		return err
	}
	return service.Publish(pac)
}

// 转发到db服务集群
func SendToDb(head *pb.PacketHead, params ...interface{}) error {
	node := service.GetLocalClusterNode()
	head.SrcClusterID = node.ClusterID
	head.SrcClusterType = node.ClusterType

	head.DstClusterType = pb.ClusterType_DB
	head.SendType = pb.SendType_POINT
	pac, err := fbasic.ReqToPacket(head, params...)
	if err != nil {
		return err
	}
	return service.Publish(pac)
}

// 转发到db服务集群
func SendToGate(head *pb.PacketHead, params ...interface{}) error {
	node := service.GetLocalClusterNode()
	head.SrcClusterID = node.ClusterID
	head.SrcClusterType = node.ClusterType

	head.DstClusterType = pb.ClusterType_GATE
	head.SendType = pb.SendType_POINT
	pac, err := fbasic.ReqToPacket(head, params...)
	if err != nil {
		return err
	}
	return service.Publish(pac)
}

// 主动发送到客户端
func SendToClient(head *pb.PacketHead, rsp proto.Message, err error) error {
	node := service.GetLocalClusterNode()
	head.SrcClusterID = node.ClusterID
	head.SrcClusterType = node.ClusterType

	head.DstClusterType = pb.ClusterType_GATE
	head.DstClusterID = 0
	pac, err := fbasic.RspToPacket(head, err, rsp)
	if err != nil {
		return err
	}
	return service.Publish(pac)
}
