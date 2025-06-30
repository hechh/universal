package node

import (
	"sync"
	"universal/common/pb"
	"universal/library/random"
)

type pool struct {
	mutex sync.RWMutex
	nodes map[int32]*pb.Node
}

type Node struct {
	pools map[pb.NodeType]*pool
}

func NewNode() *Node {
	pools := make(map[pb.NodeType]*pool)
	for i := pb.NodeType_NodeTypeBegin + 1; i < pb.NodeType_NodeTypeEnd; i++ {
		pools[i] = &pool{nodes: make(map[int32]*pb.Node)}
	}
	return &Node{pools: pools}
}

func (c *Node) GetCount(nodeType pb.NodeType) int {
	nn := c.pools[nodeType]
	return len(nn.nodes)
}

// 随机获取节点
func (c *Node) Get(nodeType pb.NodeType, nodeId int32) *pb.Node {
	nn := c.pools[nodeType]
	nn.mutex.RLock()
	defer nn.mutex.RUnlock()
	return nn.nodes[nodeId]
}

func (c *Node) Del(nodeType pb.NodeType, nodeId int32) bool {
	nn := c.pools[nodeType]
	nn.mutex.Lock()
	defer nn.mutex.Unlock()
	delete(nn.nodes, nodeId)
	return true
}

func (c *Node) Add(node *pb.Node) bool {
	nn := c.pools[pb.NodeType(node.Type)]
	nn.mutex.Lock()
	defer nn.mutex.Unlock()
	nn.nodes[node.Id] = node
	return true
}

func (c *Node) gets(nodeType pb.NodeType) (rets []*pb.Node) {
	nn := c.pools[nodeType]
	nn.mutex.RLock()
	defer nn.mutex.RUnlock()
	for _, item := range nn.nodes {
		rets = append(rets, item)
	}
	return
}

// 随机获取节点
func (c *Node) Random(nodeType pb.NodeType, seed uint64) *pb.Node {
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
