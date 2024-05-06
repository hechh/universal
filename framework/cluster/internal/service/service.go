package service

import (
	"log"
	"net"
	"universal/common/pb"
	"universal/framework/cluster/domain"
	"universal/framework/cluster/internal/discovery/etcd"
	"universal/framework/cluster/internal/nodes"
	"universal/framework/cluster/internal/routine"
	"universal/framework/fbasic"

	"github.com/spf13/cast"
	"google.golang.org/protobuf/proto"
)

var (
	dis      domain.IDiscovery // 服务发现
	selfNode *pb.ClusterNode   // 服务自身节点
)

func GetLocalClusterNode() *pb.ClusterNode {
	return selfNode
}

func GetDiscovery() domain.IDiscovery {
	return dis
}

func Stop() {
	dis.Close()
}

// 初始化
func Init(ends []string, types ...pb.ClusterType) (err error) {
	// 初始化etcd
	etcd, err := etcd.NewEtcdClient(ends...)
	if err != nil {
		return err
	}
	// 初始化节点类型
	nodes.Init(types...)
	// 设置服务发现
	dis = etcd
	return
}

func watchClusterNode(action int, key string, value string) {
	vv := &pb.ClusterNode{}
	if err := proto.Unmarshal(fbasic.StringToBytes(value), vv); err != nil {
		panic(err)
	}
	switch action {
	case domain.ActionTypeDel:
		clusterType, clusterID, _ := fbasic.ParseChannel(key)
		nodes.Delete(clusterType, clusterID)
	default:
		// 添加服务节点
		nodes.Add(vv)
	}
	log.Println("发现服务节点: ", action, key, vv)
}

// 服务发现
func Discovery(typ pb.ClusterType, addr string) error {
	ip, port, err := net.SplitHostPort(addr)
	if err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_SocketAddr, err)
	}
	// 构建自身节点
	selfNode = &pb.ClusterNode{
		ClusterType: typ,
		ClusterID:   fbasic.GetCrc32(addr),
		Ip:          ip,
		Port:        cast.ToInt32(port),
	}

	// 注册自身服务（保活，服务下线会自动删除）
	buf, err := proto.Marshal(selfNode)
	if err != nil {
		return fbasic.NewUError(1, pb.ErrorCode_ProtoMarshal, err)
	}
	dis.KeepAlive(fbasic.GetNodeChannel(selfNode.ClusterType, selfNode.ClusterID), string(buf), 10)

	// 设置监听 + 发现其他服务
	if err := dis.Watch(fbasic.GetRootDir(), watchClusterNode); err != nil {
		return err
	}
	return nil
}

// 路由
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

/*
// 路由到game集群
func ToDispatcher(head *pb.PacketHead, sendType pb.SendType, dst pb.ClusterType) (*pb.PacketHead, error) {
	// 源节点
	newHead := *head
	newHead.SrcClusterID = selfNode.ClusterID
	newHead.SrcClusterType = selfNode.ClusterType
	// 目的节点
	newHead.DstClusterType = dst
	newHead.SendType = sendType
	// 路由
	if err := dispatcher(&newHead); err != nil {
		return nil, err
	}
	return &newHead, nil
}

func Dispatcher(head *pb.PacketHead) error {
	// 本地信息
	head.SrcClusterType = selfNode.ClusterType
	head.SrcClusterID = selfNode.ClusterID
	head.DstClusterType = fbasic.ApiCodeToClusterType(head.ApiCode)
	// 路由
	if head.DstClusterType == selfNode.ClusterType {
		head.DstClusterID = selfNode.ClusterID
		return nil
	}
	return dispatcher(head)
}
*/
