package service

import (
	"universal/common/pb"
	"universal/framework/common/uerror"
	"universal/framework/network/domain"
	"universal/framework/network/internal/base"
	"universal/framework/network/internal/middle"

	"google.golang.org/protobuf/proto"
)

var (
	client domain.IMiddle
)

func InitMiddle(typ int, url string) error {
	switch typ {
	case domain.NetworkTypeNats:
		if cli, err := middle.NewNatsClient(url); err != nil {
			return err
		} else {
			client = cli
		}
	default:
		return uerror.NewUErrorf(1, -1, "Network type not supported, type: %d", typ)
	}
	return nil
}

// 消息订阅
func Subscribe(key string, f func(*pb.Packet)) error {
	return client.Subscribe(key, f)
}

// 发送消息
func Publish(key string, pac *pb.Packet) error {
	return client.Publish(key, pac)
}

func PublishReq(key string, head *pb.PacketHead, req proto.Message, params ...interface{}) error {
	// 封装发送包
	pac, err := base.ReqToPacket(head, req, params...)
	if err != nil {
		return err
	}
	// 发送
	return Publish(key, pac)
}

func PublishRsp(key string, head *pb.PacketHead, rsp proto.Message, params ...interface{}) error {
	// 封装发送包
	pac, err := base.RspToPacket(head, rsp, params...)
	if err != nil {
		return err
	}
	// 发送
	return Publish(key, pac)
}
