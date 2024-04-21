package balance

import (
	"fmt"
	"sort"
	"sync"
	"universal/common/pb"
	"universal/framework/basic"
)

type Balance struct {
	sync.RWMutex
	nodes map[pb.ClusterType][]*pb.ClusterNode // type -> []clusterNode
}

func NewBalance(types ...pb.ClusterType) *Balance {
	ret := &Balance{
		nodes: make(map[pb.ClusterType][]*pb.ClusterNode),
	}
	for _, typ := range types {
		ret.nodes[typ] = make([]*pb.ClusterNode, 0)
	}
	return ret
}

// 获取节点信息
func (d *Balance) GetNode(head *pb.PacketHead) *pb.ClusterNode {
	d.RLock()
	defer d.RUnlock()
	for _, item := range d.nodes[head.DstClusterType] {
		if item.ClusterID == head.DstClusterID {
			return item
		}
	}
	return nil
}

// 删除节点
func (d *Balance) DelNode(node *pb.ClusterNode) {
	node.ClusterID = basic.GetCrc32(fmt.Sprintf("%s:%d", node.Ip, node.Port))

	d.Lock()
	defer d.Unlock()
	vals := d.nodes[node.ClusterType]
	j := -1
	for i := 0; i < len(vals); i++ {
		if vals[i].ClusterID == node.ClusterID {
			continue
		}
		j++
		vals[j] = vals[i]
	}
	vals = vals[:j+1]
	d.nodes[node.ClusterType] = vals
}

// 添加节点
func (d *Balance) AddNode(node *pb.ClusterNode) {
	node.ClusterID = basic.GetCrc32(fmt.Sprintf("%s:%d", node.Ip, node.Port))

	d.Lock()
	defer d.Unlock()
	for _, item := range d.nodes[node.ClusterType] {
		if item.ClusterID == node.ClusterID {
			return
		}
	}

	d.nodes[node.ClusterType] = append(d.nodes[node.ClusterType], node)

	sort.Slice(d.nodes[node.ClusterType], func(i, j int) bool {
		return d.nodes[node.ClusterType][i].ClusterID < d.nodes[node.ClusterType][j].ClusterID
	})
}

func (d *Balance) RandomNode(head *pb.PacketHead) *pb.ClusterNode {
	d.RLock()
	defer d.RUnlock()
	vals := d.nodes[pb.ClusterType(head.DstClusterType)]
	lvals := len(vals)
	if lvals <= 0 {
		return nil
	}
	return vals[int(head.UID)%lvals]
}
