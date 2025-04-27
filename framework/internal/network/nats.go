package network

import (
	"fmt"
	"universal/framework/define"
	"universal/library/baselib/safe"
	"universal/library/mlog"

	"github.com/nats-io/nats.go"
)

type Nats struct {
	topic    string                 // 订阅话题
	client   *nats.Conn             // nats连接
	parseFun define.ParsePacketFunc // 解析函数
	newFun   define.NewPacketFunc   // 创建函数
}

func NewNats(url string, opts ...OpOption) (*Nats, error) {
	client, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	vals := new(Op)
	for _, opt := range opts {
		opt(vals)
	}
	return &Nats{
		topic:    vals.topic,
		parseFun: vals.parse,
		newFun:   vals.newFun,
		client:   client,
	}, nil
}

func (n *Nats) broadcastTopic(node define.INode) string {
	return fmt.Sprintf("%s/%d", n.topic, node.GetType())
}

func (n *Nats) sendTopic(node define.INode) string {
	return fmt.Sprintf("%s/%d/%d", n.topic, node.GetType(), node.GetId())
}

// 接受请求
func (n *Nats) Read(node define.INode, f func(define.IHeader, []byte)) error {
	// 单播
	_, err := n.client.Subscribe(n.sendTopic(node), func(msg *nats.Msg) {
		safe.SafeRecover(mlog.Fatal, func() {
			pack, err := n.parseFun(msg.Data)
			if err != nil {
				panic(err)
			}

			f(pack.GetHeader(), pack.GetBody())
		})
	})
	if err != nil {
		return err
	}
	// 广播
	_, err = n.client.Subscribe(n.broadcastTopic(node), func(msg *nats.Msg) {
		safe.SafeRecover(mlog.Fatal, func() {
			pack, err := n.parseFun(msg.Data)
			if err != nil {
				panic(err)
			}

			f(pack.GetHeader(), pack.GetBody())
		})
	})
	return err
}

// 发送消息
func (n *Nats) Send(node define.INode, head define.IHeader, data []byte) error {
	pack := n.newFun(head, data)
	return n.client.Publish(n.sendTopic(node), pack.ToBytes())
}

// 发送广播
func (n *Nats) Broadcast(node define.INode, head define.IHeader, data []byte) error {
	pack := n.newFun(head, data)
	return n.client.Publish(n.broadcastTopic(node), pack.ToBytes())
}
