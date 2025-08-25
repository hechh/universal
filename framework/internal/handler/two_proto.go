package handler

import (
	"sync/atomic"
	"time"
	"universal/common/base"
	"universal/common/pb"
	"universal/framework/define"
	"universal/library/mlog"
)

type TwoProto[S any, T any, R any] func(*S, *pb.Head, *T, *R) error

func (f TwoProto[S, T, R]) NewReq() *T {
	return new(T)
}

func (f TwoProto[S, T, R]) NewRsp() *R {
	return new(R)
}

func (f TwoProto[S, T, R]) Call(sendrsp define.SendRspFunc, s interface{}, head *pb.Head, args ...interface{}) func() {
	ref := atomic.AddUint32(&head.Reference, 1)
	return func() {
		// 参数解析
		req, ok := any(args[0]).(*T)
		if !ok {
			mlog.Errorf("调用%s.%s参数类型错误%v", head.ActorName, head.FuncName, args[0])
			return
		}
		rsp, ok := any(args[1]).(*R)
		if !ok {
			mlog.Errorf("调用%s.%s参数类型错误%v", head.ActorName, head.FuncName, args[1])
			return
		}
		obj, ok := s.(*S)
		if !ok {
			mlog.Errorf("调用%s.%s参数类型错误%v", head.ActorName, head.FuncName, s)
			return
		}

		// 接口调用
		startMs := time.Now().UnixMilli()
		err := f(obj, head, req, rsp)
		endMs := time.Now().UnixMilli()
		if err != nil {
			mlog.Errorf("调用%s.%s耗时(%dms)|Req<%v>|Rsp<%v>|Error<%v>", head.ActorName, head.FuncName, endMs-startMs, req, rsp, err)
		} else {
			mlog.Tracef("调用%s.%s耗时(%dms)|Req<%v>|Rsp<%v>", head.ActorName, head.FuncName, endMs-startMs, req, rsp)
		}

		// 是否回复
		if sendrsp != nil && atomic.CompareAndSwapUint32(&head.Reference, ref, ref) {
			rspProto, ok := any(rsp).(define.IRspProto)
			if !ok {
				return
			}
			rspProto.SetHead(base.ToRspHead(err))
			if err = sendrsp(head, rspProto); err != nil {
				mlog.Errorf("回复失败|Rsp<%v>|%v", rsp, err)
			}
		}
	}
}

func (f TwoProto[S, T, R]) Rpc(sendrsp define.SendRspFunc, s interface{}, head *pb.Head, buf []byte) func() {
	ref := atomic.AddUint32(&head.Reference, 1)
	return func() {
		// 参数解析
		req := f.NewReq()
		if err := base.Unmarshal[T](buf, req); err != nil {
			mlog.Errorf("调用%s.%s参数解析失败%v", head.ActorName, head.FuncName, err)
			return
		}
		obj, ok := s.(*S)
		if !ok {
			mlog.Errorf("调用%s.%s参数类型错误%v", head.ActorName, head.FuncName, s)
			return
		}
		rsp := f.NewRsp()

		// 接口调用
		startMs := time.Now().UnixMilli()
		err := f(obj, head, req, rsp)
		endMs := time.Now().UnixMilli()
		if err != nil {
			mlog.Errorf("调用%s.%s耗时(%dms)|Req<%v>|Rsp<%v>|Error<%v>", head.ActorName, head.FuncName, endMs-startMs, req, rsp, err)
		} else {
			mlog.Tracef("调用%s.%s耗时(%dms)|Req<%v>|Rsp<%v>", head.ActorName, head.FuncName, endMs-startMs, req, rsp)
		}

		// 是否回复
		if sendrsp != nil && atomic.CompareAndSwapUint32(&head.Reference, ref, ref) {
			rspProto, ok := any(rsp).(define.IRspProto)
			if !ok {
				return
			}
			rspProto.SetHead(base.ToRspHead(err))
			if err = sendrsp(head, rspProto); err != nil {
				mlog.Errorf("回复失败|Rsp<%v>|%v", rsp, err)
			}
		}
	}
}
