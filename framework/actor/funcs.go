package actor

import (
	"reflect"
	"universal/common/pb"
	"universal/framework/domain"
	"universal/library/encode"
	"universal/library/mlog"
	"universal/library/uerror"
	"universal/library/util"

	"github.com/golang/protobuf/proto"
)

const (
	HEAD_FLAG   = 1 << 4
	PROTO_FLAG  = 1 << 3
	RSP_FLAG    = 1 << 2
	NOTIFY_FLAG = 1 << 1
	BYTES_FLAG  = 1 << 0
	CMD_FLAG    = HEAD_FLAG | PROTO_FLAG | RSP_FLAG
)

var (
	headType  = reflect.TypeOf((*pb.Head)(nil))
	protoType = reflect.TypeOf((*proto.Message)(nil)).Elem()
	rspType   = reflect.TypeOf((*domain.IRspProto)(nil)).Elem()
	bytesType = reflect.TypeOf((*[]byte)(nil)).Elem()
	errorType = reflect.TypeOf((*error)(nil)).Elem()
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
	numin int
	flag  uint8
}

func parseFuncInfo(m reflect.Method) *FuncInfo {
	ins, outs := m.Type.NumIn(), m.Type.NumOut()
	if outs > 1 {
		return nil
	}
	hasHead := util.Ifelse[int](ins > 1 && m.Type.In(1).AssignableTo(headType), 1, 0)
	isBytes, isProto := 1, 1
	for i := util.Ifelse[int](hasHead < 1, 1, 2); i < ins; i++ {
		if !m.Type.In(i).AssignableTo(bytesType) {
			isBytes = 0
		}
		if !m.Type.In(i).Implements(protoType) {
			isProto = 0
		}
	}
	flag := hasHead > 0 && isProto > 0 && ins > 3
	isRsp := util.Ifelse[int](flag, 1, 0)
	for i := util.Ifelse[int](flag, 2, ins); i < ins; i++ {
		if !m.Type.In(i).Implements(rspType) {
			isRsp = 0
		}
	}
	isNotify := util.Ifelse[int](hasHead > 0 && isProto > 0 && ins == 3, 1, 0)
	return &FuncInfo{
		Method: m,
		numin:  ins,
		flag:   uint8(hasHead<<4 | isProto<<3 | isRsp<<2 | isNotify<<1 | isBytes),
	}
}

// 本地直接调用notify函数
func (f *FuncInfo) localProto(rval reflect.Value, head *pb.Head, args ...interface{}) func() {
	return func() {
		in := get(f.numin)
		defer put(in)
		in[0] = rval
		pos := 1
		if f.flag&HEAD_FLAG > 0 {
			in[1] = reflect.ValueOf(head)
			pos++
		}
		for i := pos; i < f.numin; i++ {
			in[i] = reflect.ValueOf(args[i-pos])
		}

		// 调用函数
		rets := f.Func.Call(in)

		// 处理返回值
		result(head, rets, args)
	}
}

// 本地直接调用cmd函数
func (f *FuncInfo) localCmd(rval reflect.Value, head *pb.Head, args ...interface{}) func() {
	return func() {
		in := get(f.numin)
		defer put(in)
		in[0] = rval
		in[1] = reflect.ValueOf(head)
		for i := 2; i < f.numin; i++ {
			in[i] = reflect.ValueOf(args[i-2])
		}

		// 调用函数
		aname, fname := head.ActorName, head.FuncName
		rets := f.Func.Call(in)

		// 处理返回值
		if aname == head.ActorName && fname == head.FuncName && node.Type != pb.NodeType_NodeTypeGate {
			response(head, rets, in[3:]...)
		} else {
			result(head, rets, args)
		}
	}
}

// rpc跨服务调用notify函数
func (f *FuncInfo) rpcNotify(rval reflect.Value, head *pb.Head, buf []byte) func() {
	return func() {
		in := get(f.numin)
		defer put(in)
		in[0] = rval
		in[1] = reflect.ValueOf(head)
		in[2] = reflect.New(f.Type.In(2).Elem())
		if err := proto.Unmarshal(buf, in[2].Interface().(proto.Message)); err != nil {
			mlog.Error(1, "参数解析失败: head:%v, error:%v", head, err)
			return
		}
		// 调用函数
		rets := f.Func.Call(in)

		// 处理返回值
		result(head, rets, buf)
	}
}

// rpc跨服务调用cmd函数
func (f *FuncInfo) rpcCmd(rval reflect.Value, head *pb.Head, buf []byte) func() {
	return func() {
		in := get(f.numin)
		defer put(in)
		in[0] = rval
		in[1] = reflect.ValueOf(head)
		for i := 2; i < f.numin; i++ {
			in[i] = reflect.New(f.Type.In(i).Elem())
		}
		if err := proto.Unmarshal(buf, in[2].Interface().(proto.Message)); err != nil {
			mlog.Error(1, "参数解析失败: head:%v, error:%v", head, err)
			return
		}

		// 调用函数
		aname, fname := head.ActorName, head.FuncName
		rets := f.Func.Call(in)

		// 处理返回值
		if aname == head.ActorName && fname == head.FuncName && node.Type != pb.NodeType_NodeTypeGate {
			response(head, rets, in[3:]...)
		} else {
			result(head, rets, buf)
		}
	}
}

// rpc跨服务调用gob函数
func (f *FuncInfo) rpcGob(rval reflect.Value, head *pb.Head, buf []byte) func() {
	return func() {
		pos := 1
		if f.flag&HEAD_FLAG > 0 {
			pos++
		}
		// 解析参数参数
		in, err := encode.Decode(buf, f.Method, pos)
		if err != nil {
			mlog.Error(1, "参数解析失败: head:%v, error:%v", head, err)
			return
		}
		in[0] = rval
		if f.flag&HEAD_FLAG > 0 {
			in[1] = reflect.ValueOf(head)
		}

		// 调用函数
		rets := f.Func.Call(in)

		// 处理返回值
		result(head, rets, buf)
	}
}

func result(head *pb.Head, rets []reflect.Value, buf interface{}) {
	if len(rets) > 0 && !rets[0].IsNil() {
		mlog.Error(1, "head:%v, error:%v, args:%v", head, rets[0].Interface(), buf)
	} else {
		mlog.Debug(1, "head:%v, args:%v", head, buf)
	}
}

func response(head *pb.Head, results []reflect.Value, rsps ...reflect.Value) {
	var err error
	if len(results) > 0 && !results[0].IsNil() {
		err = results[0].Interface().(error)
	}
	rsphead := reflect.ValueOf(toRspHead(err))

	// 处理返回值
	for _, rsp := range rsps {
		if !rsphead.IsNil() {
			if val := rsp.Field(3); val.CanSet() {
				val.Set(rsphead)
			}
		}
		if err != nil {
			mlog.Error(1, "head:%v, error:%v, rsp:%v", head, err, rsp.Interface())
		} else {
			mlog.Debug(1, "head:%v, rsp:%v", head, rsp.Interface())
		}
		sendRspFunc(head, rsp.Interface().(domain.IRspProto))
	}
}

func toRspHead(err error) *pb.RspHead {
	switch vv := err.(type) {
	case *uerror.UError:
		return &pb.RspHead{Code: int32(vv.GetCode()), Msg: vv.GetMsg()}
	case nil:
		return nil
	}
	return &pb.RspHead{Code: int32(pb.ErrorCode_Unknown), Msg: err.Error()}
}
