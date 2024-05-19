package handler

import (
	"fmt"
	"reflect"
	"universal/common/pb"
	"universal/framework/common/fbasic"
	"universal/framework/common/uerror"

	"google.golang.org/protobuf/proto"
)

type IApi interface {
	GetReturns() []reflect.Type
	Call(*fbasic.Context, proto.Message, proto.Message) error
}

type index struct {
	ActorName string
	FuncName  string
}

type ActorHandler struct {
	req  reflect.Type
	rsp  reflect.Type
	apis map[index]IApi
}

func NewActorHandler(req, rsp proto.Message) *ActorHandler {
	return &ActorHandler{
		req:  reflect.TypeOf(req).Elem(),
		rsp:  reflect.TypeOf(rsp).Elem(),
		apis: make(map[index]IApi),
	}
}

func (d *ActorHandler) GetReturns(actorName, funcName string) ([]reflect.Type, error) {
	index := index{actorName, funcName}
	val, ok := d.apis[index]
	if !ok {
		return nil, uerror.NewUErrorf(1, -1, "%s.%s not supported", actorName, funcName)
	}
	return val.GetReturns(), nil
}

func (d *ActorHandler) RegisterStruct(st interface{}) {
	for _, attr := range NewStruct(st) {
		index := index{attr.ActorName(), attr.FuncName()}
		if _, ok := d.apis[index]; ok {
			panic(fmt.Sprintf("%s.%s() has already registered", attr.ActorName(), attr.FuncName()))
		}
		d.apis[index] = attr
	}
}

func (d *ActorHandler) RegisterFunc(f interface{}) {
	attr := NewFunc(f)
	index := index{"", attr.FuncName()}
	if _, ok := d.apis[index]; ok {
		panic(fmt.Sprintf("%s() has already registered", attr.FuncName()))
	}
	d.apis[index] = attr
}

func (d *ActorHandler) Call(ctx *fbasic.Context, buf []byte) (newRsp proto.Message) {
	// 设置rsp的head
	rspHead := new(pb.RpcHead)
	newRsp = reflect.New(d.rsp).Interface().(proto.Message)
	if vv := reflect.ValueOf(newRsp).Elem().Field(3); vv.IsNil() {
		vv.Set(reflect.ValueOf(rspHead))
	}
	// 解析请求参数
	newReq := reflect.New(d.req).Interface().(proto.Message)
	if err := proto.Unmarshal(buf, newReq); err != nil {
		uerr := uerror.NewUError(1, -1, err)
		rspHead.Code, rspHead.ErrMsg = uerr.GetCode(), uerr.GetErrMsg()
		return
	}
	// 获取api
	req := newReq.(*pb.ActorRequest)
	api, ok := d.apis[index{req.ActorName, req.FuncName}]
	if !ok {
		uerr := uerror.NewUError(1, -1, newReq)
		rspHead.Code, rspHead.ErrMsg = uerr.GetCode(), uerr.GetErrMsg()
		return
	}
	// 执行API
	if err := api.Call(ctx, newReq, newRsp); err != nil {
		rspHead.Code, rspHead.ErrMsg = uerror.GetCodeMsg(err)
	}
	return
}
