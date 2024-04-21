package repository

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"universal/common/pb"
	"universal/framework/basic"
)

type FuncPacket struct {
	name       string         // 函数名
	isVariadic bool           // 是否为可变参数
	handle     reflect.Value  // 函数
	params     []reflect.Type // 参数
}

func NewFuncPacket(f interface{}) *FuncPacket {
	v := reflect.ValueOf(f)
	t := v.Type()
	// 获取函数参数
	params := make([]reflect.Type, t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		params[i] = t.In(i)
	}
	// 返回数据
	return &FuncPacket{
		name:       basic.GetFuncName(v),
		isVariadic: t.IsVariadic(),
		handle:     v,
		params:     params,
	}
}

func (d *FuncPacket) GetFuncName() string {
	return d.name
}

func (d *FuncPacket) Call(ctx *basic.Context, pac *pb.Packet) *pb.Packet {
	// 解析参数
	params := make([]reflect.Value, len(d.params))
	decode := gob.NewDecoder(bytes.NewReader(pac.Buff))
	for i, paramType := range d.params {
		params[i] = reflect.New(paramType).Elem()
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
		basic.ToErrorPacket(result, results[ll-1].Interface().(error))
		if ll > 1 {
			result.Buff = basic.ToGobBytes(results[:ll-1])
		}
	}
	return result
}
