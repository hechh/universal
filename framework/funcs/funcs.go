package funcs

import (
	"reflect"
	"sync/atomic"
	"universal/common/pb"
	"universal/framework/domain"
	"universal/library/encode"
	"universal/library/mlog"
	"universal/library/util"

	"google.golang.org/protobuf/proto"
)

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
	hasHead := util.Or[int](ins > 1 && m.Type.In(1).AssignableTo(headType), 1, 0)
	hasReq := util.Or[int](ins > 2 && m.Type.In(2).Implements(reqType), 1, 0)
	hasRsp := util.Or[int](ins > 3 && m.Type.In(3).Implements(rspType), 1, 0)
	hasBytes := 1
	for i := util.Or[int](hasHead > 0, 1, 2); i < ins; i++ {
		if !m.Type.In(i).AssignableTo(bytesType) {
			hasBytes = 0
		}
	}
	return &Method{Method: m, ins: ins, flag: uint32(hasHead | hasReq<<1 | hasRsp<<2 | hasBytes<<3)}
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
		case CMD_HANDLER, NOTIFY_HANDLER:
			for i := pos; i < m.ins; i++ {
				params[i] = reflect.New(m.Type.In(i).Elem())
			}
			if err := proto.Unmarshal(buf, params[pos].Interface().(proto.Message)); err != nil {
				mlog.Errorf("参数解析失败 %v", head)
				return
			}
		case HEAD_BYTES_HANDLER, BYTES_HANDLER:
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
	case CMD_HANDLER:
		req := params[2].Interface().(proto.Message)
		rsp := params[3].Interface().(domain.IRspProto)
		rsp.SetHead(ToRspHead(err))
		var reterr error
		if atomic.CompareAndSwapUint32(&head.Reference, ref, ref) {
			reterr = sendRsp(head, rsp)
		}
		if err != nil || reterr != nil {
			mlog.Error(1, "head:%v, req:%v, rsp:%v, error:%v", head, req, rsp, reterr)
		} else {
			mlog.Debug(1, "head:%v, req:%v, rsp:%v, error:%v", head, req, rsp, reterr)
		}
	case NOTIFY_HANDLER:
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
