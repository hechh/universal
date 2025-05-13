package network

import (
	"fmt"
	"universal/common/pb"
	"universal/framework/domain"
	"universal/framework/library/mlog"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
)

type Nats struct {
	client   *nats.Conn        // nats连接
	topic    string            // 订阅话题
	routeMgr domain.IRouterMgr // 路由表
}

func NewNats(url string, opts ...OpOption) (*Nats, error) {
	client, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	vals := NewOp(opts...)
	return &Nats{
		client: client,
		topic:  vals.topic,
	}, nil
}

func (n *Nats) broadTopic(t int32) string {
	return fmt.Sprintf("%s/%d", n.topic, t)
}

func (n *Nats) sendTopic(t, id int32) string {
	return fmt.Sprintf("%s/%d/%d", n.topic, t, id)
}

func (n *Nats) Listen(node *pb.Node, act domain.IActorMgr) error {
	if _, err := n.client.Subscribe(n.sendTopic(node.GetType(), node.GetId()), func(msg *nats.Msg) {
		// 解析包
		pack := &pb.Packet{}
		if err := proto.Unmarshal(msg.Data, pack); err != nil {
			mlog.Error("详细解析失败: %v", err)
			return
		}
		// 更新路由表
		n.routeMgr.Set(pack.Head.RouteId, pack.Router)
		// 消息转发
		if err := act.SendRpc(pack.Head, pack.Body); err != nil {
			mlog.Error("请求的 Actor 不存在: %v", err)
		}
	}); err != nil {
		return err
	}

	_, err := n.client.Subscribe(n.broadTopic(node.GetType()), func(msg *nats.Msg) {
		// 解析包
		pack := &pb.Packet{}
		if err := proto.Unmarshal(msg.Data, pack); err != nil {
			mlog.Error("详细解析失败: %v", err)
			return
		}
		// 更新路由表
		n.routeMgr.Set(pack.Head.RouteId, pack.Router)
		// 消息转发
		if err := act.SendRpc(pack.Head, pack.Body); err != nil {
			mlog.Error("请求的 Actor 不存在: %v", err)
		}
	})
	return err
}

// 发送消息
func (n *Nats) Send(head *pb.Head, data []byte) error {
	// 封装消息
	pack := &pb.Packet{
		Head:   head,
		Router: n.routeMgr.Get(head.RouteId),
		Body:   data,
	}
	buf, err := proto.Marshal(pack)
	if err != nil {
		return err
	}
	// 发送消息
	return n.client.Publish(n.sendTopic(head.DstNodeType, head.DstNodeId), buf)
}

// 发送广播
func (n *Nats) Broadcast(head *pb.Head, data []byte) error {
	// 封装消息
	pack := &pb.Packet{
		Head:   head,
		Router: n.routeMgr.Get(head.RouteId),
		Body:   data,
	}
	buf, err := proto.Marshal(pack)
	if err != nil {
		return err
	}
	// 发送消息
	return n.client.Publish(n.broadTopic(head.DstNodeType), buf)
}

func (n *Nats) Close() error {
	if n.client != nil {
		n.client.Close()
	}
	return nil
}
