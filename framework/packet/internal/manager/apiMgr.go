package manager

import (
	"fmt"
	"universal/common/pb"
	"universal/framework/basic"
	"universal/framework/packet/domain"
	"universal/framework/packet/internal/base"

	"google.golang.org/protobuf/proto"
)

var (
	apiPool = make(map[int32]domain.IPacket)
)

func RegisterApi(apiCode int32, h domain.ApiFunc, req, rsp proto.Message) {
	attr := base.NewApiPacket(h, req, rsp)
	if _, ok := apiPool[apiCode]; ok {
		panic(fmt.Sprintf("ApiCode(%d) has already registered", apiCode))
	}
	apiPool[apiCode] = attr
}

func RegisterStruct(apiCode int32, h interface{}) {
	if _, ok := apiPool[apiCode]; !ok {
		apiPool[apiCode] = base.NewActorPacket(&pb.ActorRequest{}, &pb.ActorResponse{})
	}
	mgr, ok := apiPool[apiCode].(*base.ActorPacket)
	if !ok {
		panic(fmt.Sprintf("%d is not ActorPacket", apiCode))
	}
	mgr.RegisterStruct(h)
}

func RegisterFunc(apiCode int32, h interface{}) {
	if _, ok := apiPool[apiCode]; !ok {
		apiPool[apiCode] = base.NewActorPacket(&pb.ActorRequest{}, &pb.ActorResponse{})
	}
	mgr, ok := apiPool[apiCode].(*base.ActorPacket)
	if !ok {
		panic(fmt.Sprintf("%d is not ActorPacket", apiCode))
	}
	mgr.RegisterFunc(h)
}

func Call(ctx *basic.Context, buf []byte) (*pb.Packet, error) {
	val, ok := apiPool[ctx.ApiCode]
	if !ok {
		return nil, basic.NewUError(1, pb.ErrorCode_ApiCodeNotFound, ctx.String())
	}
	return val.Call(ctx, buf)
}
