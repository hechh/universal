package cluster

import (
	"sync"
	"universal/common/pb"
	"universal/library/mlog"
	"universal/library/random"
)

type pool struct {
	mutex sync.RWMutex
	nodes []*pb.Node
}

type Cluster struct {
	pools map[pb.NodeType]*pool
}

func NewCluster() *Cluster {
	pools := make(map[pb.NodeType]*pool)
	for i := pb.NodeType_NodeTypeBegin + 1; i < pb.NodeType_NodeTypeEnd; i++ {
		pools[i] = new(pool)
	}
	return &Cluster{pools: pools}
}

func (c *Cluster) GetCount(nodeType pb.NodeType) int {
	nn := c.pools[nodeType]
	return len(nn.nodes)
}

// 随机获取节点
func (c *Cluster) Get(nodeType pb.NodeType, nodeId int32) *pb.Node {
	nn := c.pools[nodeType]
	nn.mutex.RLock()
	defer nn.mutex.RUnlock()
	for _, val := range nn.nodes {
		if val.GetId() == nodeId {
			return val
		}
	}
	return nil
}

func (c *Cluster) Del(nodeType pb.NodeType, nodeId int32) bool {
	if nn := c.Get(nodeType, nodeId); nn == nil {
		return false
	}

	nn := c.pools[nodeType]
	nn.mutex.Lock()
	defer nn.mutex.Unlock()
	j := -1
	for _, val := range nn.nodes {
		if val.GetId() == nodeId {
			mlog.Debugf("删除服务节点：%s", val.String())
			continue
		}
		j++
		nn.nodes[j] = val
	}
	nn.nodes = nn.nodes[:j+1]
	return true
}

func (c *Cluster) Add(node *pb.Node) bool {
	if nn := c.Get(node.Type, node.Id); nn != nil {
		return false
	}

	nn := c.pools[pb.NodeType(node.Type)]
	nn.mutex.Lock()
	defer nn.mutex.Unlock()
	for i, item := range nn.nodes {
		if item.GetId() == node.GetId() {
			nn.nodes[i] = node
			return true
		}
	}
	nn.nodes = append(nn.nodes, node)
	return true
}

// 随机获取节点
func (c *Cluster) Random(nodeType pb.NodeType, seed uint64) *pb.Node {
	nn := c.pools[nodeType]
	llen := len(nn.nodes)
	if llen <= 0 {
		return nil
	}

	nn.mutex.RLock()
	defer nn.mutex.RUnlock()
	if seed <= 0 {
		return nn.nodes[random.Int32n(int32(llen))]
	}
	return nn.nodes[seed%uint64(llen)]
}
