package handler

import (
	"sync/atomic"
	"time"
	"universal/common/base"
	"universal/common/pb"
	"universal/framework/define"
	"universal/library/mlog"
)

type Handler[S any, T any, R any] func(*S, *pb.Head, *T, *R) error

func (f Handler[S, T, R]) NewReq() *T {
	return new(T)
}

func (f Handler[S, T, R]) NewRsp() *R {
	return new(R)
}

func (f Handler[S, T, R]) Call(sendrsp define.SendRspFunc, s interface{}, head *pb.Head, args ...interface{}) func() {
	ref := atomic.AddUint32(&head.Reference, 1)
	return func() {
		// 参数解析
		req, ok := any(args[0]).(*T)
		if !ok {
			mlog.Errorf("调用%s参数类型错误%v", val2str[head.ActorFunc], args[0])
			return
		}
		rsp, ok := any(args[1]).(*R)
		if !ok {
			mlog.Errorf("调用%s参数类型错误%v", val2str[head.ActorFunc], args[1])
			return
		}
		obj, ok := s.(*S)
		if !ok {
			mlog.Errorf("调用%s参数类型错误%v", val2str[head.ActorFunc], s)
			return
		}

		// 接口调用
		startMs := time.Now().UnixMilli()
		err := f(obj, head, req, rsp)
		endMs := time.Now().UnixMilli()
		if err != nil {
			mlog.Errorf("耗时(%dms)|Req<%v>|Rsp<%v>|Error<%v>", endMs-startMs, req, rsp, err)
		} else {
			mlog.Tracef("耗时(%dms)|Req<%v>|Rsp<%v>", endMs-startMs, req, rsp)
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

func (f Handler[S, T, R]) Rpc(sendrsp define.SendRspFunc, s interface{}, head *pb.Head, buf []byte) func() {
	ref := atomic.AddUint32(&head.Reference, 1)
	return func() {
		// 参数解析
		req := f.NewReq()
		if err := base.Unmarshal[T](buf, req); err != nil {
			mlog.Errorf("参数解析失败%v", err)
			return
		}
		obj, ok := s.(*S)
		if !ok {
			mlog.Errorf("调用%s参数类型错误%v", val2str[head.ActorFunc], s)
			return
		}
		rsp := f.NewRsp()

		// 接口调用
		startMs := time.Now().UnixMilli()
		err := f(obj, head, req, rsp)
		endMs := time.Now().UnixMilli()
		if err != nil {
			mlog.Errorf("耗时(%dms)|Req<%v>|Rsp<%v>|Error<%v>", endMs-startMs, req, rsp, err)
		} else {
			mlog.Tracef("耗时(%dms)|Req<%v>|Rsp<%v>", endMs-startMs, req, rsp)
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
