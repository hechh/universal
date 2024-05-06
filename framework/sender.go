package framework

import (
	"universal/common/pb"
	"universal/framework/cluster"
	"universal/framework/fbasic"
	"universal/framework/notify"

	"google.golang.org/protobuf/proto"
)

// 发送到其他服务
func SendTo(sendType pb.SendType, apiCode int32, uid uint64, req proto.Message) error {
	self := cluster.GetLocalClusterNode()
	head := &pb.PacketHead{
		Status:         pb.StatusType_REQUEST,
		SendType:       sendType,
		SrcClusterType: self.ClusterType,
		SrcClusterID:   self.ClusterID,
		DstClusterType: fbasic.ApiCodeToClusterType(apiCode),
		ApiCode:        apiCode,
		Time:           fbasic.GetNow(),
		UID:            uid,
	}
	// 路由
	if err := cluster.Dispatcher(head); err != nil {
		return err
	}
	pp, err := fbasic.ReqToPacket(head, req)
	if err != nil {
		return err
	}
	return notify.Publish(pp)
}

// 发送客户端
func SendToClient(sendType pb.SendType, apiCode int32, uid uint64, rsp proto.Message) error {
	self := cluster.GetLocalClusterNode()
	head := &pb.PacketHead{
		Status:         pb.StatusType_RESPONSE,
		SendType:       sendType,
		SrcClusterType: self.ClusterType,
		SrcClusterID:   self.ClusterID,
		DstClusterType: pb.ClusterType_GATE,
		ApiCode:        apiCode,
		Time:           fbasic.GetNow(),
		UID:            uid,
	}
	// 路由
	if err := cluster.Dispatcher(head); err != nil {
		return err
	}
	pp, err := fbasic.RspToPacket(head, rsp)
	if err != nil {
		return err
	}
	return notify.Publish(pp)
}
