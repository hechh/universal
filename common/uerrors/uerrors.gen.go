package uerrors

import (
	"universal/common/pb"
	"universal/framework/fbasic"
)

func Success(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_Success, args...)
}

func ProtoMarshal(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_ProtoMarshal, args...)
}

func ProtoUnmarshal(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_ProtoUnmarshal, args...)
}

func JsonMarshal(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_JsonMarshal, args...)
}

func JsonUnmarshal(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_JsonUnmarshal, args...)
}

func Parameter(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_Parameter, args...)
}

func SocketClientSend(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_SocketClientSend, args...)
}

func SocketClientRead(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_SocketClientRead, args...)
}

func SocketFrameCheck(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_SocketFrameCheck, args...)
}

func SocketFrameHeaderSize(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_SocketFrameHeaderSize, args...)
}

func SocketFrameBodySizeMaxLimit(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_SocketFrameBodySizeMaxLimit, args...)
}

func SocketFrameSizeMaxLimit(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_SocketFrameSizeMaxLimit, args...)
}

func Unknown(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_Unknown, args...)
}

func ActorNotSupported(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_ActorNotSupported, args...)
}

func ActorHasRegistered(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_ActorHasRegistered, args...)
}

func EtcdBuildClient(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_EtcdBuildClient, args...)
}

func EtcdClientGet(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_EtcdClientGet, args...)
}

func EtcdClientPut(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_EtcdClientPut, args...)
}

func EtcdClientDelete(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_EtcdClientDelete, args...)
}

func EtcdLeaseCreate(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_EtcdLeaseCreate, args...)
}

func EtcdLeaseKeepAliveOnce(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_EtcdLeaseKeepAliveOnce, args...)
}

func NatsBuildClient(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_NatsBuildClient, args...)
}

func NatsPublish(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_NatsPublish, args...)
}

func NatsSubscribe(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_NatsSubscribe, args...)
}

func ApiCodeNotFound(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_ApiCodeNotFound, args...)
}

func TypeNotSupported(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_TypeNotSupported, args...)
}

func SendTypeNotSupported(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_SendTypeNotSupported, args...)
}

func ClusterNodeNotFound(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_ClusterNodeNotFound, args...)
}

func ArrayOutOfBounds(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_ArrayOutOfBounds, args...)
}

func ReadYaml(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_ReadYaml, args...)
}

func YamlUnmarshal(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_YamlUnmarshal, args...)
}

func NotSupported(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_NotSupported, args...)
}

func NotFound(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_NotFound, args...)
}

func SocketAddr(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_SocketAddr, args...)
}

func GateLoginRequestExpected(args ...interface{}) *fbasic.UError {
	return fbasic.NewUError(2, pb.ErrorCode_GateLoginRequestExpected, args...)
}
