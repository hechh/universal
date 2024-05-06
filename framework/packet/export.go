package packet

import (
	"universal/common/pb"
	"universal/framework/fbasic"
	"universal/framework/packet/domain"
	"universal/framework/packet/internal/manager"

	"google.golang.org/protobuf/proto"
)

func RegisterApi(apiCode pb.ApiCode, h domain.ApiFunc, req, rsp proto.Message) {
	manager.RegisterApi(int32(apiCode), h, req, rsp)
}

func RegisterStruct(apiCode pb.ApiCode, h interface{}) {
	manager.RegisterStruct(int32(apiCode), h)
}

func RegisterFunc(apiCode pb.ApiCode, h interface{}) {
	manager.RegisterFunc(int32(apiCode), h)
}

func Call(ctx *fbasic.Context, buf []byte) (proto.Message, error) {
	return manager.Call(ctx, buf)
}

// 解析返回值参数
func ParseReturns(apiCode pb.ApiCode, actorName, funcName string, buf []byte) ([]interface{}, error) {
	return manager.ParseReturns(int32(apiCode), actorName, funcName, buf)
}
