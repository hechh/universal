package actor

import (
	"reflect"
	"strings"
	"universal/framework/domain"
	"universal/library/baselib/uerror"
	"universal/library/encode"
	"universal/library/mlog"

	"github.com/golang/protobuf/proto"
)

type Actor struct {
	*Async
	name  string
	rval  reflect.Value
	funcs map[string]*FuncInfo
}

func (a *Actor) GetActorName() string {
	return a.name
}

func (a *Actor) Register(ac domain.IActor) {
	a.Async = NewAsync()
	a.rval = reflect.ValueOf(ac)
	name := a.rval.Elem().Type().Name()
	if index := strings.Index(name, "."); index > -1 {
		name = name[index+1:]
	}
	a.name = name
}

func (d *Actor) ParseFunc(tt interface{}) {
	switch vv := tt.(type) {
	case map[string]*FuncInfo:
		d.funcs = vv
	case reflect.Type:
		d.funcs = parseFuncs(vv)
	default:
		panic("注册参数错误，必须是方法列表或reflect.Type")
	}
}

func (d *Actor) Send(h domain.IHead, args ...interface{}) error {
	mm, ok := d.funcs[h.GetFuncName()]
	if !ok {
		return uerror.New(1, -1, "%s.%s未实现", h.GetActorName(), h.GetFuncName())
	}
	if mm.isVariadic {
		d.Push(handleVariadic(d.rval, mm, h, args...))
	} else {
		d.Push(handle(d.rval, mm, h, args...))
	}
	return nil
}

func (d *Actor) SendRpc(h domain.IHead, buf []byte) error {
	mm, ok := d.funcs[h.GetFuncName()]
	if !ok {
		return uerror.New(1, -1, "%s.%s未实现", h.GetActorName(), h.GetFuncName())
	}
	// 发送事件
	if mm.isProto {
		d.Push(handRpcProto(d.rval, mm, h, buf))
	} else {
		d.Push(handRpcGob(d.rval, mm, h, buf))
	}
	return nil
}

func handle(rval reflect.Value, mm *FuncInfo, head domain.IHead, args ...interface{}) func() {
	return func() {
		in := make([]reflect.Value, mm.Type.NumIn())
		// 设置 this
		in[0] = rval
		// 设置 head
		pos := 1
		if mm.hasHead {
			in[1] = reflect.ValueOf(head)
			pos++
		}
		// 设置参数
		for i := pos; i < mm.Type.NumIn(); i++ {
			in[i] = reflect.ValueOf(args[i-pos])
		}
		// 调用函数
		result := mm.Func.Call(in)
		// 处理返回值
		if mm.hasError {
			if result[0].IsNil() {
				return
			}
			mlog.Error("调用%s.%s报错：%v", head.GetActorName(), head.GetFuncName(), result[0].Interface())
		}
	}
}

// 可变参数
func handleVariadic(rval reflect.Value, mm *FuncInfo, head domain.IHead, args ...interface{}) func() {
	return func() {
		in := make([]reflect.Value, mm.Type.NumIn())
		// 设置 this
		in[0] = rval
		// 设置 head
		pos := 1
		if mm.hasHead {
			in[1] = reflect.ValueOf(head)
			pos++
		}
		// 设置参数
		for i := pos; i < mm.Type.NumIn()-1; i++ {
			in[i] = reflect.ValueOf(args[i-pos])
		}
		// 设置可变参数
		args = args[mm.Type.NumIn()-pos-1:]
		if len(args) > 0 {
			arr := make([]reflect.Value, len(args))
			for i, item := range args {
				arr[i] = reflect.ValueOf(item)
			}
			in[mm.Type.NumIn()-1] = reflect.ValueOf(arr)
		}
		// 调用函数
		result := mm.Func.CallSlice(in)
		// 处理返回值
		if mm.hasError {
			if result[0].IsNil() {
				return
			}
			mlog.Error("调用%s.%s报错：%v", head.GetActorName(), head.GetFuncName(), result[0].Interface())
		}
	}
}

func handRpcProto(rval reflect.Value, mm *FuncInfo, head domain.IHead, buf []byte) func() {
	return func() {
		in := make([]reflect.Value, mm.Type.NumIn())
		// 设置 this
		in[0] = rval
		// 设置 head
		pos := 1
		if mm.hasHead {
			in[1] = reflect.ValueOf(head)
			pos++
		}
		// 解析参数
		for i := pos; i < mm.Type.NumIn(); i++ {
			req := reflect.New(mm.Type.In(i).Elem())
			if i == pos {
				if err := proto.Unmarshal(buf, req.Interface().(proto.Message)); err != nil {
					mlog.Error("%s.%s参数解析报错: %v", head.GetActorName(), head.GetFuncName(), err)
					return
				}
			}
			in[i] = reflect.ValueOf(req)
		}
		// 调用函数
		result := mm.Func.Call(in)
		if mm.hasError {
			if result[0].IsNil() {
				return
			}
			mlog.Error("调用%s.%s报错：%v", head.GetActorName(), head.GetFuncName(), result[0].Interface())
		}
	}
}

func handRpcGob(rval reflect.Value, mm *FuncInfo, head domain.IHead, buf []byte) func() {
	return func() {
		pos := 1
		if mm.hasHead {
			pos++
		}
		// 解析参数参数
		in, err := encode.Decode(buf, mm.Method, pos)
		if err != nil {
			mlog.Error("%s.%s参数解析错误: %v", head.GetActorName(), head.GetFuncName(), err)
		}
		// 设置 this
		in[0] = rval
		// 设置 head
		if mm.hasHead {
			in[1] = reflect.ValueOf(head)
		}
		// 调用函数
		var result []reflect.Value
		if mm.isVariadic {
			result = mm.Func.CallSlice(in)
		} else {
			result = mm.Func.Call(in)
		}
		// 处理返回值
		if mm.hasError {
			if result[0].IsNil() {
				return
			}
			mlog.Error("调用%s.%s报错：%v", head.GetActorName(), head.GetFuncName(), result[0].Interface())
		}
	}
}
