package network

import (
	"fmt"
	"universal/framework/define"
	"universal/library/baselib/safe"
	"universal/library/mlog"

	"github.com/nats-io/nats.go"
)

type Nats struct {
	options *Op        // 选项
	client  *nats.Conn // nats连接
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
	return &Nats{options: vals, client: client}, nil
}

func (n *Nats) getChannel() string {
	node := n.options.cluster.GetSelf()
	return fmt.Sprintf("%s/%d/%d", n.options.root, node.GetType(), node.GetId())
}

func (n *Nats) getTopChannel() string {
	node := n.options.cluster.GetSelf()
	return fmt.Sprintf("%s/%d", n.options.root, node.GetType())
}

func (n *Nats) Receive(f func(define.IHeader, []byte)) error {
	// 单播
	_, err := n.client.Subscribe(n.getChannel(), func(msg *nats.Msg) {
		pack := n.options.parse(msg.Data)
		safe.SafeRecover(mlog.Fatal, func() {
			f(pack.GetHeader(), pack.GetBody())
		})
	})
	if err != nil {
		return err
	}
	// 广播
	_, err = n.client.Subscribe(n.getTopChannel(), func(msg *nats.Msg) {
		pack := n.options.parse(msg.Data)
		safe.SafeRecover(mlog.Fatal, func() {
			f(pack.GetHeader(), pack.GetBody())
		})
	})
	return err
}

func (n *Nats) Send(header define.IHeader, data []byte) error {

	return nil
}
