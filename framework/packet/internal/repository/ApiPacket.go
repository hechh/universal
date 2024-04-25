package repository

import (
	"reflect"
	"universal/common/pb"
	"universal/framework/basic"
	"universal/framework/packet/domain"

	"google.golang.org/protobuf/proto"
)

type ApiPacket struct {
	name   string
	handle domain.ApiFunc
	req    reflect.Type
	rsp    reflect.Type
}

func NewApiPacket(h domain.ApiFunc, req, rsp proto.Message) *ApiPacket {
	return &ApiPacket{
		name:   basic.GetFuncName(h),
		handle: h,
		req:    reflect.TypeOf(req).Elem(),
		rsp:    reflect.TypeOf(rsp).Elem(),
	}
}

func (d *ApiPacket) GetFuncName() string {
	return d.name
}

func (d *ApiPacket) Call(ctx *basic.Context, pac *pb.Packet) (*pb.Packet, error) {
	newReq := reflect.New(d.req).Interface().(proto.Message)
	newRsp := reflect.New(d.rsp).Interface().(proto.Message)

	if err := proto.Unmarshal(pac.Buff, newReq); err != nil {
		return nil, basic.NewUError(1, pb.ErrorCode_Unmarshal, err)
	}
	// 执行API
	err := d.handle(ctx, newReq, newRsp)
	// 执行函数
	return basic.RspToPacket(ctx.PacketHead, err, newRsp)
}
