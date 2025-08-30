package entity

import (
	"time"
	"universal/common/pb"
	"universal/framework/define"
	"universal/library/encode"
	"universal/library/mlog"
)

type Gob1Handler[S any, T any] func(*S, *pb.Head, T) error

func (f Gob1Handler[S, T]) New() *T {
	return new(T)
}

func (f Gob1Handler[S, T]) Call(sendrsp define.SendRspFunc, s define.IActor, head *pb.Head, args ...interface{}) func() {
	return func() {
		// 参数解析
		obj, ok := any(s).(*S)
		if !ok {
			mlog.Errorf("调用%s.%s参数类型错误%v", head.ActorName, head.FuncName, s)
			return
		}
		arg1, ok := any(args[0]).(T)
		if !ok {
			mlog.Errorf("调用%s.%s参数类型错误%v", head.ActorName, head.FuncName, s)
			return
		}

		// 接口调用
		startMs := time.Now().UnixMilli()
		err := f(obj, head, arg1)
		endMs := time.Now().UnixMilli()
		if err != nil {
			mlog.Errorf("%s.%s耗时(%dms)|Arg1<%v>|Error<%v>", head.ActorName, head.FuncName, endMs-startMs, arg1, err)
		} else {
			mlog.Tracef("%s.%s耗时(%dms)|Arg1<%v>", head.ActorName, head.FuncName, endMs-startMs, arg1)
		}
	}
}

func (f Gob1Handler[S, T]) Rpc(sendrsp define.SendRspFunc, s define.IActor, head *pb.Head, buf []byte) func() {
	return func() {
		// 参数解析
		req := f.New()
		if err := encode.Decode(buf, req); err != nil {
			mlog.Errorf("调用%s.%s参数解析错误%v", head.ActorName, head.FuncName, err)
			return
		}
		obj, ok := any(s).(*S)
		if !ok {
			mlog.Errorf("调用%s.%s参数类型错误%v", head.ActorName, head.FuncName, s)
			return
		}

		// 接口调用
		startMs := time.Now().UnixMilli()
		err := f(obj, head, *req)
		endMs := time.Now().UnixMilli()
		if err != nil {
			mlog.Errorf("%s.%s耗时(%dms)|Arg1<%v>|Error<%v>", head.ActorName, head.FuncName, endMs-startMs, *req, err)
		} else {
			mlog.Tracef("%s.%s耗时(%dms)|Arg1<%v>", head.ActorName, head.FuncName, endMs-startMs, *req)
		}
	}
}
