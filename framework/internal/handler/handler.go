package handler

import (
	"reflect"
	"universal/framework/define"
	"universal/library/baselib/uerror"
)

var (
	cmds = make(map[uint32]*Handler)
)

type Handler struct {
	f   define.HandleFunc
	req reflect.Type
	rsp reflect.Type
}

func (d *Handler) Handle(ctx define.IContext, req define.IProto, rsp define.IProto) error {
	return d.f(ctx, req, rsp)
}

func (d *Handler) NewReq() define.IProto {
	return reflect.New(d.req).Interface().(define.IProto)
}

func (d *Handler) NewRsp() define.IProto {
	return reflect.New(d.rsp).Interface().(define.IProto)
}

func Register(cmd uint32, f define.HandleFunc, req define.IProto, rsp define.IProto) error {
	if _, ok := cmds[cmd]; ok {
		return uerror.New(1, -1, "cmd:%d already register", cmd)
	}
	if f == nil {
		return uerror.New(1, -1, "cmd:%d, handler is nil", cmd)
	}
	if req == nil {
		return uerror.New(1, -1, "cmd:%d, request is nil", cmd)
	}
	if rsp == nil {
		return uerror.New(1, -1, "cmd:%d, response is nil", cmd)
	}
	cmds[cmd] = &Handler{
		f:   f,
		req: reflect.TypeOf(req).Elem(),
		rsp: reflect.TypeOf(rsp).Elem(),
	}
	return nil
}

func GetHeader(cmd uint32) *Handler {
	if h, ok := cmds[cmd]; ok {
		return h
	}
	return nil
}
