package network

import (
	"fmt"
	"time"
	"universal/framework/define"
	"universal/library/baselib/safe"
	"universal/library/mlog"

	"github.com/nsqio/go-nsq"
)

// Nsq 结构体实现 INetwork 接口
type Nsq struct {
	addr      string                 // nsqd 地址
	topic     string                 // 订阅话题
	producer  *nsq.Producer          // 生产者
	consumer  *nsq.Consumer          // 消费者
	broadcast *nsq.Consumer          // 广播消费者
	parseFun  define.ParsePacketFunc // 解析函数
	newFun    define.NewPacketFunc   // 创建函数
}

// NewNSQNetwork 创建一个新的 Nsq 实例
func NewNsq(nsqdAddr string, opts ...OpOption) (*Nsq, error) {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(nsqdAddr, config)
	if err != nil {
		return nil, err
	}
	vals := &Op{}
	for _, opt := range opts {
		opt(vals)
	}
	return &Nsq{
		addr:     nsqdAddr,
		topic:    vals.topic,
		parseFun: vals.parse,
		newFun:   vals.newFun,
		producer: producer,
	}, nil
}

func (n *Nsq) broadcastTopic(node define.INode) string {
	return fmt.Sprintf("%s_%d", n.topic, node.GetType())
}

func (n *Nsq) sendTopic(node define.INode) string {
	return fmt.Sprintf("%s_%d_%d", n.topic, node.GetType(), node.GetId())
}

type handler struct {
	nsq      *Nsq
	listener func(define.IHeader, []byte)
}

func (h *handler) HandleMessage(m *nsq.Message) error {
	safe.SafeRecover(mlog.Fatal, func() {
		pack, err := h.nsq.parseFun(m.Body)
		if err != nil {
			panic(err)
		}
		h.listener(pack.GetHeader(), pack.GetBody())
	})
	return nil
}

// Read 实现接收消息的功能
func (n *Nsq) Read(node define.INode, listen func(define.IHeader, []byte)) error {
	cfg := nsq.NewConfig()
	cfg.MaxInFlight = 5
	cfg.MsgTimeout = 5 * time.Second

	// 单播
	consumer, err := nsq.NewConsumer(n.sendTopic(node), node.GetName(), cfg)
	if err != nil {
		return err
	}
	consumer.AddHandler(&handler{listener: listen, nsq: n})
	if err := consumer.ConnectToNSQD(n.addr); err != nil {
		return err
	}
	n.consumer = consumer

	// 广播
	bro, err := nsq.NewConsumer(n.broadcastTopic(node), node.GetName(), cfg)
	if err != nil {
		return err
	}
	bro.AddHandler(&handler{listener: listen})
	if err := bro.ConnectToNSQD(n.addr); err != nil {
		return err
	}
	n.broadcast = bro
	return nil
}

func (n *Nsq) Send(node define.INode, head define.IHeader, body []byte) error {
	pack := n.newFun(head, body)
	return n.producer.Publish(n.sendTopic(node), pack.ToBytes())
}

func (n *Nsq) Broadcast(node define.INode, head define.IHeader, body []byte) error {
	pack := n.newFun(head, body)
	return n.producer.Publish(n.broadcastTopic(node), pack.ToBytes())
}
