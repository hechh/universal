package uerrors

import (
	"fmt"
	"universal/common/pb"
	"universal/framework/fbasic"
)

func Success(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_Success, fmt.Sprint(args...))
}

func ProtoMarshal(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_ProtoMarshal, fmt.Sprint(args...))
}

func ProtoUnmarshal(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_ProtoUnmarshal, fmt.Sprint(args...))
}

func JsonMarshal(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_JsonMarshal, fmt.Sprint(args...))
}

func JsonUnmarshal(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_JsonUnmarshal, fmt.Sprint(args...))
}

func Parameter(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_Parameter, fmt.Sprint(args...))
}

func SocketClientSend(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_SocketClientSend, fmt.Sprint(args...))
}

func SocketClientRead(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_SocketClientRead, fmt.Sprint(args...))
}

func SocketFrameCheck(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_SocketFrameCheck, fmt.Sprint(args...))
}

func SocketFrameHeaderSize(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_SocketFrameHeaderSize, fmt.Sprint(args...))
}

func SocketFrameBodySizeMaxLimit(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_SocketFrameBodySizeMaxLimit, fmt.Sprint(args...))
}

func SocketFrameSizeMaxLimit(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_SocketFrameSizeMaxLimit, fmt.Sprint(args...))
}

func Unknown(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_Unknown, fmt.Sprint(args...))
}

func ActorNotSupported(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_ActorNotSupported, fmt.Sprint(args...))
}

func ActorHasRegistered(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_ActorHasRegistered, fmt.Sprint(args...))
}

func EtcdBuildClient(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_EtcdBuildClient, fmt.Sprint(args...))
}

func EtcdClientGet(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_EtcdClientGet, fmt.Sprint(args...))
}

func EtcdClientPut(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_EtcdClientPut, fmt.Sprint(args...))
}

func EtcdLeaseCreate(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_EtcdLeaseCreate, fmt.Sprint(args...))
}

func EtcdLeaseKeepAliveOnce(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_EtcdLeaseKeepAliveOnce, fmt.Sprint(args...))
}

func NatsBuildClient(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_NatsBuildClient, fmt.Sprint(args...))
}

func NatsPublish(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_NatsPublish, fmt.Sprint(args...))
}

func NatsSubscribe(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_NatsSubscribe, fmt.Sprint(args...))
}

func ApiCodeNotFound(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_ApiCodeNotFound, fmt.Sprint(args...))
}

func NotFound(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_NotFound, fmt.Sprint(args...))
}

func TypeNotSupported(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_TypeNotSupported, fmt.Sprint(args...))
}

func SendTypeNotSupported(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_SendTypeNotSupported, fmt.Sprint(args...))
}

func ClusterNodeNotFound(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_ClusterNodeNotFound, fmt.Sprint(args...))
}

func ArrayOutOfBounds(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_ArrayOutOfBounds, fmt.Sprint(args...))
}
