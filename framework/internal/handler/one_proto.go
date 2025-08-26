package handler

import (
	"sync/atomic"
	"time"
	"universal/common/base"
	"universal/common/pb"
	"universal/framework/define"
	"universal/library/mlog"
)

type OneProto[S any, T any] func(*S, *pb.Head, *T) error

func (f OneProto[S, T]) New() *T {
	return new(T)
}

func (f OneProto[S, T]) Call(sendrsp define.SendRspFunc, s define.IActor, head *pb.Head, args ...interface{}) func() {
	ref := atomic.AddUint32(&head.Reference, 1)
	return func() {
		// 参数解析
		req, ok := any(args[0]).(*T)
		if !ok {
			mlog.Errorf("调用%s.%s参数类型错误%v", head.ActorName, head.FuncName, args[0])
			return
		}
		obj, ok := any(s).(*S)
		if !ok {
			mlog.Errorf("调用%s.%s参数类型错误%v", head.ActorName, head.FuncName, s)
			return
		}

		// 接口调用
		startMs := time.Now().UnixMilli()
		err := f(obj, head, req)
		endMs := time.Now().UnixMilli()
		if err != nil {
			mlog.Errorf("%s.%s耗时(%dms)|Event<%v>|Error<%v>", head.ActorName, head.FuncName, endMs-startMs, args[0], err)
		} else {
			mlog.Tracef("%s.%s耗时(%dms)|Event<%v>", head.ActorName, head.FuncName, endMs-startMs, args[0])
		}

		// 是否回复
		if sendrsp != nil && atomic.CompareAndSwapUint32(&head.Reference, ref, ref) {
			rspProto, ok := any(req).(define.IRspProto)
			if !ok {
				return
			}
			rspProto.SetHead(base.ToRspHead(err))
			if err := sendrsp(head, rspProto); err != nil {
				mlog.Errorf("回复失败|Rsp<%v>|%v", rspProto, err)
			}
		}
	}
}

func (f OneProto[S, T]) Rpc(sendrsp define.SendRspFunc, s define.IActor, head *pb.Head, buf []byte) func() {
	ref := atomic.AddUint32(&head.Reference, 1)
	return func() {
		// 参数解析
		req := f.New()
		if err := base.Unmarshal[T](buf, req); err != nil {
			mlog.Errorf("协议解析失败<%v>", err)
			return
		}
		obj, ok := any(s).(*S)
		if !ok {
			mlog.Errorf("调用%s.%s参数类型错误%v", head.ActorName, head.FuncName, s)
			return
		}

		// 接口调用
		startMs := time.Now().UnixMilli()
		err := f(obj, head, req)
		endMs := time.Now().UnixMilli()
		if err != nil {
			mlog.Errorf("%s.%s耗时(%dms)|Event<%v>|Error<%v>", head.ActorName, head.FuncName, endMs-startMs, req, err)
		} else {
			mlog.Tracef("%s.%s耗时(%dms)|Event<%v>", head.ActorName, head.FuncName, endMs-startMs, req)
		}

		// 是否回复
		if sendrsp != nil && atomic.CompareAndSwapUint32(&head.Reference, ref, ref) {
			rspProto, ok := any(req).(define.IRspProto)
			if !ok {
				return
			}
			rspProto.SetHead(base.ToRspHead(err))
			if err = sendrsp(head, rspProto); err != nil {
				mlog.Errorf("回复失败|Rsp<%v>|%v", rspProto, err)
			}
		}
	}
}
