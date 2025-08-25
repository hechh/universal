package attribute

import (
	"time"
	"universal/common/pb"
	"universal/framework/define"
	"universal/library/mlog"

	"github.com/golang/protobuf/proto"
)

type Trigger func(*pb.Head) error

func (f Trigger) Call(sendrsp define.SendRspFunc, head *pb.Head, args ...proto.Message) func() {
	return func() {
		startMs := time.Now().UnixMilli()
		err := f(head)
		endMs := time.Now().UnixMilli()
		if err != nil {
			mlog.Errorf("耗时(%dms)|Error<%v>", endMs-startMs, err)
		} else {
			mlog.Tracef("耗时(%dms)", endMs-startMs)
		}
	}
}

func (f Trigger) Rpc(sendrsp define.SendRspFunc, head *pb.Head, buf []byte) func() {
	return func() {
		startMs := time.Now().UnixMilli()
		err := f(head)
		endMs := time.Now().UnixMilli()
		if err != nil {
			mlog.Errorf("耗时(%dms)|Error<%v>", endMs-startMs, err)
		} else {
			mlog.Tracef("耗时(%dms)", endMs-startMs)
		}
	}
}
