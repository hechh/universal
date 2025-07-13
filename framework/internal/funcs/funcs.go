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
	sendRsp       func(*pb.Head, proto.Message) error
	args          = util.ArrayPool[reflect.Value](6)
	headType      = reflect.TypeOf((*pb.Head)(nil))
	reqType       = reflect.TypeOf((*proto.Message)(nil)).Elem()
	rspType       = reflect.TypeOf((*domain.IRspProto)(nil)).Elem()
	bytesType     = reflect.TypeOf((*[]byte)(nil)).Elem()
	errorType     = reflect.TypeOf((*error)(nil)).Elem()
	interfaceType = reflect.TypeOf((*interface{})(nil)).Elem()
	nilError      = reflect.ValueOf((*error)(nil)).Elem()
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

type Method struct {
	reflect.Method
	mask uint32
	act  domain.IActor
}

func NewMethod(act domain.IActor, m reflect.Method) *Method {
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
	return &Method{Method: m, mask: mask, act: act}
}

func (r *Method) Parse(head *pb.Head) error {
	actorType := head.Dst.ActorId & 0xFF
	if r.act.GetActorType() != uint32(actorType) {
		return uerror.New(1, -1, "ActorType错误%v", head.Dst)
	}
	head.ActorName = r.act.GetActorName()
	head.FuncName = r.Name
	head.ActorId = (head.Dst.ActorId >> 8)
	return nil
}

func (r *Method) Call(rval reflect.Value, head *pb.Head, args ...interface{}) func() {
	ref := r.getReference(head)
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
	ref := r.getReference(head)
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
		case (HEAD_FLAG | REQ_FLAG | RSP_FLAG), (HEAD_FLAG | REQ_FLAG), (HEAD_FLAG | RSP_FLAG), (REQ_FLAG | RSP_FLAG), REQ_FLAG, RSP_FLAG:
			for i := pos; i < r.Type.NumIn(); i++ {
				params[i] = reflect.New(r.Type.In(i).Elem())
			}
			if err := proto.Unmarshal(buf, params[pos].Interface().(proto.Message)); err != nil {
				mlog.Errorf("参数解析失败 %v", head)
				return
			}
		case (HEAD_FLAG | BYTES_FLAG), BYTES_FLAG, (HEAD_FLAG | INTERFACE_FLAG), INTERFACE_FLAG:
			params[pos] = reflect.ValueOf(buf)
		default:
			if err := encode.Decode(buf, r.Method, params, pos); err != nil {
				mlog.Errorf("参数解析失败 %v", head)
				return
			}
		}

		rets := r.Func.Call(params)
		r.result(pos, ref, head, params, util.Index[reflect.Value](rets, 0, nilError).Interface().(error))
	}
}

func (r *Method) getReference(head *pb.Head) uint32 {
	if r.mask == (HEAD_FLAG | REQ_FLAG | RSP_FLAG) {
		return atomic.AddUint32(&head.Reference, 1)
	}
	return 0
}

func (r *Method) result(pos int, ref uint32, head *pb.Head, params []reflect.Value, err error) {
	switch r.mask {
	case (HEAD_FLAG | REQ_FLAG | RSP_FLAG), (REQ_FLAG | RSP_FLAG):
		req := params[pos].Interface().(proto.Message)
		rsp := params[pos+1].Interface().(domain.IRspProto)
		rsp.SetHead(toRspHead(err))
		var reterr error
		if r.mask == (HEAD_FLAG|REQ_FLAG|RSP_FLAG) && atomic.CompareAndSwapUint32(&head.Reference, ref, ref) {
			reterr = sendRsp(head, rsp)
		}
		if err != nil || reterr != nil {
			mlog.Error(1, "%d|%s|%s|%d|%v|%v req:%v, rsp:%v, error:%v", head.Uid, head.ActorName, head.FuncName, head.ActorId, head.Src, head.Dst, req, rsp, reterr)
		} else {
			mlog.Debug(1, "%d|%s|%s|%d|%v|%v req:%v, rsp:%v, error:%v", head.Uid, head.ActorName, head.FuncName, head.ActorId, head.Src, head.Dst, req, rsp, reterr)
		}
	case (HEAD_FLAG | REQ_FLAG), REQ_FLAG, RSP_FLAG:
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
