package service

import (
	"universal/common/pb"
	"universal/framework/fbasic"
	"universal/framework/notify/domain"
	"universal/framework/notify/internal/middle/nats"

	"google.golang.org/protobuf/proto"
)

var (
	client domain.IMiddle
)

func Init(typ int, url string) error {
	switch typ {
	case domain.NotifyTypeNats:
		if cli, err := nats.NewNatsClient(url); err != nil {
			return err
		} else {
			client = cli
		}
	default:
		return fbasic.NewUError(1, pb.ErrorCode_NotSupported, typ)
	}
	return nil
}

// 消息订阅
func Subscribe(key string, f domain.NotifyHandle) error {
	return client.Subscribe(key, f)
}

// 发送消息
func Publish(key string, pac *pb.Packet) error {
	return client.Publish(key, pac)
}

func PublishReq(key string, head *pb.PacketHead, req proto.Message, params ...interface{}) error {
	// 封装发送包
	pac, err := fbasic.ReqToPacket(head, req, params...)
	if err != nil {
		return err
	}
	// 发送
	return Publish(key, pac)
}

func PublishRsp(key string, head *pb.PacketHead, rsp proto.Message, params ...interface{}) error {
	// 封装发送包
	pac, err := fbasic.RspToPacket(head, rsp, params...)
	if err != nil {
		return err
	}
	// 发送
	return Publish(key, pac)
}
