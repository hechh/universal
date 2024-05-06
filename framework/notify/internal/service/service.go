package service

import (
	"universal/common/pb"
	"universal/framework/fbasic"
	"universal/framework/notify/domain"
	"universal/framework/notify/internal/middle/nats"
)

var (
	notify domain.INotify
)

func Init(typ int, url string) error {
	switch typ {
	case domain.NotifyTypeNats:
		if client, err := nats.NewNatsClient(url); err != nil {
			return err
		} else {
			notify = client
		}
	default:
		return fbasic.NewUError(1, pb.ErrorCode_NotSupported, typ)
	}
	return nil
}

// 消息订阅
func Subscribe(key string, f func(*pb.Packet)) error {
	return notify.Subscribe(key, f)
}

// 发送消息
func Publish(key string, pac *pb.Packet) error {
	return notify.Publish(key, pac)
}
