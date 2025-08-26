package node

import (
	"hash"
	"hash/fnv"
	"sort"
	"sync"
	"universal/common/pb"
	"universal/library/random"
)

type pool struct {
	mutex   sync.RWMutex
	buckets int
	nodes   map[int32]*pb.Node
}

type Node struct {
	hashPool sync.Pool
	self     *pb.Node
	pools    map[pb.NodeType]*pool
}

func NewNode(nn *pb.Node) *Node {
	pools := make(map[pb.NodeType]*pool)
	for i := pb.NodeType_Begin + 1; i < pb.NodeType_End; i++ {
		pools[i] = &pool{nodes: make(map[int32]*pb.Node), buckets: 128}
	}
	return &Node{
		pools: pools,
		self:  nn,
		hashPool: sync.Pool{
			New: func() interface{} {
				return fnv.New64a()
			},
		},
	}
}

func (c *Node) GetSelf() *pb.Node {
	return c.self
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

// 随机获取节点
func (c *Node) Random(nodeType pb.NodeType, seed uint64) *pb.Node {
	items := c.gets(nodeType)
	sort.Slice(items, func(i, j int) bool {
		return items[i].Id < items[j].Id
	})
	llen := len(items)
	if llen <= 0 {
		return nil
	}
	if seed <= 0 {
		return items[random.Int32n(int32(llen))]
	}
	// hash路由
	h := c.hashPool.Get().(hash.Hash64)
	defer c.hashPool.Put(h)
	h.Reset()
	var b [8]byte
	b[0] = byte(seed >> 56)
	b[1] = byte(seed >> 48)
	b[2] = byte(seed >> 40)
	b[3] = byte(seed >> 32)
	b[4] = byte(seed >> 24)
	b[5] = byte(seed >> 16)
	b[6] = byte(seed >> 8)
	b[7] = byte(seed)
	h.Write(b[:])
	return items[h.Sum64()%uint64(llen)]
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
