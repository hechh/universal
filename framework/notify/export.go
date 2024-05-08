package notify

import (
	"universal/common/pb"
	"universal/framework/notify/domain"
	"universal/framework/notify/internal/base"
	"universal/framework/notify/internal/service"

	"google.golang.org/protobuf/proto"
)

func Init(url string) error {
	return service.Init(domain.NotifyTypeNats, url)
}

// 消息订阅
func Subscribe(key string, f func(*pb.Packet)) error {
	return service.Subscribe(key, f)
}

// 发送
func Publish(key string, pac *pb.Packet) error {
	return service.Publish(key, pac)
}

func PublishReq(key string, head *pb.PacketHead, req proto.Message, params ...interface{}) error {
	return service.PublishReq(key, head, req, params...)
}

func PublishRsp(key string, head *pb.PacketHead, rsp proto.Message, params ...interface{}) error {
	return service.PublishRsp(key, head, rsp, params...)
}

// 创建一个广播
func NewBroadcast(key string) *base.Broadcast {
	return base.NewBroadcast(key)
}

// 创建一个单波
func NewUnicast(key string, f domain.NotifyHandle) *base.Unicast {
	return base.NewUnicast(key, f)
}
