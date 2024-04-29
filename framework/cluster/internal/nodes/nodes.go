package nodes

import (
	"fmt"
	"sort"
	"sync/atomic"
	"universal/common/pb"
	"universal/framework/fbasic"
)

const (
	CLUSTER_SIZE = 50 // 集群大小
)

type Nodes [CLUSTER_SIZE]*atomic.Value

var (
	_nodes = make(map[pb.ClusterType]Nodes)
)

func (d Nodes) Walk(f func(i int, node *pb.ClusterNode) bool) {
	for i, n := range d {
		val, _ := n.Load().(*pb.ClusterNode)
		if !f(i, val) {
			return
		}
	}
}

func InitNodes(types ...pb.ClusterType) {
	for _, typ := range types {
		list := [CLUSTER_SIZE]*atomic.Value{}
		for i := 0; i < CLUSTER_SIZE; i++ {
			list[i] = new(atomic.Value)
		}
		_nodes[typ] = list
	}
}

// 获取节点信息
func GetNode(head *pb.PacketHead) (ret *pb.ClusterNode) {
	vals, ok := _nodes[head.DstClusterType]
	if !ok {
		return nil
	}
	Nodes(vals).Walk(func(_ int, node *pb.ClusterNode) bool {
		if node != nil && node.ClusterID == head.DstClusterID {
			ret = node
			return false
		}
		return true
	})
	return
}

// 删除节点
func DeleteNode(node *pb.ClusterNode) {
	if node.ClusterID <= 0 {
		node.ClusterID = fbasic.GetCrc32(fmt.Sprintf("%s:%d", node.Ip, node.Port))
	}
	vals, ok := _nodes[node.ClusterType]
	if !ok {
		return
	}
	Nodes(vals).Walk(func(i int, item *pb.ClusterNode) bool {
		if item != nil && item.ClusterID == node.ClusterID {
			vals[i].Store(nil)
			return false
		}
		return true
	})
}

// 添加节点
func AddNode(node *pb.ClusterNode) {
	if node.ClusterID <= 0 {
		node.ClusterID = fbasic.GetCrc32(fmt.Sprintf("%s:%d", node.Ip, node.Port))
	}
	vals, ok := _nodes[node.ClusterType]
	if !ok {
		return
	}
	Nodes(vals).Walk(func(i int, item *pb.ClusterNode) bool {
		if item == nil {
			vals[i].Store(node)
			return false
		}
		return true
	})
}

// 随机路由一个节点
func RandomNode(head *pb.PacketHead) *pb.ClusterNode {
	vals, ok := _nodes[head.DstClusterType]
	if !ok {
		return nil
	}
	// 读取所有节点
	rets := []*pb.ClusterNode{}
	Nodes(vals).Walk(func(i int, item *pb.ClusterNode) bool {
		if item != nil {
			rets = append(rets, item)
		}
		return true
	})
	// 排序
	sort.Slice(rets, func(i, j int) bool {
		return rets[i].ClusterID < rets[j].ClusterID
	})
	return rets[int(head.UID)%len(rets)]
}
