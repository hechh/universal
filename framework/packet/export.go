package packet

import (
	"universal/common/pb"
	"universal/framework/fbasic"
	"universal/framework/packet/domain"
	"universal/framework/packet/internal/manager"

	"google.golang.org/protobuf/proto"
)

func RegisterApi(apiCode int32, h domain.ApiFunc, req, rsp proto.Message) {
	manager.RegisterApi(apiCode, h, req, rsp)
}

func RegisterStruct(apiCode int32, h interface{}) {
	manager.RegisterStruct(apiCode, h)
}

func RegisterFunc(apiCode int32, h interface{}) {
	manager.RegisterFunc(apiCode, h)
}

func Call(ctx *fbasic.Context, buf []byte) (*pb.Packet, error) {
	return manager.Call(ctx, buf)
}

// 解析返回值参数
func ParseReturns(apiCode int32, actorName, funcName string, buf []byte) ([]interface{}, error) {
	return manager.ParseReturns(apiCode, actorName, funcName, buf)
}
