package framework

import (
	"universal/common/pb"
	"universal/framework/cluster"
	"universal/framework/fbasic"
	"universal/framework/notify"

	"google.golang.org/protobuf/proto"
)

func SendTo(uid uint64, sendType pb.SendType, apiCode int32, req proto.Message) error {
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
