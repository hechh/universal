package cluster

import (
	"sync/atomic"
	"universal/common/pb"
	"universal/common/yaml"
	"universal/framework/domain"
	"universal/framework/internal/bus"
	"universal/framework/internal/discovery"
	"universal/framework/internal/node"
	"universal/framework/internal/router"
	"universal/library/encode"
	"universal/library/mlog"
	"universal/library/pprof"
	"universal/library/uerror"

	"github.com/golang/protobuf/proto"
)

var (
	tab  domain.ITable
	cls  domain.INode
	dis  domain.IDiscovery
	buss domain.IBus
)

func Init(cfg *yaml.Config, srvCfg *yaml.NodeConfig, nn *pb.Node) (err error) {
	tab = router.NewTable(srvCfg.RouterTTL)
	cls = node.NewNode(nn)
	pprof.Init("localhost", srvCfg.Port+10000)

	dis, err = discovery.NewEtcd(cfg.Etcd)
	if err != nil {
		return
	}
	if err = dis.Watch(cls); err != nil {
		return
	}
	if err = dis.Register(cls, srvCfg.DiscoveryTTL); err != nil {
		return
	}
	buss, err = bus.NewNats(cfg.Nats)
	return
}

func Close() {
	tab.Close()
	dis.Close()
	buss.Close()
}

func GetSelf() *pb.Node {
	return cls.GetSelf()
}

func UpdateRouter(rrs ...*pb.NodeRouter) {
	for _, rr := range rrs {
		if rr == nil || rr.Router == nil || rr.ActorId <= 0 {
			return
		}
		tab.GetOrNew(rr.ActorId, cls.GetSelf()).SetData(rr.Router)
	}
}

func QueryRouter(rrs ...*pb.NodeRouter) {
	for _, rr := range rrs {
		if rr == nil || rr.Router == nil || rr.ActorId <= 0 {
			return
		}
		rr.Router = tab.GetOrNew(rr.ActorId, cls.GetSelf()).GetData()
	}
}

func SetBroadcastHandler(f func(*pb.Head, []byte)) error {
	return buss.SetBroadcastHandler(cls.GetSelf(), f)
}

func SetSendHandler(f func(*pb.Head, []byte)) error {
	return buss.SetSendHandler(cls.GetSelf(), func(head *pb.Head, body []byte) {
		UpdateRouter(head.Src, head.Dst)
		f(head, body)
	})
}

func SetReplyHandler(f func(*pb.Head, []byte)) error {
	return buss.SetReplyHandler(cls.GetSelf(), func(head *pb.Head, body []byte) {
		UpdateRouter(head.Src, head.Dst)
		f(head, body)
	})
}

func Broadcast(head *pb.Head, args ...interface{}) error {
	QueryRouter(head.Dst, head.Src)
	buf, err := encode.Marshal(args...)
	if err != nil {
		return uerror.Err(1, int32(pb.ErrorCode_ProtoMarshalFailed), err)
	}
	return buss.Broadcast(head, buf)
}

func Send(head *pb.Head, args ...interface{}) error {
	if err := Dispatcher(head); err != nil {
		return err
	}
	QueryRouter(head.Dst, head.Src)
	atomic.AddUint32(&head.Reference, 1)
	buf, err := encode.Marshal(args...)
	if err != nil {
		return uerror.Err(1, int32(pb.ErrorCode_ProtoMarshalFailed), err)
	}
	return buss.Send(head, buf)
}

func Request(head *pb.Head, msg interface{}, rsp proto.Message) error {
	if err := Dispatcher(head); err != nil {
		return err
	}
	QueryRouter(head.Dst, head.Src)
	buf, err := encode.Marshal(msg)
	if err != nil {
		return uerror.Err(1, int32(pb.ErrorCode_ProtoMarshalFailed), err)
	}
	return buss.Request(head, buf, rsp)
}

func Response(head *pb.Head, msg interface{}) error {
	QueryRouter(head.Dst, head.Src)
	buf, err := encode.Marshal(msg)
	if err != nil {
		return uerror.Err(1, int32(pb.ErrorCode_ProtoMarshalFailed), err)
	}
	return buss.Response(head, buf)
}

func SendToClient(head *pb.Head, msg proto.Message, uids ...uint64) error {
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
	QueryRouter(head.Src)
	atomic.AddUint32(&head.Reference, 1)
	head.Dst = &pb.NodeRouter{Type: pb.NodeType_NodeTypeGate}
	for _, uid := range uids {
		head.Dst.ActorId = uid
		if err := Dispatcher(head); err == nil {
			mlog.Errorf("玩家已经不在线 %v", err)
			continue
		}
		QueryRouter(head.Dst)
		if err := buss.Send(head, buf); err != nil {
			mlog.Errorf("发送客户端失败：%v", err)
		}
	}
	return nil
}

func SendResponse(head *pb.Head, rsp proto.Message) error {
	if len(head.Reply) > 0 {
		return Response(head, rsp)
	}
	if head.Cmd > 0 {
		head.Src = head.Dst
		return SendToClient(head, rsp)
	}
	if head.Src != nil && head.Src.ActorId > 0 && head.Src.ActorFunc > 0 {
		head.Src, head.Dst = head.Dst, head.Src
		return Send(head, rsp)
	}
	return nil
}

func Dispatcher(head *pb.Head) error {
	if head.Dst == nil || head.Dst.ActorId <= 0 {
		return uerror.New(1, int32(pb.ErrorCode_NodeRouterIsNil), "%v", head)
	}
	if head.Dst.Type >= pb.NodeType_NodeTypeEnd || head.Dst.Type <= pb.NodeType_NodeTypeBegin {
		return uerror.New(1, int32(pb.ErrorCode_NodeTypeNotSupported), "%v", head.Dst)
	}
	if head.Dst.Type == cls.GetSelf().Type {
		return uerror.New(1, int32(pb.ErrorCode_NodeTypeInvalid), "%v", head.Dst)
	}
	if head.Dst.Id > 0 {
		if cls.Get(head.Dst.Type, head.Dst.Id) != nil {
			return nil
		}
		return uerror.New(1, int32(pb.ErrorCode_NodeNotFound), "%v", head.Dst)
	}
	if nodeId := tab.Get(head.Dst.ActorId).Get(head.Dst.Type); nodeId > 0 {
		if cls.Get(head.Dst.Type, nodeId) != nil {
			head.Dst.Id = nodeId
			return nil
		}
	}
	if nn := cls.Random(head.Dst.Type, head.Dst.ActorId); nn != nil {
		head.Dst.Id = nn.Id
		return nil
	}
	return uerror.New(1, int32(pb.ErrorCode_NodeNotFound), "%v", head.Dst)
}
