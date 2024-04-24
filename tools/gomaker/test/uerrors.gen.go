package uerrors

import (
	"universal/common/pb"
	"universal/framework/basic"
)

func Success() *basic.UError {
	return basic.NewUError(3, test.ErrorCode_Success, "Success")
}

func Unmarshal() *basic.UError {
	return basic.NewUError(3, test.ErrorCode_Unmarshal, "Unmarshal")
}

func NotSupported() *basic.UError {
	return basic.NewUError(3, test.ErrorCode_NotSupported, "NotSupported")
}

func BuildEtcdClient() *basic.UError {
	return basic.NewUError(3, test.ErrorCode_BuildEtcdClient, "BuildEtcdClient")
}

func NotExist() *basic.UError {
	return basic.NewUError(3, test.ErrorCode_NotExist, "NotExist")
}

func Parameter() *basic.UError {
	return basic.NewUError(3, test.ErrorCode_Parameter, "Parameter")
}

func NotFound() *basic.UError {
	return basic.NewUError(3, test.ErrorCode_NotFound, "NotFound")
}

func Marhsal() *basic.UError {
	return basic.NewUError(3, test.ErrorCode_Marhsal, "Marhsal")
}

func NatsPublish() *basic.UError {
	return basic.NewUError(3, test.ErrorCode_NatsPublish, "NatsPublish")
}

func NewClient() *basic.UError {
	return basic.NewUError(3, test.ErrorCode_NewClient, "NewClient")
}

func SocketClose() *basic.UError {
	return basic.NewUError(3, test.ErrorCode_SocketClose, "SocketClose")
}

func NatsSubscribe() *basic.UError {
	return basic.NewUError(3, test.ErrorCode_NatsSubscribe, "NatsSubscribe")
}
