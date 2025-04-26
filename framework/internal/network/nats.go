package network

import (
	"fmt"
	"universal/framework/define"
	"universal/library/baselib/safe"
	"universal/library/mlog"

	"github.com/nats-io/nats.go"
)

type Nats struct {
	root   string
	client *nats.Conn // nats连接
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
		root:   vals.root,
		client: client,
	}, nil
}

func (n *Nats) broadcastChannel(node define.INode) string {
	return fmt.Sprintf("%s/%d", n.root, node.GetType())
}

func (n *Nats) sendChannel(node define.INode) string {
	return fmt.Sprintf("%s/%d/%d", n.root, node.GetType(), node.GetId())
}

// 接受请求
func (n *Nats) Read(node define.INode, f func([]byte)) error {
	// 单播
	_, err := n.client.Subscribe(n.sendChannel(node), func(msg *nats.Msg) {
		safe.SafeRecover(mlog.Fatal, func() { f(msg.Data) })
	})
	if err != nil {
		return err
	}
	// 广播
	_, err = n.client.Subscribe(n.broadcastChannel(node), func(msg *nats.Msg) {
		safe.SafeRecover(mlog.Fatal, func() { f(msg.Data) })
	})
	return err
}

// 发送消息
func (n *Nats) Send(node define.INode, data []byte) error {
	return n.client.Publish(n.sendChannel(node), data)
}

// 发送广播
func (n *Nats) Broadcast(node define.INode, data []byte) error {
	return n.client.Publish(n.broadcastChannel(node), data)
}
