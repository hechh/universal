package player

import (
	"universal/common/pb"
	"universal/framework/common/fbasic"
	"universal/framework/common/ulog"
	"universal/framework/packet"

	"google.golang.org/protobuf/proto"
)

func LoginRequest(ctx *fbasic.Context, req, rsp proto.Message) error {
	ulog.Debug(1, "ctx: %v, req: %v", ctx, req)
	return nil
}

func init() {
	packet.RegisterApi(pb.ApiCode_GATE_LOGIN_REQUEST, LoginRequest, &pb.GateLoginRequest{}, &pb.GateLoginResponse{})
}
