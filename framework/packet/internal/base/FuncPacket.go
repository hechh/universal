package base

import (
	"reflect"
	"universal/common/pb"
	"universal/framework/fbasic"

	"google.golang.org/protobuf/proto"
)

type FuncPacket struct {
	name       string         // 函数名
	isVariadic bool           // 是否为可变参数
	handle     reflect.Value  // 函数
	params     []reflect.Type // 参数
	returns    []reflect.Type // 返回值
}

func NewFuncPacket(f interface{}) *FuncPacket {
	v := reflect.ValueOf(f)
	t := v.Type()
	// 获取函数参数
	params := make([]reflect.Type, t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		params[i] = t.In(i)
	}
	// 函数返回值
	returns := make([]reflect.Type, t.NumOut())
	for i := 0; i < t.NumOut(); i++ {
		returns[i] = t.Out(i)
	}
	// 返回数据
	return &FuncPacket{
		name:       fbasic.GetFuncName(v),
		isVariadic: t.IsVariadic(),
		handle:     v,
		params:     params,
		returns:    returns,
	}
}

func (d *FuncPacket) FuncName() string {
	return d.name
}

func (d *FuncPacket) GetReturns() []reflect.Type {
	return d.returns
}

func (d *FuncPacket) Call(ctx *fbasic.Context, req, rsp proto.Message) (err error) {
	// 解析参数
	params := make([]reflect.Value, len(d.params))
	newReq := req.(*pb.ActorRequest)
	fbasic.DecodeTypes(d.params).DecodeValues(newReq.Buff, params)
	// 执行函数
	var results []reflect.Value
	if !d.isVariadic {
		results = d.handle.Call(params)
	} else {
		results = d.handle.CallSlice(params)
	}
	// 返回
	newRsp := rsp.(*pb.ActorResponse)
	newRsp.Buff = fbasic.EncodeValues(results).Encode()
	return
}
