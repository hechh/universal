package cluster

import (
	"sync/atomic"
	"universal/common/pb"
	"universal/framework/define"
	"universal/framework/internal/request"
	"universal/library/encode"
	"universal/library/mlog"
	"universal/library/uerror"

	"github.com/golang/protobuf/proto"
)

type Cluster struct {
	nodeMgr  define.INode
	tableMgr define.ITable
	dis      define.IDiscovery
	client   define.IBus
}

func NewCluster(n define.INode, t define.ITable, d define.IDiscovery, c define.IBus) *Cluster {
	return &Cluster{n, t, d, c}
}

func (c *Cluster) GetSelf() *pb.Node {
	return c.nodeMgr.GetSelf()
}

func (c *Cluster) Close() {
	c.tableMgr.Close()
	c.dis.Close()
	c.client.Close()
}

func (c *Cluster) updateRouter(rrs ...*pb.NodeRouter) {
	for _, rr := range rrs {
		if rr != nil && rr.ActorId > 0 {
			c.tableMgr.GetOrNew(rr.ActorId, c.nodeMgr.GetSelf()).SetData(rr.Router)
			rr.Router = nil
		}
	}
}

func (c *Cluster) queryRouter(rrs ...*pb.NodeRouter) {
	for _, rr := range rrs {
		if rr != nil && rr.ActorId > 0 {
			rr.Router = c.tableMgr.GetOrNew(rr.ActorId, c.nodeMgr.GetSelf()).GetData()
		}
	}
}

func (c *Cluster) SetBroadcastHandler(f func(*pb.Head, []byte)) error {
	return c.client.SetBroadcastHandler(c.nodeMgr.GetSelf(), f)
}

func (c *Cluster) SetSendHandler(f func(*pb.Head, []byte)) error {
	return c.client.SetSendHandler(c.nodeMgr.GetSelf(), func(head *pb.Head, body []byte) {
		if err := request.Parse(head, "Player.SendToClient"); err != nil {
			mlog.Errorf("Nats解析消息失败: %v", err)
			return
		}
		c.updateRouter(head.Src, head.Dst)
		f(head, body)
	})
}

func (c *Cluster) SetReplyHandler(f func(*pb.Head, []byte)) error {
	return c.client.SetReplyHandler(c.nodeMgr.GetSelf(), func(head *pb.Head, body []byte) {
		if err := request.Parse(head, "Player.SendToClient"); err != nil {
			mlog.Errorf("Nats解析消息失败: %v", err)
			return
		}
		c.updateRouter(head.Src, head.Dst)
		f(head, body)
	})
}

func (c *Cluster) Broadcast(head *pb.Head, args ...interface{}) error {
	if head.Dst == nil || (head.Dst.NodeType >= pb.NodeType_End && head.Dst.NodeType <= pb.NodeType_Begin) {
		return uerror.New(1, int32(pb.ErrorCode_NodeTypeNotSupported), "节点类型不支持")
	}
	buf, err := encode.Marshal(args...)
	if err != nil {
		return uerror.Err(1, int32(pb.ErrorCode_ProtoMarshalFailed), err)
	}
	return c.client.Broadcast(head, buf)
}

func (c *Cluster) Send(head *pb.Head, args ...interface{}) error {
	if err := c.dispatcher(head); err != nil {
		return err
	}
	c.queryRouter(head.Dst, head.Src)
	buf, err := encode.Marshal(args...)
	if err != nil {
		return uerror.Err(1, int32(pb.ErrorCode_ProtoMarshalFailed), err)
	}
	return c.client.Send(head, buf)
}

func (c *Cluster) SendCmd(head *pb.Head, args ...interface{}) error {
	if err := c.dispatcher(head); err != nil {
		return err
	}
	c.queryRouter(head.Dst, head.Src)
	atomic.AddUint32(&head.Reference, 1)
	buf, err := encode.Marshal(args...)
	if err != nil {
		return uerror.Err(1, int32(pb.ErrorCode_ProtoMarshalFailed), err)
	}
	return c.client.Send(head, buf)
}

func (c *Cluster) SendToClient(head *pb.Head, msg proto.Message, uids ...uint64) error {
	buf, err := encode.Marshal(msg)
	if err != nil {
		return uerror.Err(1, int32(pb.ErrorCode_ProtoMarshalFailed), err)
	}
	if head.Uid > 0 {
		uids = append(uids, head.Uid)
	}
	if head.Cmd%2 == 0 {
		if _, ok := pb.CMD_name[int32(head.Cmd)+1]; ok {
			head.Cmd++
			head.Seq++
		}
	}
	c.queryRouter(head.Src)
	atomic.AddUint32(&head.Reference, 1)
	head.Dst = &pb.NodeRouter{NodeType: pb.NodeType_Gate}
	for _, uid := range uids {
		head.Dst.ActorId = define.UidToActorId(uid)
		if err := c.dispatcher(head); err == nil {
			mlog.Errorf("路由失败:%v", err)
			continue
		}
		c.queryRouter(head.Dst)
		if err := c.client.Send(head, buf); err != nil {
			mlog.Errorf("发送客户端失败：%v", err)
		}
	}
	return nil
}

func (c *Cluster) Request(head *pb.Head, msg interface{}, rsp proto.Message) error {
	if err := c.dispatcher(head); err != nil {
		return err
	}
	c.queryRouter(head.Dst, head.Src)
	buf, err := encode.Marshal(msg)
	if err != nil {
		return uerror.Err(1, int32(pb.ErrorCode_ProtoMarshalFailed), err)
	}
	return c.client.Request(head, buf, rsp)
}

func (c *Cluster) Response(head *pb.Head, msg interface{}) error {
	buf, err := encode.Marshal(msg)
	if err != nil {
		return uerror.Err(1, int32(pb.ErrorCode_ProtoMarshalFailed), err)
	}
	c.queryRouter(head.Dst, head.Src)
	return c.client.Response(head, buf)
}

func (c *Cluster) dispatcher(head *pb.Head) error {
	self := c.nodeMgr.GetSelf()
	if head.Dst == nil || head.Dst.ActorId <= 0 || head.Dst.ActorFunc <= 0 {
		return uerror.New(1, int32(pb.ErrorCode_NodeRouterIsNil), "参数错误")
	}
	if head.Dst.NodeType >= pb.NodeType_End || head.Dst.NodeType <= pb.NodeType_Begin {
		return uerror.New(1, int32(pb.ErrorCode_NodeTypeNotSupported), "节点类型不支持")
	}
	if head.Dst.NodeType == self.Type {
		return uerror.New(1, int32(pb.ErrorCode_NodeTypeInvalid), "禁止同节点类型发送")
	}
	if head.Dst.NodeId > 0 {
		if c.nodeMgr.Get(head.Dst.NodeType, head.Dst.NodeId) != nil {
			return nil
		}
		return uerror.New(1, int32(pb.ErrorCode_NodeNotFound), "节点不存在")
	}

	// 从路由表中查询
	dstTab := c.tableMgr.GetOrNew(head.Dst.ActorId, self)
	if nodeId := dstTab.Get(head.Dst.NodeType); nodeId > 0 {
		if c.nodeMgr.Get(head.Dst.NodeType, nodeId) != nil {
			head.Dst.NodeId = nodeId
			return nil
		}
	}

	// 从集群节点中随机
	if nn := c.nodeMgr.Random(head.Dst.NodeType, head.Dst.ActorId); nn != nil {
		head.Dst.NodeId = nn.Id
		return nil
	}
	return uerror.New(1, int32(pb.ErrorCode_NodeNotFound), "%v", head.Dst)
}

func (c *Cluster) SendResponse(head *pb.Head, rsp proto.Message) error {
	if len(head.Reply) > 0 {
		return c.Response(head, rsp)
	}
	if head.Cmd > 0 {
		head.Src = head.Dst
		return c.SendToClient(head, rsp)
	}
	if head.Src != nil && head.Src.ActorId > 0 && head.Src.ActorFunc > 0 {
		head.Src, head.Dst = head.Dst, head.Src
		return c.Send(head, rsp)
	}
	return nil
}
