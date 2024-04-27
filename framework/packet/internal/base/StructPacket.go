package base

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"universal/common/pb"
	"universal/framework/fbasic"

	"google.golang.org/protobuf/proto"
)

type StructPacket struct {
	actorName  string         // 结构名字
	fname      string         // 函数名
	isVariadic bool           // 是否为可变参数
	handle     reflect.Value  // 函数
	params     []reflect.Type // 参数
}

func NewStructPacket(f interface{}) (rets []*StructPacket) {
	v := reflect.ValueOf(f).Type()
	for i := 0; i < v.NumMethod(); i++ {
		m := v.Method(i)
		tFunc := m.Func.Type()
		params := make([]reflect.Type, tFunc.NumIn())
		for i := 0; i < tFunc.NumIn(); i++ {
			params[i] = tFunc.In(i)
		}
		// 返回函数
		rets = append(rets, &StructPacket{
			actorName:  v.Elem().Name(),
			fname:      m.Name,
			isVariadic: m.Type.IsVariadic(),
			handle:     m.Func,
			params:     params,
		})
	}
	return
}

func (d *StructPacket) FuncName() string {
	return d.fname
}

func (d *StructPacket) ActorName() string {
	return d.actorName
}

func (d *StructPacket) Call(ctx *fbasic.Context, req, rsp proto.Message) (err error) {
	params := make([]reflect.Value, len(d.params))
	obj := ctx.GetValue(d.actorName)
	if obj == nil {
		return fbasic.NewUError(1, pb.ErrorCode_ActorNotSupported, d.actorName)
	}
	params[0] = reflect.ValueOf(obj)
	// 解析参数
	newReq := req.(*pb.ActorRequest)
	decode := gob.NewDecoder(bytes.NewReader(newReq.Buff))
	for i := 1; i < len(d.params); i++ {
		params[i] = reflect.New(d.params[i]).Elem()
		decode.DecodeValue(params[i])
	}
	// 执行函数
	var results []reflect.Value
	if !d.isVariadic {
		results = d.handle.Call(params)
	} else {
		results = d.handle.CallSlice(params)
	}
	// 返回
	var ok bool
	newRsp := rsp.(*pb.ActorResponse)
	if ll := len(results); ll > 0 {
		if err, ok = results[ll-1].Interface().(error); ok {
			newRsp.Buff = fbasic.ToGobBytes(results[:ll-1])
		} else {
			newRsp.Buff = fbasic.ToGobBytes(results)
		}
	}
	return
}
