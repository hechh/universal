package domain

import (
	"universal/common/pb"

	"github.com/golang/protobuf/proto"
)

type IHead interface {
	GetPacketHead() *pb.IPacket
}

type Result struct {
	UID      uint64
	Cost     uint64
	Error    error
	Response proto.Message
}

type ResultCallBack func(*Result)
