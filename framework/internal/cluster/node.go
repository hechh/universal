package cluster

import (
	"encoding/json"
	"universal/framework/define"
)

type Node struct {
	Name string `json:"name"` // 节点名称
	Addr string `json:"addr"` // 节点地址
	Type uint32 `json:"type"` // 节点类型
	Id   uint32 `json:"id"`   // 节点ID
}

func NewNode(buf []byte) define.INode {
	node := new(Node)
	json.Unmarshal(buf, node)
	return node
}

func (n *Node) GetName() string {
	return n.Name
}

func (n *Node) GetType() uint32 {
	return n.Type
}

func (n *Node) GetId() uint32 {
	return n.Id
}

func (n *Node) GetAddr() string {
	return n.Addr
}

func (n *Node) ToBytes() []byte {
	buf, _ := json.Marshal(n)
	return buf
}
