package framework

import (
	"universal/common/pb"
	"universal/common/yaml"
	"universal/framework/actor"
	"universal/framework/internal/service"
	"universal/library/mlog"

	"github.com/golang/protobuf/proto"
)

var (
	core *service.Service
)

func Init(node *pb.Node, cfg *yaml.Config) (err error) {
	core, err = service.NewService(node, cfg)
	if err != nil {
		return
	}

	actor.SetSend(core.SendToClient)
	actor.SetResponse(core.Response)
	return
}

// 跨服务发消息
func Send(head *pb.Head, args ...interface{}) error {
	return core.Send(head, args...)
}

// 跨服务类型广播
func Broadcast(head *pb.Head, args ...interface{}) error {
	return core.Broadcast(head, args...)
}

// 同步请求
func Request(head *pb.Head, msg proto.Message, reply proto.Message) error {
	return core.Request(head, msg, reply)
}

// 发送到客户端
func SendToClient(head *pb.Head, msg proto.Message) error {
	return core.SendToClient(head, msg)
}

// 通知客户端
func NotifyToClient(uids []uint64, head *pb.Head, msg proto.Message) {
	core.NotifyToClient(uids, head, msg)
}

// 注册消息处理函数
func RegisterBroadcastHandler(f func(*pb.Head, []byte)) {
	core.RegisterBroadcastHandler(f)
}

// 注册消息处理函数
func RegisterSendHandler(f func(*pb.Head, []byte)) {
	core.RegisterSendHandler(f)
}

// 注册消息处理函数
func RegisterReplyHandler(f func(*pb.Head, []byte)) {
	core.RegisterReplyHandler(f)
}

// 默认内网消息处理器
func defaultSendHandler(head *pb.Head, buf []byte) {
	mlog.Debugf("send调用: %v", head)
	if err := actor.Send(head, buf); err != nil {
		mlog.Errorf("跨服务调用错误: %v", err)
	}
}

// 默认内网消息处理器
func defaultReplyHandler(head *pb.Head, buf []byte) {
	mlog.Debugf("rpc调用: %v", head)
	if err := actor.Send(head, buf); err != nil {
		mlog.Errorf("跨服务调用错误: %v", err)
	}
}

// 默认内网广播消息处理器
func defaultBroadcastHandler(head *pb.Head, buf []byte) {
	mlog.Debugf("broadcast调用: %v", head)
	if err := actor.Send(head, buf); err != nil {
		mlog.Errorf("跨服务调用错误: %v", err)
	}
}
