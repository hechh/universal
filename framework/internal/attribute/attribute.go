package attribute

import (
	"sync/atomic"
	"time"
	"universal/common/base"
	"universal/common/pb"
	"universal/framework/define"
	"universal/library/mlog"
	"universal/library/templ"
	"universal/library/uerror"

	"github.com/golang/protobuf/proto"
)

type Attribute struct {
	fun     interface{}
	reqname string
	rspname string
}

func NewAttribute(f interface{}, args ...string) *Attribute {
	return &Attribute{
		fun:     f,
		reqname: templ.Index[string](args, 0, ""),
		rspname: templ.Index[string](args, 1, ""),
	}
}

func (a *Attribute) add(head *pb.Head) uint32 {
	switch a.fun.(type) {
	case define.TwoFunc, define.HeadTwoFunc:
		return atomic.AddUint32(&head.Reference, 1)
	default:
		return atomic.LoadUint32(&head.Reference)
	}
}

func (a *Attribute) Call(ff define.SendRspFunc, head *pb.Head, args ...proto.Message) func() {
	ref := a.add(head)
	return func() {
		nowMs := time.Now().UnixMilli()
		var err error
		switch f := a.fun.(type) {
		case define.ZeroFunc:
			err = f()
		case define.HeadZeroFunc:
			err = f(head)
		case define.OneFunc:
			err = f(args[0])
		case define.HeadOneFunc:
			err = f(head, args[0])
		case define.TwoFunc:
			err = f(args[0], args[1])
		case define.HeadTwoFunc:
			err = f(head, args[0], args[1])
		default:
			err = uerror.New(1, -1, "处理类型不支持: %s.%s", head.ActorName, head.FuncName)
		}
		a.result(ref, nowMs, err, ff, head, args...)
	}
}

func (a *Attribute) Rpc(pp define.IFactory, ff define.SendRspFunc, head *pb.Head, buf []byte) func() {
	ref := a.add(head)
	return func() {
		nowMs := time.Now().UnixMilli()
		var err error
		var req, rsp proto.Message
		switch f := a.fun.(type) {
		case define.ZeroFunc:
			err = f()
		case define.HeadZeroFunc:
			err = f(head)
		case define.OneFunc:
			req = pp.New(a.reqname)
			if err = proto.Unmarshal(buf, req); err == nil {
				err = f(req)
			}
		case define.HeadOneFunc:
			req = pp.New(a.reqname)
			if err = proto.Unmarshal(buf, req); err == nil {
				err = f(head, req)
			}
		case define.TwoFunc:
			req = pp.New(a.reqname)
			if err = proto.Unmarshal(buf, req); err == nil {
				rsp = pp.New(a.rspname)
				err = f(req, rsp)
			}
		case define.HeadTwoFunc:
			req = pp.New(a.reqname)
			if err = proto.Unmarshal(buf, req); err == nil {
				rsp = pp.New(a.rspname)
				err = f(head, req, rsp)
			}
		default:
			err = uerror.New(1, -1, "处理类型不支持: %s.%s", head.ActorName, head.FuncName)
		}
		a.result(ref, nowMs, err, ff, head, req, rsp)
	}
}

func (a *Attribute) result(ref uint32, nowMs int64, err error, ff define.SendRspFunc, head *pb.Head, args ...proto.Message) {
	endMs := time.Now().UnixMilli()
	var reterr error
	var rsp define.IRspProto
	switch a.fun.(type) {
	case define.TwoFunc, define.HeadTwoFunc:
		if atomic.CompareAndSwapUint32(&head.Reference, ref, ref) {
			rsp = args[1].(define.IRspProto)
			rsp.SetHead(base.ToRspHead(err))
			if ff != nil {
				reterr = ff(head, rsp)
			}
		}
	}
	if err != nil || reterr != nil {
		mlog.Errorf("耗时(%dms)|Req<%v>|Rsp<%v>|SendError<%v>|Error<%v>", endMs-nowMs, args[0], rsp, reterr, err)
	} else {
		mlog.Tracef("耗时(%dms)|Req<%v>|Rsp<%v>|SendError<%v>", endMs-nowMs, args[0], rsp, reterr)
	}
}
