package domain

import "universal/common/pb"

const (
	NotifyTypeNats  = 1
	NotifyTypeKafka = 2
)

type NotifyHandle func(*pb.Packet)

// 消息中间件的发布、订阅接口
type IMiddle interface {
	Subscribe(string, NotifyHandle) error
	Publish(string, *pb.Packet) error
}
