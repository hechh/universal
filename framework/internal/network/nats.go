package network

import (
	"fmt"
	"universal/framework/define"

	"github.com/nats-io/nats.go"
)

type Nats struct {
	client    *nats.Conn                         // nats连接
	topic     string                             // 订阅话题
	newPacket func() define.IPacket              // 解析函数
	newHeader func(define.ITable) define.IHeader // 创建函数
	newTable  func() define.ITable               // 创建函数
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
		newTable:  vals.newTable,
		newPacket: vals.newPacket,
		newHeader: vals.newHeader,
	}, nil
}

func (n *Nats) broadTopic(t uint32) string {
	return fmt.Sprintf("%s/%d", n.topic, t)
}

func (n *Nats) sendTopic(t uint32, id uint32) string {
	return fmt.Sprintf("%s/%d/%d", n.topic, t, id)
}

func (n *Nats) Read(node define.INode, f func(define.IHeader, []byte)) error {
	_, err := n.client.Subscribe(n.sendTopic(node.GetType(), node.GetId()), func(msg *nats.Msg) {
		pack := n.newPacket()
		pack.SetHeader(n.newHeader(n.newTable()))
		pack.Parse(msg.Data)
		f(pack.GetHeader(), pack.GetBody())
	})
	if err != nil {
		return err
	}

	_, err = n.client.Subscribe(n.broadTopic(node.GetType()), func(msg *nats.Msg) {
		pack := n.newPacket()
		pack.SetHeader(n.newHeader(n.newTable()))
		pack.Parse(msg.Data)
		f(pack.GetHeader(), pack.GetBody())
	})
	return err
}

// 发送消息
func (n *Nats) Send(head define.IHeader, data []byte) error {
	buf := n.newPacket().SetHeader(head).SetBody(data).ToBytes()
	topic := n.sendTopic(head.GetDstNodeType(), head.GetDstNodeId())
	return n.client.Publish(topic, buf)
}

// 发送广播
func (n *Nats) Broadcast(head define.IHeader, data []byte) error {
	buf := n.newPacket().SetHeader(head).SetBody(data).ToBytes()
	topic := n.broadTopic(head.GetDstNodeType())
	return n.client.Publish(topic, buf)
}

func (n *Nats) Close() error {
	if n.client != nil {
		n.client.Close()
	}
	return nil
}
