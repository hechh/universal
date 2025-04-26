package cluster

import (
	"encoding/json"
	"universal/framework/define"
)

type Node struct {
	Name string `json:"name"` // 节点名称
	Type int32  `json:"type"` // 节点类型
	Id   int32  `json:"id"`   // 节点ID
	Addr string `json:"addr"` // 节点地址
}

func NewNode(buf []byte) define.INode {
	node := new(Node)
	json.Unmarshal(buf, node)
	return node
}

func (n *Node) GetName() string {
	return n.Name
}

func (n *Node) GetType() int32 {
	return n.Type
}

func (n *Node) GetId() int32 {
	return n.Id
}

func (n *Node) GetAddr() string {
	return n.Addr
}

func (n *Node) ToBytes() []byte {
	buf, _ := json.Marshal(n)
	return buf
}
