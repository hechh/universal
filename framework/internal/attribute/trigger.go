package attribute

import (
	"time"
	"universal/common/pb"
	"universal/framework/define"
	"universal/library/mlog"

	"github.com/golang/protobuf/proto"
)

type Trigger[S any] func(*S, *pb.Head) error

func (f Trigger[S]) Call(sendrsp define.SendRspFunc, s interface{}, head *pb.Head, args ...proto.Message) func() {
	return func() {
		obj, ok := s.(*S)
		if !ok {
			mlog.Errorf("调用%s.%s参数类型错误%v", head.ActorName, head.FuncName, s)
			return
		}

		startMs := time.Now().UnixMilli()
		err := f(obj, head)
		endMs := time.Now().UnixMilli()
		if err != nil {
			mlog.Errorf("耗时(%dms)|Error<%v>", endMs-startMs, err)
		} else {
			mlog.Tracef("耗时(%dms)", endMs-startMs)
		}
	}
}

func (f Trigger[S]) Rpc(sendrsp define.SendRspFunc, s interface{}, head *pb.Head, buf []byte) func() {
	return func() {
		obj, ok := s.(*S)
		if !ok {
			mlog.Errorf("调用%s.%s参数类型错误%v", head.ActorName, head.FuncName, s)
			return
		}

		startMs := time.Now().UnixMilli()
		err := f(obj, head)
		endMs := time.Now().UnixMilli()
		if err != nil {
			mlog.Errorf("耗时(%dms)|Error<%v>", endMs-startMs, err)
		} else {
			mlog.Tracef("耗时(%dms)", endMs-startMs)
		}
	}
}
