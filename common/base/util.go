package base

import (
	"universal/common/pb"
	"universal/library/uerror"

	"github.com/golang/protobuf/proto"
)

func ToRspHead(err error) *pb.RspHead {
	switch vv := err.(type) {
	case *uerror.UError:
		return &pb.RspHead{Code: vv.GetCode(), ErrMsg: vv.GetMsg()}
	case nil:
		return nil
	default:
		return &pb.RspHead{Code: -1, ErrMsg: err.Error()}
	}
}

func Unmarshal[T any](buf []byte, req *T) error {
	switch vv := any(req).(type) {
	case proto.Message:
		return proto.Unmarshal(buf, vv)
	}
	return nil
}
