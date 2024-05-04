package nodes

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"universal/common/pb"
)

type NodeTable struct {
	sync.RWMutex
	list []*pb.ClusterNode
}

var (
	_nodes = make(map[pb.ClusterType]*NodeTable)
)

func Print() {
	for _, tt := range _nodes {
		tt.RLock()
		defer tt.RUnlock()
		for i, node := range tt.list {
			buf, _ := json.Marshal(node)
			fmt.Println(i, "--->", string(buf))
		}
	}
}

func Init(typs ...pb.ClusterType) {
	for _, typ := range typs {
		_nodes[typ] = new(NodeTable)
	}
}

// 获取节点信息
func Get(clusterType pb.ClusterType, clusterID uint32) *pb.ClusterNode {
	if tt, ok := _nodes[clusterType]; !ok {
		return nil
	} else {
		tt.RLock()
		defer tt.RUnlock()
		for _, node := range tt.list {
			if node.ClusterID == clusterID {
				return node
			}
		}
	}
	return nil
}

// 删除节点
func Delete(clusterType pb.ClusterType, clusterID uint32) {
	if tt, ok := _nodes[clusterType]; !ok {
		return
	} else {
		tt.Lock()
		defer tt.Unlock()
		pos := -1
		for _, item := range tt.list {
			if item.ClusterID != clusterID {
				pos++
				tt.list[pos] = item
			}
		}
		tt.list = tt.list[:pos+1]
	}
}

// 添加节点
func Add(node *pb.ClusterNode) {
	// 已经存在
	if nn := Get(node.ClusterType, node.ClusterID); nn != nil {
		return
	}
	// 新建节点
	if tt, ok := _nodes[node.ClusterType]; !ok {
		return
	} else {
		// 插入
		tt.Lock()
		defer tt.Unlock()
		tt.list = append(tt.list, node)
		sort.Slice(tt.list, func(i, j int) bool {
			return tt.list[i].ClusterID < tt.list[j].ClusterID
		})
	}
}

// 随机路由一个节点
func Random(head *pb.PacketHead) *pb.ClusterNode {
	if tt, ok := _nodes[head.DstClusterType]; !ok {
		return nil
	} else if len(tt.list) > 0 {
		return tt.list[int(head.UID)%len(tt.list)]
	}
	return nil
}
