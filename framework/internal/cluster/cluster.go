package cluster

import (
	"sync"
	"universal/framework/define"
	"universal/library/mlog"
	"universal/library/random"
)

type NodePool struct {
	mutex sync.RWMutex
	nodes []define.INode // 节点
}

type Cluster struct {
	self  define.INode
	pools map[uint32]*NodePool
}

func NewCluster(self define.INode) *Cluster {
	pools := make(map[uint32]*NodePool)
	for i := define.NodeTypeBegin + 1; i < define.NodeTypeMax; i++ {
		pools[uint32(i)] = new(NodePool)
	}
	return &Cluster{self: self, pools: pools}
}

// 获取自身节点
func (c *Cluster) GetSelf() define.INode {
	return c.self
}

// 随机获取节点
func (c *Cluster) Get(nodeType, nodeId uint32) define.INode {
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

// 添加节点
func (c *Cluster) Put(node define.INode) {
	nn := c.pools[node.GetType()]
	nn.mutex.Lock()
	defer nn.mutex.Unlock()
	for i, item := range nn.nodes {
		if item.GetId() == node.GetId() {
			nn.nodes[i] = node
			return
		}
	}
	nn.nodes = append(nn.nodes, node)
}

// 删除节点
func (c *Cluster) Del(nodeType, nodeId uint32) {
	nn := c.pools[nodeType]
	nn.mutex.Lock()
	defer nn.mutex.Unlock()
	j := -1
	for _, val := range nn.nodes {
		if val.GetId() == nodeId {
			mlog.Debug("删除服务节点：%s", string(val.ToBytes()))
			continue
		}
		j++
		nn.nodes[j] = val
	}
	nn.nodes = nn.nodes[:j+1]
}

// 随机获取节点
func (c *Cluster) Random(nodeType uint32, seed uint64) define.INode {
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
