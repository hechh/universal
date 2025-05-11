package cluster

import (
	"sync"
	"universal/framework/domain"
	"universal/library/mlog"
	"universal/library/random"
)

type pool struct {
	mutex sync.RWMutex
	nodes []domain.INode // 节点
}

type Cluster struct {
	pools map[int32]*pool
}

func NewCluster() *Cluster {
	pools := make(map[int32]*pool)
	for i := domain.NodeTypeBegin + 1; i < domain.NodeTypeMax; i++ {
		pools[int32(i)] = new(pool)
	}
	return &Cluster{pools: pools}
}

// 随机获取节点
func (c *Cluster) Get(nodeType, nodeId int32) domain.INode {
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

// 删除节点
func (c *Cluster) Del(nodeType, nodeId int32) {
	nn := c.pools[nodeType]
	nn.mutex.Lock()
	defer nn.mutex.Unlock()
	j := -1
	for _, val := range nn.nodes {
		if val.GetId() == nodeId {
			mlog.Debug("删除服务节点：%s", val.String())
			continue
		}
		j++
		nn.nodes[j] = val
	}
	nn.nodes = nn.nodes[:j+1]
}

// 添加节点
func (c *Cluster) Add(node domain.INode) {
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

// 随机获取节点
func (c *Cluster) Random(nodeType int32, seed uint64) domain.INode {
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
