package actor

import (
	"reflect"
	"universal/common/pb"
	"universal/framework/library/encode"
	"universal/framework/library/mlog"

	"github.com/golang/protobuf/proto"
)

var (
	headType  = reflect.TypeOf((**pb.Head)(nil)).Elem()
	errorType = reflect.TypeOf((*error)(nil)).Elem()
	protoType = reflect.TypeOf((*proto.Message)(nil)).Elem()
)

type FuncInfo struct {
	reflect.Method
	hasHead    bool
	hasError   bool
	isVariadic bool
	isProto    bool
}

func NewFuncInfo(m reflect.Method) *FuncInfo {
	hasError := false
	if m.Type.NumOut() > 0 && m.Type.Out(0).Implements(errorType) {
		hasError = true
	}
	pos := 1
	hasHead := false
	if m.Type.NumIn() > 1 && m.Type.In(1).Implements(headType) {
		pos++
		hasHead = true
	}
	isProto := true
	for i := pos; i < m.Type.NumIn(); i++ {
		if m.Type.In(i).Implements(protoType) {
			continue
		} else {
			isProto = false
		}
		break
	}
	return &FuncInfo{
		Method:   m,
		hasHead:  hasHead,
		hasError: hasError,
		isProto:  isProto,
	}
}

// 非可变参数调用
func (f *FuncInfo) handle(rval reflect.Value, head *pb.Head, args ...interface{}) func() {
	return func() {
		in := make([]reflect.Value, f.Type.NumIn())
		in[0] = rval
		pos := 1
		if f.hasHead {
			pos++
			in[pos] = reflect.ValueOf(head)
		}
		for i := pos; i < f.Type.NumIn(); i++ {
			in[i] = reflect.ValueOf(args[i-pos])
		}
		// 调用函数
		result := f.Func.Call(in)
		if f.hasError {
			if result[0].IsNil() {
				return
			}
			mlog.Error("调用%s.%s报错：%v", head.ActorName, head.FuncName, result[0].Interface())
		}
	}
}

// 可变参数调用
func (f *FuncInfo) handleVariadic(rval reflect.Value, head *pb.Head, args ...interface{}) func() {
	return func() {
		in := make([]reflect.Value, f.Type.NumIn())
		in[0] = rval
		pos := 1
		if f.hasHead {
			pos++
			in[1] = reflect.ValueOf(head)
		}
		for i := pos; i < f.Type.NumIn()-1; i++ {
			in[i] = reflect.ValueOf(args[i-pos])
		}
		if args2 := args[f.Type.NumIn()-pos-1:]; len(args2) > 0 {
			arr := make([]reflect.Value, len(args2))
			for i, item := range args2 {
				arr[i] = reflect.ValueOf(item)
			}
			in[f.Type.NumIn()-1] = reflect.ValueOf(arr)
		}
		// 调用函数
		result := f.Func.CallSlice(in)
		if f.hasError {
			if result[0].IsNil() {
				return
			}
			mlog.Error("调用%s.%s报错：%v", head.ActorName, head.FuncName, result[0].Interface())
		}
	}
}

// 远程调用
func (f *FuncInfo) handleRpcProto(rval reflect.Value, head *pb.Head, buf []byte) func() {
	return func() {
		in := make([]reflect.Value, f.Type.NumIn())
		in[0] = rval
		pos := 1
		if f.hasHead {
			pos++
			in[1] = reflect.ValueOf(head)
		}
		for i := pos; i < f.Type.NumIn(); i++ {
			req := reflect.New(f.Type.In(i).Elem())
			if i == pos {
				if err := proto.Unmarshal(buf, req.Interface().(proto.Message)); err != nil {
					mlog.Error("%s.%s参数解析报错: %v", head.ActorName, head.FuncName, err)
					return
				}
			}
			in[i] = reflect.ValueOf(req)
		}
		// 调用函数
		result := f.Func.Call(in)
		if f.hasError {
			if result[0].IsNil() {
				return
			}
			mlog.Error("调用%s.%s报错：%v", head.ActorName, head.FuncName, result[0].Interface())
		}
	}
}

func (f *FuncInfo) handleRpcGob(rval reflect.Value, head *pb.Head, buf []byte) func() {
	return func() {
		pos := 1
		if f.hasHead {
			pos++
		}
		// 解析参数参数
		in, err := encode.Decode(buf, f.Method, pos)
		if err != nil {
			mlog.Error("%s.%s参数解析错误: %v", head.ActorName, head.FuncName, err)
			return
		}
		// 设置 this
		in[0] = rval
		if f.hasHead {
			in[1] = reflect.ValueOf(head)
		}
		// 调用函数
		var result []reflect.Value
		if f.Type.IsVariadic() {
			result = f.Func.CallSlice(in)
		} else {
			result = f.Func.Call(in)
		}
		// 处理返回值
		if f.hasError {
			if result[0].IsNil() {
				return
			}
			mlog.Error("调用%s.%s报错：%v", head.ActorName, head.FuncName, result[0].Interface())
		}
	}
}
