package funcs

import (
	"reflect"
	"sync/atomic"
	"universal/common/pb"
	"universal/framework/domain"
	"universal/library/encode"
	"universal/library/mlog"
	"universal/library/uerror"
	"universal/library/util"

	"github.com/golang/protobuf/proto"
)

const (
	HEAD_FLAG      = 1 << 0
	REQ_FLAG       = 1 << 1
	RSP_FLAG       = 1 << 2
	BYTES_FLAG     = 1 << 3
	INTERFACE_FLAG = 1 << 4
	GOB_FLAG       = 1 << 5
)

var (
	headType      = reflect.TypeOf((*pb.Head)(nil))
	reqType       = reflect.TypeOf((*proto.Message)(nil)).Elem()
	rspType       = reflect.TypeOf((*domain.IRspProto)(nil)).Elem()
	bytesType     = reflect.TypeOf((*[]byte)(nil)).Elem()
	errorType     = reflect.TypeOf((*error)(nil)).Elem()
	interfaceType = reflect.TypeOf((*interface{})(nil)).Elem()
	nilValue      = reflect.ValueOf((*error)(nil))
	args          = util.ArrayPool[reflect.Value](6)
	sendRsp       func(*pb.Head, proto.Message) error
)

func Init(f func(*pb.Head, proto.Message) error) {
	sendRsp = f
}

func get(size int) []reflect.Value {
	rets := args.Get().([]reflect.Value)
	return rets[:size]
}

func put(rets []reflect.Value) {
	args.Put(rets)
}

func toRspHead(err error) *pb.RspHead {
	switch vv := err.(type) {
	case nil:
		return nil
	case *uerror.UError:
		return &pb.RspHead{Code: vv.GetCode(), ErrMsg: vv.GetMsg()}
	case error:
		return &pb.RspHead{Code: -1, ErrMsg: vv.Error()}
	default:
		return nil
	}
}

type Method struct {
	reflect.Method
	ins  int
	flag uint32
}

func NewMethod(m reflect.Method) *Method {
	ins, outs := m.Type.NumIn(), m.Type.NumOut()
	if outs > 1 || outs == 1 && !m.Type.Out(0).Implements(errorType) {
		return nil
	}
	flag := uint32(0)
	for i := 1; i < ins; i++ {
		if m.Type.In(i).AssignableTo(headType) {
			flag = flag | HEAD_FLAG
		} else if m.Type.In(i).Implements(rspType) {
			flag = flag | RSP_FLAG
		} else if m.Type.In(i).Implements(reqType) {
			flag = flag | REQ_FLAG
		} else if m.Type.In(i).AssignableTo(bytesType) {
			flag = flag | BYTES_FLAG
		} else if m.Type.In(i) == interfaceType {
			flag = flag | INTERFACE_FLAG
		} else {
			flag = flag | GOB_FLAG
		}
	}
	return &Method{Method: m, ins: ins, flag: flag}
}

func (m *Method) Call(rval reflect.Value, head *pb.Head, args ...interface{}) func() {
	ref := atomic.AddUint32(&head.Reference, 1)
	return func() {
		params := get(m.ins)
		defer put(params)
		params[0] = rval
		pos := 1
		if m.flag&HEAD_FLAG == HEAD_FLAG {
			params[1] = reflect.ValueOf(head)
			pos++
		}
		for i := pos; i < m.ins; i++ {
			params[i] = reflect.ValueOf(args[i-pos])
		}
		// 调用函数
		rets := m.Func.Call(params)
		// 返回值
		m.result(ref, head, params, util.Index[reflect.Value](rets, 0, nilValue).Interface().(error))
	}
}

func (m *Method) Rpc(rval reflect.Value, head *pb.Head, buf []byte) func() {
	ref := atomic.AddUint32(&head.Reference, 1)
	return func() {
		params := get(m.ins)
		defer put(params)
		params[0] = rval
		pos := 1
		if m.flag&HEAD_FLAG == HEAD_FLAG {
			params[1] = reflect.ValueOf(head)
			pos++
		}
		switch m.flag {
		case (HEAD_FLAG | REQ_FLAG | RSP_FLAG), (HEAD_FLAG | REQ_FLAG), (HEAD_FLAG | RSP_FLAG), (REQ_FLAG | RSP_FLAG), REQ_FLAG:
			for i := pos; i < m.ins; i++ {
				params[i] = reflect.New(m.Type.In(i).Elem())
			}
			if err := proto.Unmarshal(buf, params[pos].Interface().(proto.Message)); err != nil {
				mlog.Errorf("参数解析失败 %v", head)
				return
			}
		case (HEAD_FLAG | BYTES_FLAG), (BYTES_FLAG):
			params[pos] = reflect.ValueOf(buf)
		default:
			if err := encode.Decode(buf, m.Method, params, pos); err != nil {
				mlog.Errorf("参数解析失败 %v", head)
				return
			}
		}
		// 调用函数
		rets := m.Func.Call(params)
		// 返回值
		m.result(ref, head, params, util.Index[reflect.Value](rets, 0, nilValue).Interface().(error))
	}
}

func (m *Method) result(ref uint32, head *pb.Head, params []reflect.Value, err error) {
	switch m.flag {
	case (HEAD_FLAG | REQ_FLAG | RSP_FLAG):
		req := params[2].Interface().(proto.Message)
		rsp := params[3].Interface().(domain.IRspProto)
		rsp.SetHead(toRspHead(err))
		var reterr error
		if atomic.CompareAndSwapUint32(&head.Reference, ref, ref) {
			reterr = sendRsp(head, rsp)
		}
		if err != nil || reterr != nil {
			mlog.Error(1, "head:%v, req:%v, rsp:%v, error:%v", head, req, rsp, reterr)
		} else {
			mlog.Debug(1, "head:%v, req:%v, rsp:%v, error:%v", head, req, rsp, reterr)
		}
	case (HEAD_FLAG | REQ_FLAG):
		req := params[2].Interface().(proto.Message)
		if err != nil {
			mlog.Error(1, "head:%v, notify:%v, error:%v", head, req, err)
		} else {
			mlog.Debug(1, "head:%v, notify:%v, error:%v", head, req, err)
		}
	default:
		if err != nil {
			mlog.Error(1, "head:%v, error:%v", head, err)
		} else {
			mlog.Debug(1, "head:%v, error:%v", head, err)
		}
	}
}
