package method

import (
	"reflect"
	"sync/atomic"
	"universal/common/pb"
	"universal/framework/define"
	"universal/library/encode"
	"universal/library/mlog"
	"universal/library/uerror"
	"universal/library/util"

	"github.com/golang/protobuf/proto"
)

var (
	args         = util.ArrayPool[reflect.Value](6)
	sendResponse func(*pb.Head, proto.Message) error
)

func Init(f func(*pb.Head, proto.Message) error) {
	sendResponse = f
}

func get(size int) []reflect.Value {
	rets := args.Get().([]reflect.Value)
	return rets[:size]
}

func put(rets []reflect.Value) {
	args.Put(rets)
}

type Method struct {
	define.IActor
	reflect.Method
	mask uint32
}

func NewMethod(act define.IActor, m reflect.Method) *Method {
	ins, outs := m.Type.NumIn(), m.Type.NumOut()
	if outs > 1 || outs == 1 && !m.Type.Out(0).Implements(errorType) {
		return nil
	}
	mask := uint32(0)
	for i := 1; i < ins; i++ {
		if m.Type.In(i).AssignableTo(headType) {
			mask = mask | HEAD_FLAG
			continue
		}
		if m.Type.In(i).Implements(rspType) {
			mask = mask | RSP_FLAG
			continue
		}
		if m.Type.In(i).Implements(reqType) {
			mask = mask | REQ_FLAG
			continue
		}
		if m.Type.In(i).AssignableTo(bytesType) {
			mask = mask | BYTES_FLAG
			continue
		}
		if m.Type.In(i) == interfaceType {
			mask = mask | INTERFACE_FLAG
		} else {
			mask = mask | GOB_FLAG
		}
	}
	return &Method{Method: m, mask: mask, IActor: act}
}

func (r *Method) GetFuncName() string {
	return r.Name
}

func (r *Method) addReference(head *pb.Head) uint32 {
	if r.mask == HEAD_REQ_RSP_MASK {
		return atomic.AddUint32(&head.Reference, 1)
	}
	return 0
}

func (r *Method) Call(rval reflect.Value, head *pb.Head, args ...interface{}) func() {
	ref := r.addReference(head)
	return func() {
		params := get(r.Type.NumIn())
		defer put(params)

		params[0] = rval
		pos := 1
		if r.mask&HEAD_FLAG == HEAD_FLAG {
			params[pos] = reflect.ValueOf(head)
			pos++
		}
		for i := pos; i < r.Type.NumIn(); i++ {
			params[i] = reflect.ValueOf(args[i-pos])
		}

		rets := r.Func.Call(params)
		r.result(pos, ref, head, params, util.Index[reflect.Value](rets, 0, nilError).Interface().(error))
	}
}

func (r *Method) Rpc(rval reflect.Value, head *pb.Head, buf []byte) func() {
	ref := r.addReference(head)
	return func() {
		params := get(r.Type.NumIn())
		defer put(params)

		params[0] = rval
		pos := 1
		if r.mask&HEAD_FLAG == HEAD_FLAG {
			params[pos] = reflect.ValueOf(head)
			pos++
		}
		switch r.mask {
		case HEAD_REQ_RSP_MASK, HEAD_REQ_MASK, HEAD_RSP_MASK, REQ_RSP_MASK, REQ_FLAG, RSP_FLAG:
			for i := pos; i < r.Type.NumIn(); i++ {
				params[i] = reflect.New(r.Type.In(i).Elem())
			}
			if err := proto.Unmarshal(buf, params[pos].Interface().(proto.Message)); err != nil {
				mlog.Errorf("参数解析失败 %v", err)
				return
			}
		case HEAD_BYTES_MASK, BYTES_FLAG, HEAD_INTERFACE_MASK, INTERFACE_FLAG:
			params[pos] = reflect.ValueOf(buf)
		default:
			if err := encode.Decode(buf, r.Method, params, pos); err != nil {
				mlog.Errorf("参数解析失败 %v", err)
				return
			}
		}

		rets := r.Func.Call(params)
		r.result(pos, ref, head, params, util.Index[reflect.Value](rets, 0, nilError).Interface().(error))
	}
}

func (r *Method) result(pos int, ref uint32, head *pb.Head, params []reflect.Value, err error) {
	switch r.mask {
	case HEAD_REQ_RSP_MASK, REQ_RSP_MASK:
		req := params[pos].Interface().(proto.Message)
		rsp := params[pos+1].Interface().(define.IRspProto)
		rsp.SetHead(toRspHead(err))
		var reterr error
		if r.mask == HEAD_REQ_RSP_MASK && atomic.CompareAndSwapUint32(&head.Reference, ref, ref) {
			reterr = sendResponse(head, rsp)
		}
		if err != nil || reterr != nil {
			mlog.Error(1, "%d|%s|%s|%d|%v|%v req:%v, rsp:%v, error:%v", head.Uid, head.ActorName, head.FuncName, head.ActorId, head.Src, head.Dst, req, rsp, reterr)
		} else {
			mlog.Debug(1, "%d|%s|%s|%d|%v|%v req:%v, rsp:%v, error:%v", head.Uid, head.ActorName, head.FuncName, head.ActorId, head.Src, head.Dst, req, rsp, reterr)
		}
	case HEAD_REQ_MASK, REQ_FLAG, RSP_FLAG:
		req := params[pos].Interface().(proto.Message)
		if err != nil {
			mlog.Error(1, "%d|%s|%s|%d|%v|%v req:%v, error:%v", head.Uid, head.ActorName, head.FuncName, head.ActorId, head.Src, head.Dst, req, err)
		} else {
			mlog.Debug(1, "%d|%s|%s|%d|%v|%v req:%v, error:%v", head.Uid, head.ActorName, head.FuncName, head.ActorId, head.Src, head.Dst, req, err)
		}
	default:
		if err != nil {
			mlog.Error(1, "%d|%s|%s|%d|%v|%v error:%v", head.Uid, head.ActorName, head.FuncName, head.ActorId, head.Src, head.Dst, err)
		} else {
			mlog.Debug(1, "%d|%s|%s|%d|%v|%v error:%v", head.Uid, head.ActorName, head.FuncName, head.ActorId, head.Src, head.Dst, err)
		}
	}
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
