package cluster

import (
	"sync"
	"universal/common/pb"
	"universal/library/mlog"
	"universal/library/random"
)

type pool struct {
	mutex sync.RWMutex
	nodes []*pb.Node // 节点
}

type Cluster struct {
	pools map[pb.NodeType]*pool
}

func New() *Cluster {
	pools := make(map[pb.NodeType]*pool)
	for i := pb.NodeType_Begin + 1; i < pb.NodeType_End; i++ {
		pools[i] = new(pool)
	}
	return &Cluster{pools: pools}
}

func (c *Cluster) GetCount(nodeType pb.NodeType) int {
	nn := c.pools[nodeType]
	return len(nn.nodes)
}

func (c *Cluster) List(nodeType pb.NodeType) (rets []*pb.Node) {
	nn := c.pools[nodeType]
	nn.mutex.RLock()
	defer nn.mutex.RUnlock()
	for _, val := range nn.nodes {
		rets = append(rets, val)
	}
	return
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

// 删除节点
func (c *Cluster) Del(nodeType pb.NodeType, nodeId int32) {
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
}

// 添加节点
func (c *Cluster) Add(node *pb.Node) {
	nn := c.pools[pb.NodeType(node.Type)]
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
