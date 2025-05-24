package domain

import (
	"github.com/golang/protobuf/proto"
)

type Result struct {
	UID      uint64
	Cost     int64
	Error    error
	Response proto.Message
}

type ResultCB func(*Result)

func DefaultResult(_ *Result) {}
