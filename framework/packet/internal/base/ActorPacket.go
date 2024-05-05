package base

import (
	"fmt"
	"reflect"
	"universal/common/pb"
	"universal/framework/fbasic"
	"universal/framework/packet/domain"

	"google.golang.org/protobuf/proto"
)

type index struct {
	ActorName string
	FuncName  string
}

type ActorPacket struct {
	req  reflect.Type
	rsp  reflect.Type
	apis map[index]domain.IApi
}

func NewActorPacket(req, rsp proto.Message) *ActorPacket {
	return &ActorPacket{
		req:  reflect.TypeOf(req).Elem(),
		rsp:  reflect.TypeOf(rsp).Elem(),
		apis: make(map[index]domain.IApi),
	}
}

func (d *ActorPacket) GetReturns(actorName, funcName string) ([]reflect.Type, error) {
	index := index{actorName, funcName}
	val, ok := d.apis[index]
	if !ok {
		return nil, fbasic.NewUError(1, pb.ErrorCode_ActorNotSupported, actorName, funcName)
	}
	return val.GetReturns(), nil
}

func (d *ActorPacket) RegisterStruct(st interface{}) {
	for _, attr := range NewStructPacket(st) {
		index := index{attr.ActorName(), attr.FuncName()}
		if _, ok := d.apis[index]; ok {
			panic(fmt.Sprintf("%s.%s() has already registered", attr.ActorName(), attr.FuncName()))
		}
		d.apis[index] = attr
	}
}

func (d *ActorPacket) RegisterFunc(f interface{}) {
	attr := NewFuncPacket(f)
	index := index{"", attr.FuncName()}
	if _, ok := d.apis[index]; ok {
		panic(fmt.Sprintf("%s() has already registered", attr.FuncName()))
	}
	d.apis[index] = attr
}

func (d *ActorPacket) Call(ctx *fbasic.Context, buf []byte) (*pb.Packet, error) {
	// 解析请求参数
	newReq := reflect.New(d.req).Interface().(proto.Message)
	if err := proto.Unmarshal(buf, newReq); err != nil {
		return nil, fbasic.NewUError(1, pb.ErrorCode_ProtoUnmarshal, err)
	}
	// 获取api
	req := newReq.(*pb.ActorRequest)
	api, ok := d.apis[index{req.ActorName, req.FuncName}]
	if !ok {
		return nil, fbasic.NewUError(1, pb.ErrorCode_ActorNotSupported, newReq)
	}
	// 设置rsp的head
	newRsp := reflect.New(d.rsp).Interface().(proto.Message)
	if vv := reflect.ValueOf(newRsp).Elem().Field(3); vv.IsNil() {
		vv.Set(reflect.ValueOf(&pb.RpcHead{}))
	}
	// 执行API
	err := api.Call(ctx, newReq, newRsp)
	// 执行函数
	return fbasic.RspToPacket(ctx.PacketHead, err, newRsp)
}
