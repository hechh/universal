package handler

import (
	"reflect"
	"universal/common/pb"
	"universal/framework/common/fbasic"
	"universal/framework/common/uerror"
	"universal/framework/packet/domain"

	"google.golang.org/protobuf/proto"
)

type ApiHandler struct {
	name   string
	handle domain.ApiFunc
	req    reflect.Type
	rsp    reflect.Type
}

func NewApiHandler(h domain.ApiFunc, req, rsp proto.Message) *ApiHandler {
	return &ApiHandler{
		name:   fbasic.GetFuncName(h),
		handle: h,
		req:    reflect.TypeOf(req).Elem(),
		rsp:    reflect.TypeOf(rsp).Elem(),
	}
}

func (d *ApiHandler) GetFuncName() string {
	return d.name
}

func (d *ApiHandler) Call(ctx *fbasic.Context, buf []byte) (newRsp proto.Message) {
	// 设置rsp的head
	rspHead := new(pb.RpcHead)
	newRsp = reflect.New(d.rsp).Interface().(proto.Message)
	if vv := reflect.ValueOf(newRsp).Elem().Field(3); vv.IsNil() {
		vv.Set(reflect.ValueOf(rspHead))
	}
	// 解析req请求
	newReq := reflect.New(d.req).Interface().(proto.Message)
	if err := proto.Unmarshal(buf, newReq); err != nil {
		uerr := uerror.NewUError(1, -1, err)
		rspHead.Code, rspHead.ErrMsg = uerr.GetCode(), uerr.GetErrMsg()
		return
	}
	// 执行API
	if err := d.handle(ctx, newReq, newRsp); err != nil {
		rspHead.Code, rspHead.ErrMsg = uerror.GetCodeMsg(err)
	}
	return
}
