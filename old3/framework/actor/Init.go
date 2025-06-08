package actor

import (
	"universal/common/pb"
	"universal/framework/domain"
	"universal/library/async"
	"universal/library/mlog"
	"universal/library/uerror"

	"github.com/golang/protobuf/proto"
)

var (
	node        *pb.Node
	actors      = make(map[string]domain.IActor)
	sendRspFunc func(*pb.Head, proto.Message) error
)

func Init(nn *pb.Node, f func(*pb.Head, proto.Message) error) {
	node = nn
	async.Init(mlog.Fatalf)
	sendRspFunc = f
}

func Register(ac domain.IActor) {
	actors[ac.GetActorName()] = ac
}

func SendMsg(head *pb.Head, args ...interface{}) error {
	if act, ok := actors[head.ActorName]; ok {
		return act.SendMsg(head, args...)
	}
	return uerror.New(1, pb.ErrorCode_ACTOR_NOT_SUPPORTED, "Actor%s不存在", head.ActorName)
}

func Send(head *pb.Head, body []byte) error {
	if act, ok := actors[head.ActorName]; ok {
		return act.Send(head, body)
	}
	return uerror.New(1, pb.ErrorCode_ACTOR_NOT_SUPPORTED, "Actor%s不存在", head.ActorName)
}
