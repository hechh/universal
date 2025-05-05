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
	addr      string                             // nsqd 地址
	topic     string                             // 订阅话题
	producer  *nsq.Producer                      // 生产者
	consumer  *nsq.Consumer                      // 消费者
	broadcast *nsq.Consumer                      // 广播消费者
	newPacket func() define.IPacket              // 解析函数
	newHeader func(define.ITable) define.IHeader // 创建函数
	newTable  func() define.ITable               // 创建函数
}

// NewNSQNetwork 创建一个新的 Nsq 实例
func NewNsq(nsqdAddr string, opts ...OpOption) (*Nsq, error) {
	producer, err := nsq.NewProducer(nsqdAddr, nsq.NewConfig())
	if err != nil {
		return nil, err
	}
	vals := NewOp(opts...)
	return &Nsq{
		addr:      nsqdAddr,
		topic:     vals.topic,
		producer:  producer,
		newTable:  vals.newTable,
		newPacket: vals.newPacket,
		newHeader: vals.newHeader,
	}, nil
}

func (n *Nsq) broadTopic(t uint32) string {
	return fmt.Sprintf("%s_%d", n.topic, t)
}

func (n *Nsq) sendTopic(t, id uint32) string {
	return fmt.Sprintf("%s_%d_%d", n.topic, t, id)
}

type handler struct {
	nsq      *Nsq
	listener func(define.IHeader, []byte)
}

func (h *handler) HandleMessage(m *nsq.Message) error {
	safe.SafeRecover(mlog.Fatal, func() {
		pack := h.nsq.newPacket().SetHeader(h.nsq.newHeader(h.nsq.newTable())).Parse(m.Body)
		h.listener(pack.GetHeader(), pack.GetBody())
	})
	return nil
}

// Read 实现接收消息的功能
func (n *Nsq) Read(node define.INode, listen func(define.IHeader, []byte)) error {
	cfg := nsq.NewConfig()
	cfg.MaxInFlight = 5
	cfg.MsgTimeout = 5 * time.Second
	consumer, err := nsq.NewConsumer(n.sendTopic(node.GetType(), node.GetId()), node.GetName(), cfg)
	if err != nil {
		return err
	}
	consumer.AddHandler(&handler{listener: listen, nsq: n})
	if err := consumer.ConnectToNSQD(n.addr); err != nil {
		return err
	}
	n.consumer = consumer

	bro, err := nsq.NewConsumer(n.broadTopic(node.GetType()), node.GetName(), cfg)
	if err != nil {
		return err
	}
	bro.AddHandler(&handler{listener: listen, nsq: n})
	if err := bro.ConnectToNSQD(n.addr); err != nil {
		return err
	}
	n.broadcast = bro
	return nil
}

func (n *Nsq) Send(head define.IHeader, body []byte) error {
	buf := n.newPacket().SetHeader(n.newHeader(n.newTable())).SetBody(body).ToBytes()
	topic := n.sendTopic(head.GetDstNodeType(), head.GetDstNodeId())
	return n.producer.Publish(topic, buf)
}

func (n *Nsq) Broadcast(head define.IHeader, body []byte) error {
	buf := n.newPacket().SetHeader(n.newHeader(n.newTable())).SetBody(body).ToBytes()
	topic := n.broadTopic(head.GetDstNodeType())
	return n.producer.Publish(topic, buf)
}

func (n *Nsq) Close() error {
	if n.consumer != nil {
		n.consumer.Stop()
	}
	if n.broadcast != nil {
		n.broadcast.Stop()
	}
	if n.producer != nil {
		n.producer.Stop()
	}
	return nil
}
