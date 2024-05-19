package handler

import (
	"reflect"
	"universal/common/pb"
	"universal/framework/common/fbasic"

	"google.golang.org/protobuf/proto"
)

type Func struct {
	name       string         // 函数名
	isVariadic bool           // 是否为可变参数
	handle     reflect.Value  // 函数
	params     []reflect.Type // 参数
	returns    []reflect.Type // 返回值
}

func NewFunc(f interface{}) *Func {
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
	return &Func{
		name:       fbasic.GetFuncName(v),
		isVariadic: t.IsVariadic(),
		handle:     v,
		params:     params,
		returns:    returns,
	}
}

func (d *Func) FuncName() string {
	return d.name
}

func (d *Func) GetReturns() []reflect.Type {
	return d.returns
}

func (d *Func) Call(ctx *fbasic.Context, req, rsp proto.Message) (err error) {
	// 解析参数
	newReq := req.(*pb.ActorRequest)
	params := fbasic.DecodeValue(newReq.Buff, d.params, 0)
	// 执行函数
	var results []reflect.Value
	if !d.isVariadic {
		results = d.handle.Call(params)
	} else {
		results = d.handle.CallSlice(params)
	}
	// 返回
	newRsp := rsp.(*pb.ActorResponse)
	newRsp.Buff = fbasic.EncodeValue(results...)
	return
}
