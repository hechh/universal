package service

import (
	"fmt"
	"universal/common/pb"
	"universal/framework/cluster/domain"
	"universal/framework/cluster/internal/discovery/etcd"
	"universal/framework/cluster/internal/nodes"
	"universal/framework/cluster/internal/routine"
	"universal/framework/fbasic"
	"universal/framework/network"

	"google.golang.org/protobuf/proto"
)

var (
	dis        domain.IDiscovery   // 服务发现
	natsClient *network.NatsClient // nats客户端
	selfNode   *pb.ClusterNode     // 服务自身节点
)

func GetLocalClusterNode() *pb.ClusterNode {
	return selfNode
}

// 初始化
func Init(natsUrl string, ends []string) (err error) {
	// 初始化nats
	if natsClient, err = network.NewNatsClient(natsUrl); err != nil {
		return err
	}
	// 初始化etcd
	if dis, err = etcd.NewEtcdClient(ends...); err != nil {
		return err
	}
	return
}

func watchClusterNode(action int, key string, value string) {
	vv := &pb.ClusterNode{}
	if err := proto.Unmarshal(fbasic.StringToBytes(value), vv); err != nil {
		panic(err)
	}
	switch action {
	case domain.ActionTypeDel:
		// 添加服务节点
		nodes.AddNode(vv)
	default:
		nodes.DeleteNode(vv)
	}
}

// 服务发现
func Discovery(node *pb.ClusterNode) error {
	if node == nil {
		return fbasic.NewUError(1, pb.ErrorCode_Parameter, "*pb.ClusterNode is nil")
	}
	if node.ClusterID <= 0 {
		node.ClusterID = fbasic.GetCrc32(fmt.Sprintf("%s:%d", node.Ip, node.Port))
	}
	selfNode = node
	// 注册自身服务
	key := domain.GetNodeChannel(selfNode.ClusterType, selfNode.ClusterID)
	buf, err := proto.Marshal(selfNode)
	if err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_ProtoMarshal, err)
	}
	dis.KeepAlive(key, buf, 10)
	// 设置监听 + 发现其他服务
	if err := dis.Watch(domain.ROOT_DIR, watchClusterNode); err != nil {
		return err
	}
	return nil
}

// 消息订阅
func Subscribe(f func(*pb.Packet)) (err error) {
	// 订阅单点
	if err = natsClient.Subscribe(domain.GetNodeChannel(selfNode.ClusterType, selfNode.ClusterID), f); err != nil {
		return
	}
	// 订阅广播
	err = natsClient.Subscribe(domain.GetTopicChannel(selfNode.ClusterType), f)
	return
}

// 发送消息
func Publish(pac *pb.Packet) error {
	head := pac.Head
	var key string
	switch head.SendType {
	case pb.SendType_POINT:
		// 先路由
		if err := dispatcher(head); err != nil {
			return err
		}
		key = domain.GetNodeChannel(head.DstClusterType, head.DstClusterID)
	case pb.SendType_BOARDCAST:
		key = domain.GetTopicChannel(head.DstClusterType)
	default:
		return fbasic.NewUError(1, pb.ErrorCode_SendTypeNotSupported, pac.Head.SendType)
	}
	return natsClient.Publish(key, pac)
}

func dispatcher(head *pb.PacketHead) error {
	rlist := routine.GetRoutine(head)
	if rinfo := rlist.Get(head.DstClusterType); rinfo == nil {
		// 路由
		if err := rlist.UpdateRoutine(head, nodes.RandomNode(head)); err != nil {
			return err
		}
	} else {
		if head.DstClusterID <= 0 {
			head.DstClusterID = rinfo.GetClusterID()
		}
		// 节点丢失
		if head.DstClusterID != rinfo.GetClusterID() {
			// 重新路由
			if err := rlist.UpdateRoutine(head, nodes.RandomNode(head)); err != nil {
				return err
			}
		}
		// 判断节点是否存在
		if node := nodes.GetNode(head); node == nil {
			// 重新路由
			if err := rlist.UpdateRoutine(head, nodes.RandomNode(head)); err != nil {
				return err
			}
		} else {
			rlist.UpdateRoutine(head, node)
		}
	}
	return nil
}
