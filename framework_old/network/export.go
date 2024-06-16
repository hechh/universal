package network

import (
	"universal/common/pb"
	"universal/framework/network/domain"
	"universal/framework/network/internal/base"
	"universal/framework/network/internal/service"

	"google.golang.org/protobuf/proto"
)

// nats
func InitNats(url string) error {
	return service.InitMiddle(domain.NetworkTypeNats, url)
}

func Subscribe(key string, f func(*pb.Packet)) error {
	return service.Subscribe(key, f)
}

func Publish(key string, pac *pb.Packet) error {
	return service.Publish(key, pac)
}

func PublishReq(key string, head *pb.PacketHead, req proto.Message, params ...interface{}) error {
	return service.PublishReq(key, head, req, params...)
}

func PublishRsp(key string, head *pb.PacketHead, rsp proto.Message, params ...interface{}) error {
	return service.PublishRsp(key, head, rsp, params...)
}

func ReqToPacket(head *pb.PacketHead, req proto.Message, params ...interface{}) (*pb.Packet, error) {
	return base.ReqToPacket(head, req, params...)
}

func RspToPacket(head *pb.PacketHead, rsp proto.Message, params ...interface{}) (*pb.Packet, error) {
	return base.RspToPacket(head, rsp, params...)
}
