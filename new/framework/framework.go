package framework

import (
	"poker_server/common/pb"
	"poker_server/common/yaml"
	"poker_server/framework/domain"
	"poker_server/framework/internal/core/actor"
	"poker_server/framework/internal/service"

	"github.com/golang/protobuf/proto"
)

type Actor struct{ actor.Actor }

type ActorPool struct{ actor.ActorPool }

type ActorMgr struct{ actor.ActorMgr }

// 初始化框架
func Init(cfg *yaml.Config) error {
	return service.Init(cfg)
}

// 获取自身节点信息
func GetSelf() *pb.Node {
	return service.GetSelf()
}

func SetSelf(node *pb.Node) {
	service.SetSelf(node)
}

// 注册Actor
func RegisterActor(act domain.IActor) {
	actor.Register(act)
}

func RegisterBroadcastHandler(f func(*pb.Head, []byte)) {
	service.RegisterBroadcastHandler(f)
}

func RegisterSendHandler(f func(*pb.Head, []byte)) {
	service.RegisterSendHandler(f)
}

func RegisterReplyHandler(f func(*pb.Head, []byte)) {
	service.RegisterReplyHandler(f)
}

// 直接调用本地actor
func SendMsg(head *pb.Head, args ...interface{}) error {
	return actor.SendMsg(head, args...)
}

// 直接rpc调用actor
func Send(head *pb.Head, buf []byte) error {
	return actor.Send(head, buf)
}

// 跨服务广播
func BroadcastMsgToNode(head *pb.Head, msg proto.Message) error {
	return service.BroadcastMsgToNode(head, msg)
}

// 跨服务广播
func BroadcastToNode(head *pb.Head, buf []byte) error {
	return service.BroadcastToNode(head, buf)
}

// 远程异步调用
func SendMsgToNode(head *pb.Head, msg proto.Message) error {
	return service.SendMsgToNode(head, msg)
}

// 远程异步调用
func SendToNode(head *pb.Head, body []byte) error {
	return service.SendToNode(head, body)
}

// 远程同步调用
func RequestMsgToNode(head *pb.Head, msg proto.Message, reply proto.Message) error {
	return service.RequestMsgToNode(head, msg, reply)
}

// 远程同步调用
func RequestToNode(head *pb.Head, body []byte, reply proto.Message) error {
	return service.RequestToNode(head, body, reply)
}

// 返回客户端
func SendMsgToClient(head *pb.Head, msg proto.Message) error {
	return service.SendMsgToClient(head, msg)
}

// 返回客户端
func SendToClient(head *pb.Head, body []byte) error {
	return service.SendToClient(head, body)
}

// 通知客户端
func NotifyMsgToClient(uids []uint64, head *pb.Head, msg proto.Message) error {
	return service.NotifyMsgToClient(uids, head, msg)
}

// 通知客户端
func NotifyToClient(uids []uint64, head *pb.Head, buf []byte) error {
	return service.NotifyToClient(uids, head, buf)
}

// 客户端广播
func BroadcastMsgToClient(head *pb.Head, msg proto.Message) error {
	return service.BroadcastMsgToClient(head, msg)
}

// 客户端广播
func BroadcastToClient(head *pb.Head, body []byte) error {
	return service.BroadcastToClient(head, body)
}

func CopyHead(h *pb.Head, a, f string) *pb.Head {
	return &pb.Head{
		SendType:    h.SendType,
		SrcNodeType: h.SrcNodeType,
		SrcNodeId:   h.SrcNodeId,
		DstNodeType: h.DstNodeType,
		DstNodeId:   h.DstNodeId,
		Id:          h.Id,
		RouteId:     h.RouteId,
		Cmd:         h.Cmd,
		ActorName:   a,
		FuncName:    f,
		Reply:       h.Reply,
	}
}
