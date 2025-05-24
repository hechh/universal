package service

import (
	"universal/common/pb"
	"universal/common/yaml"
	"universal/framework/domain"
	"universal/framework/internal/core/bus"
	"universal/framework/internal/core/cluster"
	"universal/framework/internal/core/discovery"
	"universal/framework/internal/core/router"
	"universal/library/encode"
	"universal/library/mlog"
	"universal/library/uerror"

	"github.com/golang/protobuf/proto"
)

type Service struct {
	node         *pb.Node          // 节点信息
	clusterObj   domain.ICluster   // 集群节点
	tableObj     domain.ITable     // 路由表
	discoveryObj domain.IDiscovery // 服务发现
	busObj       domain.IBus       // 消息总线
}

func NewService(node *pb.Node, cfg *yaml.Config) (*Service, error) {
	nodeCfg := cfg.Cluster[node.Name]

	clusterObj := cluster.New()
	tableObj := router.New()
	tableObj.SetExpire(nodeCfg.RouterExpire)

	// 服务发现
	dis, err := discovery.NewEtcd(cfg.Etcd)
	if err != nil {
		return nil, err
	}
	if err := dis.Watch(clusterObj); err != nil {
		return nil, err
	}
	if err := dis.Register(node, nodeCfg.DicoveryExpire); err != nil {
		return nil, err
	}
	// 消息中间件
	busObj, err := bus.NewNats(cfg.Nats, tableObj)
	if err != nil {
		return nil, err
	}
	return &Service{
		node:         node,
		clusterObj:   clusterObj,
		tableObj:     tableObj,
		discoveryObj: dis,
		busObj:       busObj,
	}, nil
}

func (d *Service) GetNode() *pb.Node {
	return d.node
}

func (d *Service) RegisterBroadcastHandler(f func(*pb.Head, []byte)) error {
	return d.busObj.SetBroadcastHandler(d.node, f)
}

func (d *Service) RegisterSendHandler(f func(*pb.Head, []byte)) error {
	return d.busObj.SetSendHandler(d.node, f)
}

func (d *Service) RegisterReplyHandler(f func(*pb.Head, []byte)) error {
	return d.busObj.SetReplyHandler(d.node, f)
}

func (d *Service) Send(head *pb.Head, args ...interface{}) (err error) {
	// 检测参数
	if head.RouteId <= 0 || head.Id <= 0 {
		return uerror.New(1, -1, "唯一ID为空: %v", head)
	}
	if head.DstNodeType >= pb.NodeType_End || head.DstNodeType <= pb.NodeType_Begin {
		return uerror.New(1, -1, "服务类型不支持: %v", head)
	}
	if d.clusterObj.GetCount(head.DstNodeType) <= 0 {
		return uerror.New(1, -1, "未找到服务节点: %v", head)
	}

	// 做路由分发
	head.SendType = pb.SendType_POINT
	if err := d.dispatcher(head); err != nil {
		return err
	}

	// 检测参数
	if head.DstNodeType == d.node.Type && head.DstNodeId == d.node.Id {
		return uerror.New(1, -1, "不能发送给自身节点: %v", head)
	}

	// 解析参数
	var buf []byte
	if len(args) <= 0 {
	} else if msg, ok := args[0].(proto.Message); len(args) == 1 && ok {
		buf, err = proto.Marshal(msg)
		if err != nil {
			return uerror.New(1, -1, "序列化失败：%v", err)
		}
	} else {
		buf = encode.Encode(args...)
	}

	// 发送请求
	return d.busObj.Send(head, buf)
}

func (d *Service) Broadcast(head *pb.Head, args ...interface{}) (err error) {
	// 检测参数
	if head.DstNodeType >= pb.NodeType_End || head.DstNodeType <= pb.NodeType_Begin {
		return uerror.New(1, -1, "服务类型不支持: %v", head)
	}
	if d.clusterObj.GetCount(head.DstNodeType) <= 0 {
		return uerror.New(1, -1, "未找到服务节点: %v", head)
	}

	// 设置值
	head.SendType = pb.SendType_BROADCAST
	head.SrcNodeType = d.node.Type
	head.SrcNodeId = d.node.Id

	// 更新路由
	if head.Id > 0 {
		router := d.tableObj.Get(head.IdType, head.Id)
		router.Set(d.node.Type, d.node.Id)
	}

	// 解析参数
	var buf []byte
	if len(args) <= 0 {
	} else if msg, ok := args[0].(proto.Message); len(args) == 1 && ok {
		buf, err = proto.Marshal(msg)
		if err != nil {
			return uerror.New(1, -1, "序列化失败：%v", err)
		}
	} else {
		buf = encode.Encode(args...)
	}

	// 发送请求
	return d.busObj.Broadcast(head, buf)
}

func (d *Service) Request(head *pb.Head, msg proto.Message, reply proto.Message) error {
	// 检测参数
	if head.RouteId <= 0 || head.Id <= 0 {
		return uerror.New(1, -1, "唯一ID为空: %v", head)
	}
	if head.DstNodeType >= pb.NodeType_End || head.DstNodeType <= pb.NodeType_Begin {
		return uerror.New(1, -1, "服务类型不支持: %v", head)
	}
	if d.clusterObj.GetCount(head.DstNodeType) <= 0 {
		return uerror.New(1, -1, "未找到服务节点: %v", head)
	}

	// 做路由分发
	head.SendType = pb.SendType_RPC
	if err := d.dispatcher(head); err != nil {
		return err
	}

	// 检测参数
	if head.DstNodeType == d.node.Type && head.DstNodeId == d.node.Id {
		return uerror.New(1, -1, "不能发送给自身节点: %v", head)
	}

	// 解析参数
	buf, err := proto.Marshal(msg)
	if err != nil {
		return uerror.New(1, -1, "序列化失败：%v", err)
	}

	// 发送请求
	return d.busObj.Request(head, buf, reply)
}

func (d *Service) Response(head *pb.Head, msg proto.Message) error {
	if len(head.Reply) <= 0 {
		return nil
	}
	head.SendType = pb.SendType_RPC

	// 解析参数
	buf, err := proto.Marshal(msg)
	if err != nil {
		return uerror.New(1, -1, "序列化失败：%v", err)
	}

	// 发送请求
	return d.busObj.Response(head, buf)
}

func (d *Service) dispatcher(head *pb.Head) error {
	head.SrcNodeType = d.node.Type
	head.SrcNodeId = d.node.Id
	router := d.tableObj.Get(head.IdType, head.Id)
	router.Set(d.node.Type, d.node.Id)
	// 业务层直接指定具体节点
	if head.DstNodeId > 0 {
		if d.clusterObj.Get(head.DstNodeType, head.DstNodeId) != nil {
			router.Set(head.DstNodeType, head.DstNodeId)
			return nil
		}
		return uerror.New(1, -1, "未找到服务节点: %v", head)
	}
	// 优先从路由中选择
	if nodeId := router.Get(head.DstNodeType); nodeId > 0 {
		if d.clusterObj.Get(head.DstNodeType, nodeId) != nil {
			head.DstNodeId = nodeId
			return nil
		}
	}
	//从集群中随机获取一个节点
	if node := d.clusterObj.Random(head.DstNodeType, head.RouteId); node != nil {
		head.DstNodeId = node.Id
		router.Set(head.DstNodeType, node.Id)
		return nil
	}
	return uerror.New(1, -1, "未找到服务节点: %v", head)
}

func (d *Service) SendToClient(head *pb.Head, msg proto.Message) error {
	head.SendType = pb.SendType_POINT
	head.SrcNodeType = d.node.Type
	head.SrcNodeId = d.node.Id
	head.DstNodeType = pb.NodeType_Gate

	// 检测参数
	if head.Id <= 0 {
		return uerror.New(1, -1, "唯一ID为空: %v", head)
	}

	// 读取路由节点
	router := d.tableObj.Get(head.IdType, head.Id)
	head.DstNodeId = router.Get(head.DstNodeType)

	// 判断节点是否还在
	if d.clusterObj.Get(head.DstNodeType, head.DstNodeId) == nil {
		return uerror.New(1, -1, "未找到服务节点: %v", head)
	}

	// 解析参数
	buf, err := proto.Marshal(msg)
	if err != nil {
		return uerror.New(1, -1, "序列化失败：%v", err)
	}

	// 发送请求
	return d.busObj.Send(head, buf)
}

func (d *Service) NotifyToClient(uids []uint64, head *pb.Head, msg proto.Message) {
	head.SendType = pb.SendType_POINT
	head.SrcNodeType = d.node.Type
	head.SrcNodeId = d.node.Id
	head.DstNodeType = pb.NodeType_Gate

	for _, uid := range uids {
		head.Id = uid
		// 读取路由节点
		router := d.tableObj.Get(head.IdType, head.Id)
		head.DstNodeId = router.Get(head.DstNodeType)

		// 判断节点是否还在
		if d.clusterObj.Get(head.DstNodeType, head.DstNodeId) == nil {
			mlog.Errorf("未找到服务节点: %v", head)
			continue
		}

		// 序列化数据
		buf, err := proto.Marshal(msg)
		if err != nil {
			mlog.Errorf("序列化失败：%v", err)
			continue
		}

		// 发送
		if err := d.busObj.Send(head, buf); err != nil {
			mlog.Errorf("通知玩家失败：%v, error:%v", head, err)
		}
	}
}
