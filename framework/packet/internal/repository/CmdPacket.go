package repository

import (
	"reflect"
	"universal/common/pb"
	"universal/framework/basic"
	"universal/framework/packet/domain"

	"google.golang.org/protobuf/proto"
)

type CmdPacket struct {
	name   string
	handle domain.CmdFunc
	req    reflect.Type
	rsp    reflect.Type
}

func NewCmdPacket(h domain.CmdFunc, req, rsp proto.Message) *CmdPacket {
	return &CmdPacket{
		name:   basic.GetFuncName(h),
		handle: h,
		req:    reflect.TypeOf(req).Elem(),
		rsp:    reflect.TypeOf(rsp).Elem(),
	}
}

func (d *CmdPacket) GetFuncName() string {
	return d.name
}

func (d *CmdPacket) Call(ctx *basic.Context, pac *pb.Packet) *pb.Packet {
	newRsp := reflect.New(d.rsp).Interface().(proto.Message)
	newReq := reflect.New(d.req).Interface().(proto.Message)
	if err := proto.Unmarshal(pac.Buff, newReq); err != nil {
		return basic.ToRspPacket(ctx.PacketHead, basic.NewUError(2, pb.ErrorCode_Unmarshal, err), newRsp)
	}
	// 执行函数
	return basic.ToRspPacket(ctx.PacketHead, d.handle(ctx, newReq, newRsp), newRsp)
}
