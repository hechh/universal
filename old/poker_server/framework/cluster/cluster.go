package cluster

import (
	"poker_server/common/pb"
	"poker_server/framework/domain"
	"poker_server/library/encode"
	"poker_server/library/mlog"
	"poker_server/library/uerror"
	"sync/atomic"

	"github.com/golang/protobuf/proto"
)

// 集群
type Cluster struct {
	nodeMgr      domain.INode
	tableMgr     domain.ITable
	discoveryObj domain.IDiscovery
	busObj       domain.IBus
}

func New(n domain.INode, t domain.ITable, b domain.IBus, d domain.IDiscovery) *Cluster {
	return &Cluster{n, t, d, b}
}

func (c *Cluster) Close() {
	c.tableMgr.Close()
	c.discoveryObj.Close()
	c.busObj.Close()
}

func (c *Cluster) GetSelf() *pb.Node {
	return c.nodeMgr.GetSelf()
}

func (c *Cluster) SetBroadcastHandler(f func(*pb.Head, []byte)) error {
	return c.busObj.SetBroadcastHandler(c.nodeMgr.GetSelf(), f)
}

func (c *Cluster) SetSendHandler(f func(*pb.Head, []byte)) error {
	return c.busObj.SetSendHandler(c.nodeMgr.GetSelf(), func(head *pb.Head, body []byte) {
		c.updateRouter(head.Src, head.Dst)
		f(head, body)
	})
}

func (c *Cluster) SetReplyHandler(f func(*pb.Head, []byte)) error {
	return c.busObj.SetReplyHandler(c.nodeMgr.GetSelf(), func(head *pb.Head, body []byte) {
		c.updateRouter(head.Src, head.Dst)
		f(head, body)
	})
}

func (c *Cluster) updateRouter(rrs ...*pb.NodeRouter) {
	for _, rr := range rrs {
		if rr != nil && rr.Router != nil {
			c.tableMgr.GetOrNew(rr.RouterType, rr.ActorId, c.nodeMgr.GetSelf()).SetData(rr.Router)
			rr.Router = nil
		}
	}
}

func (c *Cluster) queryRouter(rrs ...*pb.NodeRouter) {
	for _, rr := range rrs {
		if rr != nil {
			rr.Router = c.tableMgr.GetOrNew(rr.RouterType, rr.ActorId, c.nodeMgr.GetSelf()).GetData()
		}
	}
}

func (c *Cluster) Broadcast(head *pb.Head, args ...interface{}) error {
	if head.Dst == nil {
		return uerror.New(1, pb.ErrorCode_PARAM_INVALID, "参数错误")
	}
	if head.Dst.NodeType <= pb.NodeType_NodeTypeBegin || head.Dst.NodeType >= pb.NodeType_NodeTypeEnd {
		return uerror.New(1, pb.ErrorCode_NODE_TYPE_NOT_SUPPORTED, "节点类型不支持")
	}
	if c.nodeMgr.GetCount(head.Dst.NodeType) <= 0 {
		return uerror.New(1, pb.ErrorCode_NODE_NOT_FOUND, "服务节点不存在")
	}
	buf, err := encode.Marshal(args...)
	if err != nil {
		return err
	}
	return c.busObj.Broadcast(head, buf)
}

func (c *Cluster) Send(head *pb.Head, args ...interface{}) error {
	if err := c.dispatcher(head); err != nil {
		return err
	}
	if head.Src == nil {
		return uerror.New(1, pb.ErrorCode_PARAM_INVALID, "参数错误")
	}
	c.queryRouter(head.Dst, head.Src)
	if head.Cmd > 0 && head.Cmd%2 == 0 {
		atomic.AddUint32(&head.Reference, 1)
	}
	buf, err := encode.Marshal(args...)
	if err != nil {
		return uerror.Err(pb.ErrorCode_MARSHAL_FAILED, err)
	}
	return c.busObj.Send(head, buf)
}

func (c *Cluster) SendToClient(head *pb.Head, msg proto.Message, uids ...uint64) error {
	if head.Uid > 0 {
		uids = append(uids, head.Uid)
	}
	if len(uids) <= 0 {
		return uerror.New(1, pb.ErrorCode_PARAM_INVALID, "接受对象为空")
	}
	buf, err := encode.Marshal(msg)
	if err != nil {
		return uerror.Err(pb.ErrorCode_MARSHAL_FAILED, err)
	}
	if head.Src == nil {
		return uerror.New(1, pb.ErrorCode_PARAM_INVALID, "参数错误")
	}
	c.queryRouter(head.Src)
	atomic.AddUint32(&head.Reference, 1)
	head.Dst = &pb.NodeRouter{NodeType: pb.NodeType_NodeTypeGate}
	for _, uid := range uids {
		head.Dst.ActorId = uid
		if err := c.dispatcher(head); err != nil {
			mlog.Errorf("路由失败:%v", err)
			continue
		}
		c.queryRouter(head.Dst)
		if err := c.busObj.Send(head, buf); err != nil {
			mlog.Errorf("发送客户端失败：%v", err)
		}
	}
	return nil
}

func (c *Cluster) Request(head *pb.Head, msg interface{}, rsp proto.Message) error {
	if err := c.dispatcher(head); err != nil {
		return err
	}
	if head.Src == nil {
		return uerror.New(1, pb.ErrorCode_PARAM_INVALID, "参数错误")
	}
	c.queryRouter(head.Dst, head.Src)
	buf, err := encode.Marshal(msg)
	if err != nil {
		return uerror.Err(pb.ErrorCode_MARSHAL_FAILED, err)
	}
	return c.busObj.Request(head, buf, rsp)
}

func (c *Cluster) Response(head *pb.Head, msg interface{}) error {
	buf, err := encode.Marshal(msg)
	if err != nil {
		return uerror.Err(pb.ErrorCode_MARSHAL_FAILED, err)
	}
	c.queryRouter(head.Dst, head.Src)
	return c.busObj.Response(head, buf)
}

func (c *Cluster) dispatcher(head *pb.Head) error {
	if head.Dst == nil {
		return uerror.New(1, pb.ErrorCode_PARAM_INVALID, "参数错误")
	}
	if head.Dst.NodeType >= pb.NodeType_NodeTypeEnd || head.Dst.NodeType <= pb.NodeType_NodeTypeBegin {
		return uerror.New(1, pb.ErrorCode_NODE_TYPE_NOT_SUPPORTED, "节点类型不支持")
	}
	self := c.nodeMgr.GetSelf()
	if head.Dst.NodeType == self.Type {
		return uerror.New(1, pb.ErrorCode_NODE_TYPE_NOT_SUPPORTED, "禁止同节点类型发送")
	}
	if head.Dst.NodeId > 0 {
		if c.nodeMgr.Get(head.Dst.NodeType, head.Dst.NodeId) != nil {
			return nil
		}
		return uerror.New(1, pb.ErrorCode_NODE_NOT_FOUND, "节点不存在")
	}
	// 从路由表中读取
	dstTab := c.tableMgr.GetOrNew(head.Dst.RouterType, head.Dst.ActorId, self)
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
	return uerror.New(1, pb.ErrorCode_NODE_NOT_FOUND, "节点不存在")
}

func (c *Cluster) SendResponse(head *pb.Head, rsp proto.Message) error {
	if len(head.Reply) > 0 {
		return c.Response(head, rsp)
	}
	if head.Cmd > 0 {
		head.Src = head.Dst
		return c.SendToClient(head, rsp)
	}
	if head.Src != nil && len(head.Src.ActorName) > 0 && len(head.Src.FuncName) > 0 {
		head.Src, head.Dst = head.Dst, head.Src
		return c.Send(head, rsp)
	}
	return nil
}
