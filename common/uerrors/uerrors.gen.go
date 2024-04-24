package uerrors

import (
	"universal/common/pb"
	"universal/framework/basic"
)

func Success() *basic.UError {
	return basic.NewUError(3, pb.ErrorCode_Success, "Success")
}

func Unmarshal() *basic.UError {
	return basic.NewUError(3, pb.ErrorCode_Unmarshal, "Unmarshal")
}

func NotSupported() *basic.UError {
	return basic.NewUError(3, pb.ErrorCode_NotSupported, "NotSupported")
}

func BuildEtcdClient() *basic.UError {
	return basic.NewUError(3, pb.ErrorCode_BuildEtcdClient, "BuildEtcdClient")
}

func NotExist() *basic.UError {
	return basic.NewUError(3, pb.ErrorCode_NotExist, "NotExist")
}

func Parameter() *basic.UError {
	return basic.NewUError(3, pb.ErrorCode_Parameter, "Parameter")
}

func NotFound() *basic.UError {
	return basic.NewUError(3, pb.ErrorCode_NotFound, "NotFound")
}

func Marhsal() *basic.UError {
	return basic.NewUError(3, pb.ErrorCode_Marhsal, "Marhsal")
}

func NatsPublish() *basic.UError {
	return basic.NewUError(3, pb.ErrorCode_NatsPublish, "NatsPublish")
}

func NewClient() *basic.UError {
	return basic.NewUError(3, pb.ErrorCode_NewClient, "NewClient")
}

func SocketClose() *basic.UError {
	return basic.NewUError(3, pb.ErrorCode_SocketClose, "SocketClose")
}

func NatsSubscribe() *basic.UError {
	return basic.NewUError(3, pb.ErrorCode_NatsSubscribe, "NatsSubscribe")
}
