package domain

import "universal/common/pb"

const (
	NotifyTypeNats  = 1
	NotifyTypeKafka = 2
)

// 消息处理接口
type IHandle interface {
	GetKey() string    // 获取key
	Handle(*pb.Packet) // 消息处理
}

type NotifyHandle func(*pb.Packet)

// 消息中间件的发布、订阅接口
type IMiddle interface {
	Subscribe(string, NotifyHandle) error
	Publish(string, *pb.Packet) error
}
