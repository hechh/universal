package notify

import (
	"universal/common/pb"
	"universal/framework/fbasic"
	"universal/framework/notify/domain"

	"universal/framework/notify/internal/service"
)

func Init(url string) error {
	return service.Init(domain.NotifyTypeNats, url)
}

// 消息订阅
func Subscribe(key string, f func(*pb.Packet)) error {
	return service.Subscribe(key, f)
}

// 发送
func Publish(pac *pb.Packet) error {
	head := pac.Head
	switch head.SendType {
	case pb.SendType_PLAYER:
		return service.Publish(fbasic.GetPlayerChannel(head.DstClusterType, head.DstClusterID, head.UID), pac)
	case pb.SendType_NODE:
		return service.Publish(fbasic.GetNodeChannel(head.DstClusterType, head.DstClusterID), pac)
	case pb.SendType_CLUSTER:
		return service.Publish(fbasic.GetClusterChannel(head.DstClusterType), pac)
	default:
		return fbasic.NewUError(1, pb.ErrorCode_SendTypeNotSupported, head.SendType)
	}
	return nil
}
