package actor

import (
	"universal/common/pb"
	"universal/framework/domain"
	"universal/library/uerror"

	"github.com/golang/protobuf/proto"
)

var (
	responseFunc func(*pb.Head, interface{}) error
	sendFunc     func(*pb.Head, proto.Message) error
	actors       = make(map[string]domain.IActor)
)

func SetResponse(f func(*pb.Head, interface{}) error) {
	responseFunc = f
}

func SetSend(f func(*pb.Head, proto.Message) error) {
	sendFunc = f
}

func Register(ac domain.IActor) {
	actors[ac.GetActorName()] = ac
}

func SendMsg(head *pb.Head, args ...interface{}) error {
	if act, ok := actors[head.ActorName]; ok {
		return act.SendMsg(head, args...)
	}
	return uerror.New(1, -1, "Actor%s不存在", head.ActorName)
}

func Send(head *pb.Head, body []byte) error {
	if act, ok := actors[head.ActorName]; ok {
		return act.Send(head, body)
	}
	return uerror.New(1, -1, "Actor%s不存在", head.ActorName)
}
