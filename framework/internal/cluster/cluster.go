package cluster

import (
	"sync"
	"universal/common/pb"
	"universal/library/random"
)

type pool struct {
	mutex sync.RWMutex
	nodes map[int32]*pb.Node
}

type Cluster struct {
	pools map[pb.NodeType]*pool
}

func NewCluster() *Cluster {
	pools := make(map[pb.NodeType]*pool)
	for i := pb.NodeType_NodeTypeBegin + 1; i < pb.NodeType_NodeTypeEnd; i++ {
		pools[i] = &pool{nodes: make(map[int32]*pb.Node)}
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
	return nn.nodes[nodeId]
}

func (c *Cluster) Del(nodeType pb.NodeType, nodeId int32) bool {
	nn := c.pools[nodeType]
	nn.mutex.Lock()
	defer nn.mutex.Unlock()
	delete(nn.nodes, nodeId)
	return true
}

func (c *Cluster) Add(node *pb.Node) bool {
	nn := c.pools[pb.NodeType(node.Type)]
	nn.mutex.Lock()
	defer nn.mutex.Unlock()
	nn.nodes[node.Id] = node
	return true
}

func (c *Cluster) gets(nodeType pb.NodeType) (rets []*pb.Node) {
	nn := c.pools[nodeType]
	nn.mutex.RLock()
	defer nn.mutex.RUnlock()
	for _, item := range nn.nodes {
		rets = append(rets, item)
	}
	return
}

// 随机获取节点
func (c *Cluster) Random(nodeType pb.NodeType, seed uint64) *pb.Node {
	items := c.gets(nodeType)
	llen := len(items)
	if llen <= 0 {
		return nil
	}
	if seed <= 0 {
		return items[random.Int32n(int32(llen))]
	}
	return items[seed%uint64(llen)]
}
