package bus

import (
	"fmt"
	"poker_server/common/pb"
	"poker_server/common/yaml"
	"poker_server/framework/domain"
	"poker_server/framework/library/mlog"
	"poker_server/framework/library/uerror"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
)

type Nats struct {
	client *nats.Conn
	topic  string
	table  domain.ITable
}

func NewNats(cfg *yaml.NatsConfig, table domain.ITable) (*Nats, error) {
	disconErr := nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
		mlog.Errorf("nats disconnect: %v", err)
	})
	reconErr := nats.ReconnectErrHandler(func(_ *nats.Conn, err error) {
		mlog.Errorf("nats reconnect: %v", err)
	})
	waitRecon := nats.ReconnectWait(5)
	close := nats.ClosedHandler(func(_ *nats.Conn) {
		mlog.Infof("nats close")
	})
	var client *nats.Conn
	var err error
	for i := 0; i < 3; i++ {
		client, err = nats.Connect(cfg.Endpoints, close, waitRecon, reconErr, disconErr)
		if err == nil {
			break
		} else {
			mlog.Errorf("nats服务连接失败：%v", err)
		}
	}
	if err != nil {
		return nil, uerror.New(1, -1, "nats服务连接失败: %v", err)
	}
	return &Nats{
		client: client,
		topic:  cfg.Channel,
		table:  table,
	}, nil
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

// 广播
func (n *Nats) SetBroadcastHandler(node *pb.Node, ff func(*pb.Head, []byte)) error {
	if ff == nil {
		return uerror.New(1, -1, "nats广播回调函数不能为空")
	}
	if _, err := n.client.Subscribe(n.broadcastChannel(node.Type), func(msg *nats.Msg) {
		pack := &pb.Packet{}
		if err := proto.Unmarshal(msg.Data, pack); err != nil {
			mlog.Errorf("nats解析packet包失败: %v", err)
		} else {
			if pack.Head.Id > 0 {
				n.table.Add(pack.Head.IdType, pack.Head.Id, pack.Router)
			}
			ff(pack.Head, pack.Body)
		}
	}); err != nil {
		return uerror.New(1, -1, "nats订阅广播失败: %v", err)
	}
	return nil
}

// 单播
func (n *Nats) SetSendHandler(node *pb.Node, ff func(*pb.Head, []byte)) error {
	if ff == nil {
		return uerror.New(1, -1, "nats单播回调函数不能为空")
	}
	if _, err := n.client.Subscribe(n.sendChannel(node.Type, node.Id), func(msg *nats.Msg) {
		pack := &pb.Packet{}
		if err := proto.Unmarshal(msg.Data, pack); err != nil {
			mlog.Errorf("nats解析packet包失败: %v", err)
		} else {
			n.table.Add(pack.Head.IdType, pack.Head.Id, pack.Router)
			ff(pack.Head, pack.Body)
		}
	}); err != nil {
		return uerror.New(1, -1, "nats订阅单播失败: %v", err)
	}
	return nil
}

// rpc调用
func (n *Nats) SetReplyHandler(node *pb.Node, ff func(*pb.Head, []byte)) error {
	if ff == nil {
		return uerror.New(1, -1, "nats的rpc回调函数不能为空")
	}
	if _, err := n.client.Subscribe(n.replyChannel(node.Type, node.Id), func(msg *nats.Msg) {
		pack := &pb.Packet{}
		if err := proto.Unmarshal(msg.Data, pack); err != nil {
			mlog.Errorf("nats解析packet包失败: %v", err)
		} else {
			pack.Head.Reply = msg.Reply
			n.table.Add(pack.Head.IdType, pack.Head.Id, pack.Router)
			ff(pack.Head, pack.Body)
		}
	}); err != nil {
		return uerror.New(1, -1, "nats订阅rpc调用失败: %v", err)
	}
	return nil
}

// 发送广播
func (d *Nats) Broadcast(head *pb.Head, msg []byte) error {
	pack := &pb.Packet{Head: head, Body: msg}
	if head.Id > 0 {
		router, err := d.table.Get(head.IdType, head.Id)
		if err != nil {
			return err
		}
		pack.Router = router.GetData()
	}
	msgBuf, err := proto.Marshal(pack)
	if err != nil {
		return uerror.New(1, -1, "nats待发送的msg序列化失败: %v", err)
	}
	return d.client.Publish(d.broadcastChannel(head.DstNodeType), msgBuf)
}

// 发送请求
func (d *Nats) Send(head *pb.Head, msg []byte) error {
	router, err := d.table.Get(head.IdType, head.Id)
	if err != nil {
		return err
	}
	pack := &pb.Packet{
		Head:   head,
		Router: router.GetData(),
		Body:   msg,
	}
	msgBuf, err := proto.Marshal(pack)
	if err != nil {
		return uerror.New(1, -1, "nats待发送的msg序列化失败: %v", err)
	}
	return d.client.Publish(d.sendChannel(head.DstNodeType, head.DstNodeId), msgBuf)
}

// 发送同步请求
func (d *Nats) Request(head *pb.Head, req []byte, rsp proto.Message) error {
	router, err := d.table.Get(head.IdType, head.Id)
	if err != nil {
		return err
	}
	pack := &pb.Packet{
		Head:   head,
		Router: router.GetData(),
		Body:   req,
	}
	msgBuf, err := proto.Marshal(pack)
	if err != nil {
		return uerror.New(1, -1, "nats.Request序列化失败: %v", err)
	}
	// 发送不同请求
	resp, err := d.client.Request(d.replyChannel(head.DstNodeType, head.DstNodeId), msgBuf, 500*time.Millisecond)
	if err != nil {
		return uerror.New(1, -1, "nats.Request发送请求失败: %v", err)
	}
	if err := proto.Unmarshal(resp.Data, rsp); err != nil {
		return uerror.New(1, -1, "nats.Request解析返回值失败: %v", err)
	}
	return nil
}

// 不对外开放接口
func (d *Nats) Response(head *pb.Head, msg proto.Message) error {
	head.SrcNodeType, head.DstNodeType = head.DstNodeType, head.SrcNodeType
	head.ActorName = ""
	head.FuncName = ""
	// 读取路由信息
	router, err := d.table.Get(head.IdType, head.Id)
	if err != nil {
		return err
	}
	// 组包
	buf, err := proto.Marshal(msg)
	if err != nil {
		return uerror.New(1, -1, "nats.Response序列化失败: %v", err)
	}
	pack := &pb.Packet{
		Head:   head,
		Router: router.GetData(),
		Body:   buf,
	}
	msgBuf, err := proto.Marshal(pack)
	if err != nil {
		return uerror.New(1, -1, "nats.Response序列化packet失败: %v", err)
	}
	return d.client.Publish(head.Reply, msgBuf)
}

func (n *Nats) Close() error {
	if n.client != nil {
		n.client.Close()
	}
	return nil
}
