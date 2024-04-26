package uerrors

import (
	"universal/common/pb"
	"universal/framework/fbasic"
)

func Success() *fbasic.UError {
	return fbasic.NewUError(1, pb.ErrorCode_Success, "Success")
}

func Unmarshal() *fbasic.UError {
	return fbasic.NewUError(1, pb.ErrorCode_Unmarshal, "Unmarshal")
}

func NotSupported() *fbasic.UError {
	return fbasic.NewUError(1, pb.ErrorCode_NotSupported, "NotSupported")
}

func BuildEtcdClient() *fbasic.UError {
	return fbasic.NewUError(1, pb.ErrorCode_BuildEtcdClient, "BuildEtcdClient")
}

func NotExist() *fbasic.UError {
	return fbasic.NewUError(1, pb.ErrorCode_NotExist, "NotExist")
}

func Parameter() *fbasic.UError {
	return fbasic.NewUError(1, pb.ErrorCode_Parameter, "Parameter")
}

func NotFound() *fbasic.UError {
	return fbasic.NewUError(1, pb.ErrorCode_NotFound, "NotFound")
}

func Marhsal() *fbasic.UError {
	return fbasic.NewUError(1, pb.ErrorCode_Marhsal, "Marhsal")
}

func NatsPublish() *fbasic.UError {
	return fbasic.NewUError(1, pb.ErrorCode_NatsPublish, "NatsPublish")
}

func NewClient() *fbasic.UError {
	return fbasic.NewUError(1, pb.ErrorCode_NewClient, "NewClient")
}

func SocketClose() *fbasic.UError {
	return fbasic.NewUError(1, pb.ErrorCode_SocketClose, "SocketClose")
}

func NatsSubscribe() *fbasic.UError {
	return fbasic.NewUError(1, pb.ErrorCode_NatsSubscribe, "NatsSubscribe")
}

func ActorNameNotFound() *fbasic.UError {
	return fbasic.NewUError(1, pb.ErrorCode_ActorNameNotFound, "ActorNameNotFound")
}

func ApiCodeNotFound() *fbasic.UError {
	return fbasic.NewUError(1, pb.ErrorCode_ApiCodeNotFound, "ApiCodeNotFound")
}
