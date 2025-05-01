package cluster

import (
	"fmt"
	"sync"
	"universal/framework/define"
	"universal/library/random"
)

type NodePool struct {
	mutex sync.RWMutex
	nodes []define.INode // 节点
}

type Cluster struct {
	self  define.INode
	pools map[int32]*NodePool
}

func NewCluster(self define.INode) *Cluster {
	pools := make(map[int32]*NodePool)
	for i := define.NodeTypeBegin + 1; i < define.NodeTypeMax; i++ {
		pools[int32(i)] = new(NodePool)
	}
	return &Cluster{self: self, pools: pools}
}

// 获取自身节点
func (c *Cluster) GetSelf() define.INode {
	return c.self
}

// 随机获取节点
func (c *Cluster) Get(nodeType, nodeId int32) define.INode {
	return c.pools[nodeType].get(nodeId)
}

// 添加节点
func (c *Cluster) Put(node define.INode) (err error) {
	if err = c.pools[node.GetType()].put(node); err == nil {
		// mlog.Debug("添加服务节点：%s", string(node.ToBytes()))
	}
	return
}

// 删除节点
func (c *Cluster) Del(nodeType, nodeId int32) error {
	return c.pools[nodeType].del(nodeId)
}

// 随机获取节点
func (c *Cluster) Random(nodeType int32, seed uint64) define.INode {
	return c.pools[nodeType].rand(seed)
}

func (c *NodePool) get(id int32) define.INode {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	for _, val := range c.nodes {
		if val.GetId() == id {
			return val
		}
	}
	return nil
}

func (c *NodePool) put(val define.INode) error {
	if node := c.get(val.GetId()); node != nil {
		return fmt.Errorf("服务节点已经存在: %v", val)
	}
	c.mutex.Lock()
	c.nodes = append(c.nodes, val)
	c.mutex.Unlock()
	return nil
}

func (c *NodePool) del(id int32) error {
	if node := c.get(id); node == nil {
		return nil
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	j := -1
	for _, val := range c.nodes {
		if val.GetId() == id {
			//mlog.Debug("删除服务节点：%s", string(val.ToBytes()))
			continue
		}
		j++
		c.nodes[j] = val
	}
	c.nodes = c.nodes[:j+1]
	return nil
}

func (c *NodePool) rand(id uint64) define.INode {
	lnodes := len(c.nodes)
	if lnodes <= 0 {
		return nil
	}
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if id <= 0 {
		return c.nodes[random.Int32n(int32(lnodes))]
	}
	return c.nodes[id%uint64(lnodes)]
}
