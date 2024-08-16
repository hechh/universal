package domain

import (
	"github.com/golang/protobuf/proto"
)

type Result struct {
	UID      uint64
	Cost     uint64
	Error    error
	Response proto.Message
}

type ResultCallBack func(*Result)
