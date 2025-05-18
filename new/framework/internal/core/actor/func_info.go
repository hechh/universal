package actor

import (
	"poker_server/common/pb"
	"poker_server/framework/library/encode"
	"poker_server/framework/library/mlog"
	"reflect"

	"github.com/golang/protobuf/proto"
)

var (
	headType  = reflect.TypeOf((*pb.Head)(nil))
	errorType = reflect.TypeOf((*error)(nil)).Elem()
	protoType = reflect.TypeOf((*proto.Message)(nil)).Elem()
)

type FuncInfo struct {
	reflect.Method
	hasHead  bool
	hasError bool
	isNotify bool
	isCmd    bool
}

func parseFunc(m reflect.Method) *FuncInfo {
	hasError := m.Type.NumOut() > 0 && m.Type.Out(0).Implements(errorType)
	hasHead := m.Type.NumIn() >= 2 && m.Type.In(1).AssignableTo(headType)
	isNotify := m.Type.NumIn() == 3 && hasHead && m.Type.In(2).Implements(protoType)
	isCmd := m.Type.NumIn() == 4 && hasHead && m.Type.In(2).Implements(protoType) && m.Type.In(3).Implements(protoType)
	return &FuncInfo{
		Method:   m,
		hasHead:  hasHead,
		hasError: hasError,
		isNotify: isNotify,
		isCmd:    isCmd,
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
			mlog.Errorf("调用%s.%s报错：head:%v, error:%v", head.ActorName, head.FuncName, head, result[0].Interface())
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
			mlog.Errorf("调用%s.%s报错：head:%v, error:%v", head.ActorName, head.FuncName, head, result[0].Interface())
		}
	}
}

// 远程调用(不可调用可变参数)
func (f *FuncInfo) handleRpcNotify(rval reflect.Value, head *pb.Head, buf []byte) func() {
	return func() {
		in := make([]reflect.Value, f.Type.NumIn())
		in[0] = rval
		in[1] = reflect.ValueOf(head)
		in[2] = reflect.New(f.Type.In(2).Elem())
		if err := proto.Unmarshal(buf, in[2].Interface().(proto.Message)); err != nil {
			mlog.Errorf("%s.%s参数解析报错: head:%v, error:%v", head.ActorName, head.FuncName, head, err)
			return
		}
		// 调用函数
		result := f.Func.Call(in)
		if f.hasError {
			if result[0].IsNil() {
				return
			}
			mlog.Errorf("调用%s.%s报错：head:%v, error:%v", head.ActorName, head.FuncName, head, result[0].Interface())
		}
	}
}

// 远程调用(不可调用可变参数)
func (f *FuncInfo) handleRpcCmd(rval reflect.Value, head *pb.Head, buf []byte) func() {
	return func() {
		in := make([]reflect.Value, f.Type.NumIn())
		in[0] = rval
		in[1] = reflect.ValueOf(head)
		in[2] = reflect.New(f.Type.In(2).Elem())
		in[3] = reflect.New(f.Type.In(3).Elem())
		if err := proto.Unmarshal(buf, in[2].Interface().(proto.Message)); err != nil {
			mlog.Errorf("%s.%s参数解析报错: head:%v, error:%v", head.ActorName, head.FuncName, head, err)
			return
		}
		// 调用函数
		result := f.Func.Call(in)
		if f.hasError {
			if result[0].IsNil() {
				return
			}
			mlog.Errorf("调用%s.%s报错：head:%v, error:%v", head.ActorName, head.FuncName, head, result[0].Interface())
		}
		// 返回应答
		if len(head.Reply) > 0 && busObj != nil {
			if err := busObj.Response(head, in[3].Interface().(proto.Message)); err != nil {
				mlog.Errorf("调用%s.%s应答报错：head:%v, error:%v", head.ActorName, head.FuncName, head, err)
			}
		}
	}
}

// 远程调用(不可调用可变参数)
func (f *FuncInfo) handleRpcGob(rval reflect.Value, head *pb.Head, buf []byte) func() {
	return func() {
		pos := 1
		if f.hasHead {
			pos++
		}
		// 解析参数参数
		in, err := encode.Decode(buf, f.Method, pos)
		if err != nil {
			mlog.Errorf("%s.%s参数解析错误: head:%v, error:%v", head.ActorName, head.FuncName, head, err)
			return
		}
		// 设置 this
		in[0] = rval
		if f.hasHead {
			in[1] = reflect.ValueOf(head)
		}
		// 调用函数
		result := f.Func.Call(in)
		if f.hasError {
			if result[0].IsNil() {
				return
			}
			mlog.Errorf("调用%s.%s报错：head:%v, error:%v", head.ActorName, head.FuncName, head, result[0].Interface())
		}
	}
}
