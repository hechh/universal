package network

import (
	"fmt"
	"universal/framework/domain"
	"universal/framework/library/mlog"

	"github.com/nats-io/nats.go"
)

type Nats struct {
	client    *nats.Conn            // nats连接
	topic     string                // 订阅话题
	routeMgr  domain.IRouterMgr     // 路由表
	newPacket func() domain.IPacket // 解析函数
	newHeader func() domain.IHead   // 创建函数
	newRoute  func() domain.IRouter // 创建函数
}

func NewNats(url string, opts ...OpOption) (*Nats, error) {
	client, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	vals := NewOp(opts...)
	return &Nats{
		client:    client,
		topic:     vals.topic,
		newRoute:  vals.newRoute,
		newPacket: vals.newPacket,
		newHeader: vals.newHeader,
	}, nil
}

func (n *Nats) broadTopic(t int32) string {
	return fmt.Sprintf("%s/%d", n.topic, t)
}

func (n *Nats) sendTopic(t, id int32) string {
	return fmt.Sprintf("%s/%d/%d", n.topic, t, id)
}

func (n *Nats) Listen(node domain.INode, act domain.IActorMgr) error {
	if _, err := n.client.Subscribe(n.sendTopic(node.GetType(), node.GetId()), func(msg *nats.Msg) {
		// 解析包
		pack := n.newPacket()
		if err := pack.ReadFrom(msg.Data); err != nil {
			mlog.Error("详细解析失败: %v", err)
			return
		}
		// 更新路由表
		n.routeMgr.Set(pack.GetHead().GetRouteId(), pack.GetRoute())
		// 消息转发
		if err := act.SendRpc(pack.GetHead(), pack.GetBody()); err != nil {
			mlog.Error("请求的 Actor 不存在: %v", err)
		}
	}); err != nil {
		return err
	}

	_, err := n.client.Subscribe(n.broadTopic(node.GetType()), func(msg *nats.Msg) {
		// 解析包
		pack := n.newPacket()
		if err := pack.ReadFrom(msg.Data); err != nil {
			mlog.Error("详细解析失败: %v", err)
			return
		}
		// 更新路由表
		n.routeMgr.Set(pack.GetHead().GetRouteId(), pack.GetRoute())
		// 消息转发
		if err := act.SendRpc(pack.GetHead(), pack.GetBody()); err != nil {
			mlog.Error("请求的 Actor 不存在: %v", err)
		}
	})
	return err
}

// 发送消息
func (n *Nats) Send(head domain.IHead, data []byte) error {
	// 封装消息
	pack := n.newPacket().SetHead(head).SetBody(data).SetRoute(n.routeMgr.Get(head.GetRouteId()))
	buf := make([]byte, pack.GetSize())
	if err := pack.WriteTo(buf); err != nil {
		return err
	}
	// 发送消息
	topic := n.sendTopic(head.GetDstNodeType(), head.GetDstNodeId())
	return n.client.Publish(topic, buf)
}

// 发送广播
func (n *Nats) Broadcast(head domain.IHead, data []byte) error {
	// 封装消息
	pack := n.newPacket().SetHead(head).SetBody(data).SetRoute(n.routeMgr.Get(head.GetRouteId()))
	buf := make([]byte, pack.GetSize())
	if err := pack.WriteTo(buf); err != nil {
		return err
	}
	// 发送消息
	topic := n.broadTopic(head.GetDstNodeType())
	return n.client.Publish(topic, buf)
}

func (n *Nats) Close() error {
	if n.client != nil {
		n.client.Close()
	}
	return nil
}
