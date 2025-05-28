package actor

import (
	"reflect"
	"universal/common/pb"
	"universal/library/encode"
	"universal/library/mlog"
	"universal/library/util"

	"github.com/golang/protobuf/proto"
)

var (
	bytesType = reflect.TypeOf((*[]byte)(nil)).Elem()
	headType  = reflect.TypeOf((*pb.Head)(nil))
	protoType = reflect.TypeOf((*proto.Message)(nil)).Elem()
	pool      = util.PoolSlice[reflect.Value](6)
)

func get(size int) []reflect.Value {
	rets := pool.Get().([]reflect.Value)
	return rets[:size]
}

func put(v []reflect.Value) {
	pool.Put(v)
}

type FuncInfo struct {
	reflect.Method
	hasHead bool
	isProto bool
	isBytes bool
}

func parseFuncInfo(m reflect.Method) *FuncInfo {
	ins, outs := m.Type.NumIn(), m.Type.NumOut()
	if outs > 1 || ins > 5 {
		return nil
	}
	hasHead, isProto, isBytes := false, true, true
	for i := 1; i < ins; i++ {
		if i == 1 {
			hasHead = m.Type.In(i).AssignableTo(headType)
		}
		if !m.Type.In(i).Implements(protoType) {
			isProto = false
		}
		if !m.Type.In(i).AssignableTo(bytesType) {
			isBytes = false
		}
	}
	return &FuncInfo{Method: m, hasHead: hasHead, isProto: isProto, isBytes: isBytes}
}

func (f *FuncInfo) handle(rval reflect.Value, head *pb.Head, args ...interface{}) func() {
	return func() {
			mlog.Debugf("Actor(%s.%s) head:%v, args:%v", head.ActorName, head.FuncName, head, args)
		size := f.Type.NumIn()
		in := get(size)
		defer put(in)
		in[0] = rval
		pos := 1
		if f.hasHead {
			in[1] = reflect.ValueOf(head)
			pos++
		}
		for i := pos; i < size-1; i++ {
			in[i] = reflect.ValueOf(args[i-pos])
		}
		// 可变参数
		if size > pos {
			if !f.Type.IsVariadic() {
				in[size-1] = reflect.ValueOf(args[size-pos-1])
			} else {
				if args2 := args[size-pos-1:]; len(args2) > 0 {
					arr := make([]reflect.Value, len(args2))
					for i, item := range args2 {
						arr[i] = reflect.ValueOf(item)
					}
					in[size-1] = reflect.ValueOf(arr)
				}
			}
		}
		// 调用函数
		var result []reflect.Value
		if !f.Type.IsVariadic() {
			result = f.Func.Call(in)
		} else {
			result = f.Func.CallSlice(in)
		}
		// 处理返回值
		response(f.Method, head, result)
	}
}

func (f *FuncInfo) handleRpc(rval reflect.Value, head *pb.Head, buf []byte) func() {
	return func() {
		size := f.Type.NumIn()
		in := get(size)
		defer put(in)
		in[0] = rval
		pos := 1
		if f.hasHead {
			in[1] = reflect.ValueOf(head)
			pos++
		}
		for i := pos; i < size; i++ {
			in[i] = reflect.New(f.Type.In(i).Elem())
		}
		if err := proto.Unmarshal(buf, in[pos].Interface().(proto.Message)); err != nil {
			mlog.Errorf("调用%s.%s报错：参数解析失败: head:%v, error:%v", head.ActorName, head.FuncName, head, err)
			return
		}
		// 调用函数
		result := f.Func.Call(in)
		// 处理返回值
		response(f.Method, head, result)
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
		// 处理返回值
		response(f.Method, head, result)
	}
}

func response(m reflect.Method, head *pb.Head, result []reflect.Value) {
	if m.Type.NumOut() <= 0 {
			mlog.Debugf("Actor(%s.%s) head:%v", head.ActorName, head.FuncName, head)
		return
	}

	switch vv := result[0].Interface().(type) {
	case error:
		mlog.Errorf("Actor(%s.%s) head:%v, error:%v", head.ActorName, head.FuncName, head, vv)
	case proto.Message:
		mlog.Debugf("Actor(%s.%s) head:%v, rsp:%v", head.ActorName, head.FuncName, head, vv)
		if len(head.Reply) > 0 {
			if err := responseFunc(head, vv); err != nil {
				mlog.Errorf("Reponse head:%v, error:%v", head, err)
			}
		} else {
			if err := sendFunc(head, vv); err != nil {
				mlog.Errorf("SendToClient head:%v, error:%v", head, err)
			}
		}
	case nil:
		mlog.Debugf("Actor(%s.%s) head:%v", head.ActorName, head.FuncName, head)
	default:
		mlog.Errorf("Actor(%s.%s) head:%v, unknown:%v", head.ActorName, head.FuncName, head, vv)
	}
}
