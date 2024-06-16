package domain

import "universal/common/pb"

const (
	NetworkTypeNats  = 1
	NetworkTypeKafka = 2
)

// 消息中间件的发布、订阅接口
type IMiddle interface {
	Subscribe(string, func(*pb.Packet)) error
	Publish(string, *pb.Packet) error
}
