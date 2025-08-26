package cluster

import (
	"sync/atomic"
	"universal/common/pb"
	"universal/common/yaml"
	"universal/framework/define"
	"universal/framework/internal/bus"
	"universal/framework/internal/discovery"
	"universal/framework/internal/node"
	"universal/framework/internal/router"
	"universal/library/encode"
	"universal/library/mlog"
	"universal/library/uerror"

	"github.com/golang/protobuf/proto"
)

var (
	nodeMgr      define.INode
	tableMgr     define.ITable
	discoveryObj define.IDiscovery
	busObj       define.IBus
)

func Init(nn *pb.Node, srvCfg *yaml.NodeConfig, cfg *yaml.Config) (err error) {
	nodeMgr = node.NewNode(nn)
	tableMgr = router.NewTable(srvCfg.RouterTTL)

	// 服务注册与发现
	discoveryObj, err = discovery.NewEtcdMonitor(cfg.Etcd.Topic, cfg.Etcd.Endpoints)
	if err != nil {
		return
	}
	if err := discoveryObj.Watch(nodeMgr); err != nil {
		return err
	}
	if err := discoveryObj.Register(nodeMgr, srvCfg.DiscoveryTTL); err != nil {
		return err
	}

	// 消息中间件
	busObj, err = bus.NewNats(cfg.Nats.Topic, cfg.Nats.Endpoints)
	return
}

func Close() {
	tableMgr.Close()
	discoveryObj.Close()
	busObj.Close()
}

func SetBroadcastHandler(f func(*pb.Head, []byte)) error {
	return busObj.SetBroadcastHandler(nodeMgr.GetSelf(), f)
}

func SetSendHandler(f func(*pb.Head, []byte)) error {
	return busObj.SetSendHandler(nodeMgr.GetSelf(), func(head *pb.Head, body []byte) {
		updateRouter(head.Src, head.Dst)
		f(head, body)
	})
}

func SetReplyHandler(f func(*pb.Head, []byte)) error {
	return busObj.SetReplyHandler(nodeMgr.GetSelf(), func(head *pb.Head, body []byte) {
		updateRouter(head.Src, head.Dst)
		f(head, body)
	})
}

func updateRouter(rrs ...*pb.NodeRouter) {
	for _, rr := range rrs {
		if rr != nil && rr.Router != nil {
			tableMgr.GetOrNew(rr.RouterType, rr.RouterId, nodeMgr.GetSelf()).SetData(rr.Router)
			rr.Router = nil
		}
	}
}

func queryRouter(rrs ...*pb.NodeRouter) {
	for _, rr := range rrs {
		if rr != nil {
			rr.Router = tableMgr.GetOrNew(rr.RouterType, rr.RouterId, nodeMgr.GetSelf()).GetData()
		}
	}
}

func Broadcast(head *pb.Head, args ...interface{}) error {
	if head.Dst == nil {
		return uerror.New(1, -1, "参数错误")
	}
	if head.Dst.NodeType <= pb.NodeType_Begin || head.Dst.NodeType >= pb.NodeType_End {
		return uerror.New(1, -1, "节点类型不支持")
	}
	if nodeMgr.GetCount(head.Dst.NodeType) <= 0 {
		return uerror.New(1, -1, "服务节点不存在")
	}
	buf, err := encode.Marshal(args...)
	if err != nil {
		return err
	}
	return busObj.Broadcast(head, buf)
}

func Send(head *pb.Head, args ...interface{}) error {
	if err := dispatcher(head); err != nil {
		return err
	}
	if head.Src == nil {
		return uerror.New(1, -1, "Src参数错误")
	}
	queryRouter(head.Dst, head.Src)
	if head.Cmd > 0 && head.Cmd%2 == 0 {
		atomic.AddUint32(&head.Reference, 1)
	}
	buf, err := encode.Marshal(args...)
	if err != nil {
		return uerror.Err(1, -1, err)
	}
	return busObj.Send(head, buf)
}

func Request(head *pb.Head, msg interface{}, rsp proto.Message) error {
	if err := dispatcher(head); err != nil {
		return err
	}
	queryRouter(head.Dst, head.Src)
	buf, err := encode.Marshal(msg)
	if err != nil {
		return uerror.Err(1, -1, err)
	}
	return busObj.Request(head, buf, rsp)
}

func Response(head *pb.Head, msg interface{}) error {
	buf, err := encode.Marshal(msg)
	if err != nil {
		return uerror.Err(1, -1, err)
	}
	queryRouter(head.Dst, head.Src)
	return busObj.Response(head, buf)
}

func dispatcher(head *pb.Head) error {
	if head.Dst == nil {
		return uerror.New(1, -1, "Dst参数错误")
	}
	if head.Dst.NodeType >= pb.NodeType_End || head.Dst.NodeType <= pb.NodeType_Begin {
		return uerror.New(1, -1, "Dst节点类型不支持")
	}
	self := nodeMgr.GetSelf()
	if head.Dst.NodeType == self.Type {
		return uerror.New(1, -1, "Dst禁止发送相同节点类型")
	}
	if head.Dst.NodeId > 0 {
		if nodeMgr.Get(head.Dst.NodeType, head.Dst.NodeId) != nil {
			return nil
		}
		return uerror.New(1, -1, "Dst节点不存在")
	}
	defer mlog.Debugf("Dispatcher %v", head)
	// 从路由表中读取
	dstTab := tableMgr.GetOrNew(head.Dst.RouterType, head.Dst.RouterId, self)
	if nodeId := dstTab.Get(head.Dst.NodeType); nodeId > 0 {
		if nodeMgr.Get(head.Dst.NodeType, nodeId) != nil {
			head.Dst.NodeId = nodeId
			return nil
		}
	}
	// 从集群节点中随机
	if nn := nodeMgr.Random(head.Dst.NodeType, head.Dst.RouterId); nn != nil {
		head.Dst.NodeId = nn.Id
		return nil
	}
	return uerror.New(1, -1, "Dst路由节点不存在")
}

func SendToClient(head *pb.Head, msg proto.Message, uids ...uint64) error {
	if head.Uid > 0 {
		uids = append(uids, head.Uid)
	}
	if len(uids) <= 0 {
		return nil
	}
	buf, err := encode.Marshal(msg)
	if err != nil {
		return uerror.Err(1, -1, err)
	}
	self := nodeMgr.GetSelf()
	if self.Type == pb.NodeType_Gate {
		return uerror.New(1, -1, "禁止相同节点类型发送")
	}
	atomic.AddUint32(&head.Reference, 1)
	tmps := map[uint64]struct{}{}
	for _, uid := range uids {
		if _, ok := tmps[uid]; !ok {
			tmps[uid] = struct{}{}
		} else {
			continue
		}
		dstTab := tableMgr.Get(pb.RouterType_UID, uid)
		if dstTab == nil {
			mlog.Warnf("玩家路由不存在 uid:%d", uid)
			continue
		}
		head.Uid = uid
		head.Dst = &pb.NodeRouter{
			NodeType:   pb.NodeType_Gate,
			NodeId:     dstTab.Get(pb.NodeType_Gate),
			RouterType: pb.RouterType_UID,
			RouterId:   uid,
			Router:     dstTab.GetData(),
		}
		if err := busObj.Send(head, buf); err != nil {
			mlog.Errorf("发送客户端失败：%v", err)
		}
	}
	return nil
}

func SendResponse(head *pb.Head, rsp proto.Message) (err error) {
	if len(head.Reply) > 0 {
		err = Response(head, rsp)
		return
	}
	if head.Cmd > 0 {
		head.Src = head.Dst
		err = SendToClient(head, rsp)
		return
	}
	if head.Src != nil && head.Src.ActorFunc > 0 {
		head.Src, head.Dst = head.Dst, head.Src
		err = Send(head, rsp)
	}
	return
}
