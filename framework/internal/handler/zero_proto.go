package handler

import (
	"time"
	"universal/common/pb"
	"universal/framework/define"
	"universal/library/mlog"
)

type ZeroProto[S any] func(*S, *pb.Head) error

func (f ZeroProto[S]) Call(sendrsp define.SendRspFunc, s define.IActor, head *pb.Head, args ...interface{}) func() {
	return func() {
		obj, ok := any(s).(*S)
		if !ok {
			mlog.Errorf("调用%s.%s参数类型错误%v", head.ActorName, head.FuncName, s)
			return
		}

		startMs := time.Now().UnixMilli()
		err := f(obj, head)
		endMs := time.Now().UnixMilli()
		if err != nil {
			mlog.Errorf("调用%s.%s耗时(%dms)|Error<%v>", head.ActorName, head.FuncName, endMs-startMs, err)
		} else {
			mlog.Tracef("调用%s.%s耗时(%dms)", head.ActorName, head.FuncName, endMs-startMs)
		}
	}
}

func (f ZeroProto[S]) Rpc(sendrsp define.SendRspFunc, s define.IActor, head *pb.Head, buf []byte) func() {
	return func() {
		obj, ok := any(s).(*S)
		if !ok {
			mlog.Errorf("调用%s.%s参数类型错误%v", head.ActorName, head.FuncName, s)
			return
		}

		startMs := time.Now().UnixMilli()
		err := f(obj, head)
		endMs := time.Now().UnixMilli()
		if err != nil {
			mlog.Errorf("调用%s.%s耗时(%dms)|Error<%v>", head.ActorName, head.FuncName, endMs-startMs, err)
		} else {
			mlog.Tracef("调用%s.%s耗时(%dms)", head.ActorName, head.FuncName, endMs-startMs)
		}
	}
}
