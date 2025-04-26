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

// 接受请求
func (n *Nats) Receive(f func(define.IHeader, []byte)) error {
	// 单播
	node := n.options.cluster.GetSelf()
	sendChannel := fmt.Sprintf("%s/%d/%d", n.options.root, node.GetType(), node.GetId())
	_, err := n.client.Subscribe(sendChannel, func(msg *nats.Msg) {
		pack := n.options.parse(msg.Data)
		safe.SafeRecover(mlog.Fatal, func() {
			f(pack.GetHeader(), pack.GetBody())
		})
	})
	if err != nil {
		return err
	}
	// 广播
	topChannel := fmt.Sprintf("%s/%d", n.options.root, node.GetType())
	_, err = n.client.Subscribe(topChannel, func(msg *nats.Msg) {
		pack := n.options.parse(msg.Data)
		safe.SafeRecover(mlog.Fatal, func() {
			f(pack.GetHeader(), pack.GetBody())
		})
	})
	return err
}

// 发送消息
func (n *Nats) Send(header define.IHeader, data []byte) error {
	node := n.options.cluster.GetSelf()
	mgr := n.options.routerMgr
	id := header.GetRouteId(header.GetDstType())

	// 从路由表中加载
	nodeId := mgr.Get(id, node.GetType())

	return nil
}

// 发送广播
