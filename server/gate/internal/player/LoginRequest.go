package player

import (
	"log"
	"universal/common/pb"
	"universal/framework/common/fbasic"
	"universal/framework/packet"

	"google.golang.org/protobuf/proto"
)

func LoginRequest(ctx *fbasic.Context, req, rsp proto.Message) error {
	log.Println("-----LoginRequest-------", ctx, req, rsp)
	return nil
}

func init() {
	packet.RegisterApi(pb.ApiCode_GATE_LOGIN_REQUEST, LoginRequest, &pb.GateLoginRequest{}, &pb.GateLoginResponse{})
}
