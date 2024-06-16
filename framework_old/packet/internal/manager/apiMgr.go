package manager

import (
	"fmt"
	"universal/common/pb"
	"universal/framework/common/fbasic"
	"universal/framework/common/uerror"
	"universal/framework/packet/domain"
	"universal/framework/packet/internal/handler"

	"google.golang.org/protobuf/proto"
)

var (
	apiPool = make(map[int32]domain.IHandler)
)

func RegisterApi(apiCode int32, h domain.ApiFunc, req, rsp proto.Message) {
	attr := handler.NewApiHandler(h, req, rsp)
	if _, ok := apiPool[apiCode]; ok {
		panic(fmt.Sprintf("ApiCode(%d) has already registered", apiCode))
	}
	apiPool[apiCode] = attr
}

func RegisterStruct(apiCode int32, h interface{}) {
	if _, ok := apiPool[apiCode]; !ok {
		apiPool[apiCode] = handler.NewActorHandler(&pb.ActorRequest{}, &pb.ActorResponse{})
	}
	mgr, ok := apiPool[apiCode].(*handler.ActorHandler)
	if !ok {
		panic(fmt.Sprintf("%d is not ActorHandler", apiCode))
	}
	mgr.RegisterStruct(h)
}

func RegisterFunc(apiCode int32, h interface{}) {
	if _, ok := apiPool[apiCode]; !ok {
		apiPool[apiCode] = handler.NewActorHandler(&pb.ActorRequest{}, &pb.ActorResponse{})
	}
	mgr, ok := apiPool[apiCode].(*handler.ActorHandler)
	if !ok {
		panic(fmt.Sprintf("%d is not ActorHandler", apiCode))
	}
	mgr.RegisterFunc(h)
}

func Call(ctx *fbasic.Context, buf []byte) (proto.Message, error) {
	val, ok := apiPool[ctx.ApiCode]
	if !ok {
		return nil, uerror.NewUErrorf(1, -1, "ApiCode(%d) not supported", ctx.ApiCode)
	}
	return val.Call(ctx, buf), nil
}

func ParseReturns(apiCode int32, actorName, funcName string, buf []byte) ([]interface{}, error) {
	mgr, ok := apiPool[apiCode].(*handler.ActorHandler)
	if !ok || mgr == nil {
		return nil, uerror.NewUErrorf(1, -1, "ApiCode: %d, ActorName: %s, FuncName: %s", apiCode, actorName, funcName)
	}
	// 获取返回值类型
	typs, err := mgr.GetReturns(actorName, funcName)
	if err != nil {
		return nil, err
	}
	// 解析
	return fbasic.DecodeAny(buf, typs, 0), nil
}
