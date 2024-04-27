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

func SocketClientCheck(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_SocketClientCheck, fmt.Sprint(args...))
}

func SocketClientHeaderSize(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_SocketClientHeaderSize, fmt.Sprint(args...))
}

func SocketClientBodySizeLimit(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_SocketClientBodySizeLimit, fmt.Sprint(args...))
}

func SocketClientMaxLimit(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_SocketClientMaxLimit, fmt.Sprint(args...))
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

func NatsBuildClient(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_NatsBuildClient, fmt.Sprint(args...))
}

func NatsPublish(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_NatsPublish, fmt.Sprint(args...))
}

func ApiCodeNotFound(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_ApiCodeNotFound, fmt.Sprint(args...))
}

func NotFound(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_NotFound, fmt.Sprint(args...))
}
