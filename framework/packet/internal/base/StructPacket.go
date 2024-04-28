package base

import (
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
	returns    []reflect.Type // 返回值
}

func NewStructPacket(f interface{}) (rets []*StructPacket) {
	v := reflect.ValueOf(f).Type()
	for i := 0; i < v.NumMethod(); i++ {
		m := v.Method(i)
		tFunc := m.Func.Type()
		// 函数参数
		params := make([]reflect.Type, tFunc.NumIn())
		for i := 0; i < tFunc.NumIn(); i++ {
			params[i] = tFunc.In(i)
		}
		// 函数返回值
		returns := make([]reflect.Type, tFunc.NumOut())
		for i := 0; i < tFunc.NumOut(); i++ {
			returns[i] = tFunc.Out(i)
		}
		// 返回函数
		rets = append(rets, &StructPacket{
			actorName:  v.Elem().Name(),
			fname:      m.Name,
			isVariadic: m.Type.IsVariadic(),
			handle:     m.Func,
			params:     params,
			returns:    returns,
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

func (d *StructPacket) GetReturns() []reflect.Type {
	return d.returns
}

func (d *StructPacket) Call(ctx *fbasic.Context, req, rsp proto.Message) (err error) {
	// 设置this指针
	params := make([]reflect.Value, len(d.params))
	obj := ctx.GetValue(d.actorName)
	if obj == nil {
		return fbasic.NewUError(1, pb.ErrorCode_ActorNotSupported, d.actorName)
	}
	params[0] = reflect.ValueOf(obj)
	// 解析参数
	newReq := req.(*pb.ActorRequest)
	fbasic.DecodeTypes(d.params[1:]).DecodeValues(newReq.Buff, params[1:])
	// 执行函数
	var results []reflect.Value
	if !d.isVariadic {
		results = d.handle.Call(params)
	} else {
		results = d.handle.CallSlice(params)
	}
	// 返回函数返回值
	newRsp := rsp.(*pb.ActorResponse)
	newRsp.Buff = fbasic.EncodeValues(results).Encode()
	return
}
