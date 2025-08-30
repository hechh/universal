package entity

import (
	"time"
	"universal/common/pb"
	"universal/framework/define"
	"universal/library/encode"
	"universal/library/mlog"
)

type Gob3Handler[S any, T any, U any, A any] func(*S, *pb.Head, T, U, A) error

func (f Gob3Handler[S, T, U, A]) NewT() *T {
	return new(T)
}

func (f Gob3Handler[S, T, U, A]) NewU() *U {
	return new(U)
}

func (f Gob3Handler[S, T, U, A]) NewA() *A {
	return new(A)
}

func (f Gob3Handler[S, T, U, A]) Call(sendrsp define.SendRspFunc, s define.IActor, head *pb.Head, args ...interface{}) func() {
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
		arg2, ok := any(args[1]).(U)
		if !ok {
			mlog.Errorf("调用%s.%s参数类型错误%v", head.ActorName, head.FuncName, s)
			return
		}
		arg3, ok := any(args[2]).(A)
		if !ok {
			mlog.Errorf("调用%s.%s参数类型错误%v", head.ActorName, head.FuncName, s)
			return
		}

		// 接口调用
		startMs := time.Now().UnixMilli()
		err := f(obj, head, arg1, arg2, arg3)
		endMs := time.Now().UnixMilli()
		if err != nil {
			mlog.Errorf("%s.%s耗时(%dms)|Arg1<%v>|Arg2<%v>|Arg3<%v>|Error<%v>", head.ActorName, head.FuncName, endMs-startMs, arg1, arg2, arg3, err)
		} else {
			mlog.Tracef("%s.%s耗时(%dms)|Arg1<%v>|Arg2<%v>|Arg3<%v>", head.ActorName, head.FuncName, endMs-startMs, arg1, arg2, arg3)
		}
	}
}

func (f Gob3Handler[S, T, U, A]) Rpc(sendrsp define.SendRspFunc, s define.IActor, head *pb.Head, buf []byte) func() {
	return func() {
		// 参数解析
		arg1 := f.NewT()
		if err := encode.Decode(buf, arg1); err != nil {
			mlog.Errorf("调用%s.%s参数解析错误%v", head.ActorName, head.FuncName, err)
			return
		}
		arg2 := f.NewU()
		if err := encode.Decode(buf, arg2); err != nil {
			mlog.Errorf("调用%s.%s参数解析错误%v", head.ActorName, head.FuncName, err)
			return
		}
		arg3 := f.NewA()
		if err := encode.Decode(buf, arg3); err != nil {
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
		err := f(obj, head, *arg1, *arg2, *arg3)
		endMs := time.Now().UnixMilli()
		if err != nil {
			mlog.Errorf("%s.%s耗时(%dms)|Arg1<%v>|Arg2<%v>|Arg3<%v>|Error<%v>", head.ActorName, head.FuncName, endMs-startMs, *arg1, *arg2, *arg3, err)
		} else {
			mlog.Tracef("%s.%s耗时(%dms)|Arg1<%v>|Arg2<%v>|Arg3<%v>", head.ActorName, head.FuncName, endMs-startMs, *arg1, *arg2, *arg3)
		}
	}
}
