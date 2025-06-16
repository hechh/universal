package bus

import (
	"fmt"
	"time"
	"universal/common/pb"
	"universal/common/yaml"
	"universal/library/mlog"
	"universal/library/util"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
)

var (
	disconErr = nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
		mlog.Errorf("nats disconnect: %v", err)
	})
	reconErr = nats.ReconnectErrHandler(func(_ *nats.Conn, err error) {
		mlog.Errorf("nats reconnect: %v", err)
	})
	waitRecon = nats.ReconnectWait(5)
	close     = nats.ClosedHandler(func(_ *nats.Conn) {
		mlog.Infof("nats close")
	})
)

type Nats struct {
	topic  string
	client *nats.Conn
}

func NewNats(cfg *yaml.NatsConfig) (*Nats, error) {
	var client *nats.Conn
	if err := util.Retry(3, time.Second, func() error {
		cli, err := nats.Connect(cfg.Endpoints, close, waitRecon, reconErr, disconErr)
		if err == nil {
			client = cli
		}
		return err
	}); err != nil {
		return nil, err
	}
	return &Nats{client: client, topic: cfg.Topic}, nil
}

func (n *Nats) broadcastChannel(t pb.NodeType) string {
	return fmt.Sprintf("%s/%d", n.topic, t)
}

func (n *Nats) sendChannel(t pb.NodeType, id int32) string {
	return fmt.Sprintf("%s/%d/%d", n.topic, t, id)
}

func (n *Nats) replyChannel(t pb.NodeType, id int32) string {
	return fmt.Sprintf("%s/rpc/%d/%d", n.topic, t, id)
}

func (n *Nats) SetBroadcastHandler(node *pb.Node, ff func(*pb.Head, []byte)) error {
	_, err := n.client.Subscribe(n.broadcastChannel(node.Type), func(msg *nats.Msg) {
		pack := &pb.Packet{}
		if err := proto.Unmarshal(msg.Data, pack); err != nil {
			mlog.Errorf("nats解析packet包失败: %v", err)
			return
		}
		mlog.Debugf("收到Nats广播数据包 bodySize:%d, pack:%v", len(msg.Data), pack)

		ff(pack.Head, pack.Body)
	})
	return err
}

func (n *Nats) SetSendHandler(node *pb.Node, ff func(*pb.Head, []byte)) error {
	_, err := n.client.Subscribe(n.sendChannel(node.Type, node.Id), func(msg *nats.Msg) {
		pack := &pb.Packet{}
		if err := proto.Unmarshal(msg.Data, pack); err != nil {
			mlog.Errorf("nats解析packet包失败: %v", err)
			return
		}
		mlog.Debugf("收到Nats单播数据包 bodySize:%d, pack:%v", len(msg.Data), pack)

		ff(pack.Head, pack.Body)
	})
	return err
}

func (n *Nats) SetReplyHandler(node *pb.Node, ff func(*pb.Head, []byte)) error {
	_, err := n.client.Subscribe(n.replyChannel(node.Type, node.Id), func(msg *nats.Msg) {
		pack := &pb.Packet{}
		if err := proto.Unmarshal(msg.Data, pack); err != nil {
			mlog.Errorf("nats解析packet包失败: %v", err)
			return
		}
		mlog.Debugf("收到Nats广播数据包 bodySize:%d, pack:%v", len(msg.Data), pack)

		pack.Head.Reply = msg.Reply
		ff(pack.Head, pack.Body)
	})
	return err
}

func (d *Nats) Broadcast(head *pb.Head, msg []byte) error {
	msgBuf, err := proto.Marshal(&pb.Packet{Head: head, Body: msg})
	if err != nil {
		return err
	}
	return d.client.Publish(d.broadcastChannel(head.Dst.NodeType), msgBuf)
}

func (d *Nats) Send(head *pb.Head, msg []byte) error {
	msgBuf, err := proto.Marshal(&pb.Packet{Head: head, Body: msg})
	if err != nil {
		return err
	}
	return d.client.Publish(d.sendChannel(head.Dst.NodeType, head.Dst.NodeId), msgBuf)
}

func (d *Nats) Request(head *pb.Head, req []byte, rsp proto.Message) error {
	msgBuf, err := proto.Marshal(&pb.Packet{Head: head, Body: req})
	if err != nil {
		return err
	}

	resp, err := d.client.Request(d.replyChannel(head.Dst.NodeType, head.Dst.NodeId), msgBuf, 3000*time.Millisecond)
	if err != nil {
		return err
	}
	return proto.Unmarshal(resp.Data, rsp)
}

func (d *Nats) Response(head *pb.Head, msg []byte) error {
	return d.client.Publish(head.Reply, msg)
}

func (n *Nats) Close() error {
	if n.client != nil {
		n.client.Close()
	}
	return nil
}
