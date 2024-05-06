package base

import (
	"reflect"
	"universal/common/pb"
	"universal/framework/fbasic"
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
		name:   fbasic.GetFuncName(h),
		handle: h,
		req:    reflect.TypeOf(req).Elem(),
		rsp:    reflect.TypeOf(rsp).Elem(),
	}
}

func (d *ApiPacket) GetFuncName() string {
	return d.name
}

func (d *ApiPacket) Call(ctx *fbasic.Context, buf []byte) (newRsp proto.Message) {
	// 设置rsp的head
	rspHead := new(pb.RpcHead)
	newRsp = reflect.New(d.rsp).Interface().(proto.Message)
	if vv := reflect.ValueOf(newRsp).Elem().Field(3); vv.IsNil() {
		vv.Set(reflect.ValueOf(rspHead))
	}
	// 解析req请求
	newReq := reflect.New(d.req).Interface().(proto.Message)
	if err := proto.Unmarshal(buf, newReq); err != nil {
		uerr := fbasic.NewUError(1, pb.ErrorCode_ProtoUnmarshal, err)
		rspHead.Code, rspHead.ErrMsg = uerr.GetCode(), uerr.GetErrMsg()
		return
	}
	// 执行API
	if err := d.handle(ctx, newReq, newRsp); err != nil {
		rspHead.Code, rspHead.ErrMsg = fbasic.GetCodeMsg(err)
	}
	return
}
