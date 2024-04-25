package repository

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"universal/common/pb"
	"universal/framework/basic"
)

type StructPacket struct {
	stName     string         // 结构名字
	name       string         // 函数名
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
			stName:     v.Name(),
			name:       m.Name,
			isVariadic: m.Type.IsVariadic(),
			handle:     m.Func,
			params:     params,
		})
	}
	return
}

func (d *StructPacket) GetFuncName() string {
	return d.name
}

func (d *StructPacket) GetStructName() string {
	return d.stName
}

func (d *StructPacket) Call(ctx *basic.Context, pac *pb.Packet) (*pb.Packet, error) {
	params := make([]reflect.Value, len(d.params))
	if val := ctx.GetValue(d.stName); val != nil {
		params[0] = reflect.ValueOf(val)
	} else {
		return &pb.Packet{Head: pac.Head}, nil
	}
	// 解析参数
	decode := gob.NewDecoder(bytes.NewReader(pac.Buff))
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
	result := &pb.Packet{Head: ctx.PacketHead}
	if ll := len(results); ll > 0 {
		/*
			basic.ToErrorPacket(result, results[ll-1].Interface().(error))
			if ll > 1 {
				result.Buff = basic.ToGobBytes(results[:ll-1])
			}
		*/
	}
	return result, nil
}
