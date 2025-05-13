package network

import (
	"fmt"
	"time"
	"universal/framework/domain"
	"universal/framework/library/async"
	"universal/framework/library/mlog"

	"github.com/nsqio/go-nsq"
)

// Nsq 结构体实现 INetwork 接口
type Nsq struct {
	addr      string                // nsqd 地址
	topic     string                // 订阅话题
	producer  *nsq.Producer         // 生产者
	consumer  *nsq.Consumer         // 消费者
	broadcast *nsq.Consumer         // 广播消费者
	routeMgr  domain.IRouterMgr     // 路由表
	newPacket func() domain.IPacket // 解析函数
	newHeader func() domain.IHead   // 创建函数
	newRoute  func() domain.IRouter // 创建函数
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
		routeMgr:  vals.routeMgr,
		newRoute:  vals.newRoute,
		newPacket: vals.newPacket,
		newHeader: vals.newHeader,
	}, nil
}

func (n *Nsq) broadTopic(t int32) string {
	return fmt.Sprintf("%s_%d", n.topic, t)
}

func (n *Nsq) sendTopic(t, id int32) string {
	return fmt.Sprintf("%s_%d_%d", n.topic, t, id)
}

type handler struct {
	nsq *Nsq
	act domain.IActorMgr
}

func (h *handler) HandleMessage(m *nsq.Message) error {
	async.SafeRecover(mlog.Fatal, func() {
		// 解析包
		pack := h.nsq.newPacket()
		if err := pack.ReadFrom(m.Body); err != nil {
			mlog.Error("详细解析失败: %v", err)
			return
		}
		// 更新路由表
		h.nsq.routeMgr.Set(pack.GetHead().GetRouteId(), pack.GetRoute())
		// 消息转发
		if err := h.act.SendRpc(pack.GetHead(), pack.GetBody()); err != nil {
			mlog.Error("请求的 Actor 不存在: %v", err)
		}
	})
	return nil
}

// Read 实现接收消息的功能
func (n *Nsq) Listen(node domain.INode, act domain.IActorMgr) error {
	cfg := nsq.NewConfig()
	cfg.MaxInFlight = 5
	cfg.MsgTimeout = 5 * time.Second
	// 单播
	consumer, err := nsq.NewConsumer(n.sendTopic(node.GetType(), node.GetId()), node.GetAddr(), cfg)
	if err != nil {
		return err
	}
	consumer.AddHandler(&handler{act: act, nsq: n})
	if err := consumer.ConnectToNSQD(n.addr); err != nil {
		return err
	}
	n.consumer = consumer
	// 广播
	bro, err := nsq.NewConsumer(n.broadTopic(node.GetType()), node.GetAddr(), cfg)
	if err != nil {
		return err
	}
	bro.AddHandler(&handler{act: act, nsq: n})
	if err := bro.ConnectToNSQD(n.addr); err != nil {
		return err
	}
	n.broadcast = bro
	return nil
}

func (n *Nsq) Send(head domain.IHead, data []byte) error {
	// 封装消息
	pack := n.newPacket().SetHead(head).SetBody(data).SetRoute(n.routeMgr.Get(head.GetRouteId()))
	buf := make([]byte, pack.GetSize())
	if err := pack.WriteTo(buf); err != nil {
		return err
	}
	// 发送消息
	topic := n.sendTopic(head.GetDstNodeType(), head.GetDstNodeId())
	return n.producer.Publish(topic, buf)
}

func (n *Nsq) Broadcast(head domain.IHead, data []byte) error {
	// 封装消息
	pack := n.newPacket().SetHead(head).SetBody(data).SetRoute(n.routeMgr.Get(head.GetRouteId()))
	buf := make([]byte, pack.GetSize())
	if err := pack.WriteTo(buf); err != nil {
		return err
	}
	// 发送消息
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
