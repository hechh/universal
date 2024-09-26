package test

import (
	"corps/base/cfgEnum"
	"corps/pb"
	"fmt"

	"github.com/golang/protobuf/proto"
)

type TestPacket struct{}

// 发送给客户端
func (d *TestPacket) SendToClient(head *pb.RpcHead, packet proto.Message, uCode cfgEnum.ErrorCode) {
	fmt.Println(head, packet, uCode)
}

// 通用返回给客户端
func (d *TestPacket) SendCommonToClient(head *pb.RpcHead, uErrorCode cfgEnum.ErrorCode) {
	fmt.Println(head, uErrorCode)
}
