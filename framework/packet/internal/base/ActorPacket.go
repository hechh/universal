package base

import (
	"fmt"
	"reflect"
	"universal/common/pb"
	"universal/framework/basic"
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

func (d *ActorPacket) Call(ctx *basic.Context, buf []byte) (*pb.Packet, error) {
	newReq := reflect.New(d.req).Interface().(proto.Message)
	newRsp := reflect.New(d.rsp).Interface().(proto.Message)
	if err := proto.Unmarshal(buf, newReq); err != nil {
		return nil, basic.NewUError(1, pb.ErrorCode_Unmarshal, err)
	}
	// 获取index
	req := newReq.(*pb.ActorRequest)
	index := index{req.ActorName, req.FuncName}
	api, ok := d.apis[index]
	if !ok {
		return nil, basic.NewUError(1, pb.ErrorCode_ActorNameNotFound, fmt.Sprintf("%v", index))
	}
	// 执行API
	err := api.Call(ctx, newReq, newRsp)
	// 执行函数
	return basic.RspToPacket(ctx.PacketHead, err, newRsp)
}
