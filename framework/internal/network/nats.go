package network

import (
	"fmt"
	"universal/framework/define"
	"universal/library/baselib/safe"
	"universal/library/mlog"

	"github.com/nats-io/nats.go"
)

type Nats struct {
	root     string
	parseFun define.ParsePacketFunc
	cluster  define.ICluster
	table    define.ITable
	client   *nats.Conn // nats连接
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
		root:     vals.root,
		parseFun: vals.parse,
		cluster:  vals.cluster,
		table:    vals.table,
		client:   client,
	}, nil
}

// 接受请求
func (n *Nats) Receive(f func(define.IHeader, []byte)) error {
	// 单播
	node := n.cluster.GetSelf()
	sendChannel := fmt.Sprintf("%s/%d/%d", n.root, node.GetType(), node.GetId())
	_, err := n.client.Subscribe(sendChannel, func(msg *nats.Msg) {
		pack := n.parseFun(msg.Data)
		safe.SafeRecover(mlog.Fatal, func() {
			f(pack.GetHeader(), pack.GetBody())
		})
	})
	if err != nil {
		return err
	}
	// 广播
	topChannel := fmt.Sprintf("%s/%d", n.root, node.GetType())
	_, err = n.client.Subscribe(topChannel, func(msg *nats.Msg) {
		pack := n.parseFun(msg.Data)
		safe.SafeRecover(mlog.Fatal, func() {
			f(pack.GetHeader(), pack.GetBody())
		})
	})
	return err
}

// 发送消息
func (n *Nats) Send(header define.IHeader, data []byte) error {
	/*
		node := n.cluster.GetSelf()
		id := header.GetRouteId(header.GetDstType())
		dstNode := n.table.Get(id)

		// 从路由表中加载
		if nodeId := n.table.Get(id); nodeId > 0 {

			return nil
		}
	*/
	return nil
}

// 发送广播
