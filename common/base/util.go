package base

import (
	"universal/common/pb"
	"universal/library/uerror"
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
