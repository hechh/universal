package bus

import (
	"fmt"
	"time"
	"universal/common/pb"
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

func NewNats(topic string, endpoints string) (cli *Nats, err error) {
	err = util.Retry(3, time.Second, func() error {
		client, err := nats.Connect(endpoints, close, waitRecon, reconErr, disconErr)
		if err == nil {
			cli = &Nats{
				topic:  topic,
				client: client,
			}
		}
		return err
	})
	return
}

func (n *Nats) broadChannel(t pb.NodeType) string {
	return fmt.Sprintf("%s/%d", n.topic, t)
}

func (n *Nats) sendChannel(t pb.NodeType, id int32) string {
	return fmt.Sprintf("%s/%d/%d", n.topic, t, id)
}

func (n *Nats) replyChannel(t pb.NodeType, id int32) string {
	return fmt.Sprintf("%s/rpc/%d/%d", n.topic, t, id)
}

func (n *Nats) SetBroadcastHandler(node *pb.Node, ff func(*pb.Head, []byte)) error {
	_, err := n.client.Subscribe(n.broadChannel(node.Type), func(msg *nats.Msg) {
		pack := &pb.Packet{}
		if err := proto.Unmarshal(msg.Data, pack); err != nil {
			mlog.Errorf("nats解析packet包失败: %v", err)
			return
		}

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

		pack.Head.Reply = msg.Reply
		ff(pack.Head, pack.Body)
	})
	return err
}

func (d *Nats) Broadcast(head *pb.Head, msg []byte) error {
	head.ActorId = 0
	head.ActorName = ""
	head.FuncName = ""
	msgBuf, err := proto.Marshal(&pb.Packet{Head: head, Body: msg})
	if err != nil {
		return err
	}
	return d.client.Publish(d.broadChannel(head.Dst.NodeType), msgBuf)
}

func (d *Nats) Send(head *pb.Head, msg []byte) error {
	head.ActorId = 0
	head.ActorName = ""
	head.FuncName = ""
	msgBuf, err := proto.Marshal(&pb.Packet{Head: head, Body: msg})
	if err != nil {
		return err
	}
	return d.client.Publish(d.sendChannel(head.Dst.NodeType, head.Dst.NodeId), msgBuf)
}

func (d *Nats) Request(head *pb.Head, req []byte, rsp proto.Message) error {
	head.ActorId = 0
	head.ActorName = ""
	head.FuncName = ""
	msgBuf, err := proto.Marshal(&pb.Packet{Head: head, Body: req})
	if err != nil {
		return err
	}
	resp, err := d.client.Request(d.replyChannel(head.Dst.NodeType, head.Dst.NodeId), msgBuf, 3000*time.Millisecond)
	if err != nil {
		return err
	}
	if err := proto.Unmarshal(resp.Data, rsp); err != nil {
		return err
	}
	return nil
}

func (d *Nats) Response(head *pb.Head, msg []byte) error {
	head.ActorId = 0
	head.ActorName = ""
	head.FuncName = ""
	return d.client.Publish(head.Reply, msg)
}

func (n *Nats) Close() {
	if n.client != nil {
		n.client.Close()
	}
}
