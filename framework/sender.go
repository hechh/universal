package framework

import (
	"universal/common/pb"
	"universal/framework/cluster"
	"universal/framework/fbasic"
	"universal/framework/notify"
	"universal/framework/routine"

	"google.golang.org/protobuf/proto"
)

// 发送到其他服务
func SendTo(sendType pb.SendType, apiCode int32, uid uint64, req proto.Message) error {
	self := cluster.GetLocalClusterNode()
	head := &pb.PacketHead{
		SendType:       sendType,
		SrcClusterType: self.ClusterType,
		SrcClusterID:   self.ClusterID,
		DstClusterType: fbasic.ApiCodeToClusterType(apiCode),
		ApiCode:        apiCode,
		Time:           fbasic.GetNow(),
		UID:            uid,
	}
	// 路由
	if err := Dispatcher(head); err != nil {
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
		SendType:       sendType,
		SrcClusterType: self.ClusterType,
		SrcClusterID:   self.ClusterID,
		DstClusterType: pb.ClusterType_GATE,
		ApiCode:        apiCode + 1,
		Time:           fbasic.GetNow(),
		UID:            uid,
	}
	// 路由
	if err := Dispatcher(head); err != nil {
		return err
	}
	pp, err := fbasic.RspToPacket(head, rsp)
	if err != nil {
		return err
	}
	return notify.Publish(pp)
}

// 对玩家路由
func Dispatcher(head *pb.PacketHead) error {
	rlist := routine.GetRoutine(head.UID)
	if rinfo := rlist.Get(head.DstClusterType); rinfo == nil {
		// 路由
		if err := rlist.UpdateRoutine(head, cluster.RandomNode(head)); err != nil {
			return err
		}
	} else {
		if head.DstClusterID <= 0 {
			head.DstClusterID = rinfo.ClusterID
		}
		// 节点丢失
		if head.DstClusterID != rinfo.ClusterID {
			// 重新路由
			if err := rlist.UpdateRoutine(head, cluster.RandomNode(head)); err != nil {
				return err
			}
		}
		// 判断节点是否存在
		if node := cluster.GetNode(head.DstClusterType, head.DstClusterID); node == nil {
			// 重新路由
			if err := rlist.UpdateRoutine(head, cluster.RandomNode(head)); err != nil {
				return err
			}
		} else {
			rlist.UpdateRoutine(head, node)
		}
	}
	return nil
}
