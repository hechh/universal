package service

import (
	"log"
	"net"
	"strings"
	"universal/common/pb"
	"universal/framework/cluster/domain"
	"universal/framework/cluster/internal/discovery/etcd"
	"universal/framework/cluster/internal/nodes"
	"universal/framework/cluster/internal/routine"
	"universal/framework/fbasic"
	"universal/framework/network"

	"github.com/spf13/cast"
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

func GetDiscovery() domain.IDiscovery {
	return dis
}

func Stop() {
	natsClient.Close()
	dis.Close()
}

// 初始化
func Init(natsUrl string, ends []string, types ...pb.ClusterType) (err error) {
	// 初始化nats
	if natsClient, err = network.NewNatsClient(natsUrl); err != nil {
		return err
	}
	// 初始化etcd
	if dis, err = etcd.NewEtcdClient(ends...); err != nil {
		return err
	}
	nodes.Init(types...)
	return
}

func watchClusterNode(action int, key string, value string) {
	vv := &pb.ClusterNode{}
	if err := proto.Unmarshal(fbasic.StringToBytes(value), vv); err != nil {
		panic(err)
	}
	switch action {
	case domain.ActionTypeDel:
		strs := strings.Split(strings.TrimPrefix(key, domain.ROOT_DIR), "/")
		clusterType := pb.ClusterType(pb.ClusterType_value[strings.ToUpper(strs[0])])
		clusterID := cast.ToUint32(strs[1])
		nodes.Delete(clusterType, clusterID)
	default:
		// 添加服务节点
		nodes.Add(vv)
	}
	log.Println(action, key, "-----watch----->", vv)
}

// 服务发现
func Discovery(typ pb.ClusterType, addr string) error {
	ip, port, err := net.SplitHostPort(addr)
	if err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_SocketAddr, err)
	}
	selfNode = &pb.ClusterNode{
		ClusterType: typ,
		ClusterID:   fbasic.GetCrc32(addr),
		Ip:          ip,
		Port:        cast.ToInt32(port),
	}
	// 注册自身服务
	key := domain.GetNodeChannel(selfNode.ClusterType, selfNode.ClusterID)
	buf, err := proto.Marshal(selfNode)
	if err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_ProtoMarshal, err)
	}
	dis.KeepAlive(key, string(buf), 10)
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
		if err := Dispatcher(head); err != nil {
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

func Dispatcher(head *pb.PacketHead) error {
	rlist := routine.GetRoutine(head)
	if rinfo := rlist.Get(head.DstClusterType); rinfo == nil {
		// 路由
		if err := rlist.UpdateRoutine(head, nodes.Random(head)); err != nil {
			return err
		}
	} else {
		if head.DstClusterID <= 0 {
			head.DstClusterID = rinfo.ClusterID
		}
		// 节点丢失
		if head.DstClusterID != rinfo.ClusterID {
			// 重新路由
			if err := rlist.UpdateRoutine(head, nodes.Random(head)); err != nil {
				return err
			}
		}
		// 判断节点是否存在
		if node := nodes.Get(head.DstClusterType, head.DstClusterID); node == nil {
			// 重新路由
			if err := rlist.UpdateRoutine(head, nodes.Random(head)); err != nil {
				return err
			}
		} else {
			rlist.UpdateRoutine(head, node)
		}
	}
	return nil
}
