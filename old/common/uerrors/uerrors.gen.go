package uerrors

import (
	"universal/common/pb"
	"universal/framework/common/uerror"
)

func Success(args ...interface{}) *uerror.UError {
	return uerror.NewUError(1, int32(pb.ErrorCode_Success), args...)
}
