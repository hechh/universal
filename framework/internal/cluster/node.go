package cluster

import (
	"encoding/binary"
	"encoding/json"
	"universal/framework/domain"
)

type Node struct {
	Name string `json:"name"` // 节点名称
	Addr string `json:"addr"` // 节点地址
	Type int32  `json:"type"` // 节点类型
	Id   int32  `json:"id"`   // 节点ID
}

func NewNode() domain.INode {
	return &Node{}
}

func (n *Node) GetAddr() string {
	return n.Addr
}

func (n *Node) GetType() int32 {
	return n.Type
}

func (n *Node) GetId() int32 {
	return n.Id
}

func (n *Node) GetName() string {
	return n.Name
}

func (n *Node) SetName(name string) domain.INode {
	n.Name = name
	return n
}

func (n *Node) SetAddr(addr string) domain.INode {
	n.Addr = addr
	return n
}

func (n *Node) SetType(t int32) domain.INode {
	n.Type = t
	return n
}

func (n *Node) SetId(id int32) domain.INode {
	n.Id = id
	return n
}

func (n *Node) String() string {
	buf, _ := json.Marshal(n)
	return string(buf)
}

func (n *Node) GetSize() int {
	return len(n.Addr) + 12
}

func (n *Node) WriteTo(buf []byte) error {
	pos := 0
	binary.BigEndian.PutUint32(buf[pos:], uint32(n.Type))
	pos += 4
	binary.BigEndian.PutUint32(buf[pos:], uint32(n.Id))
	pos += 4
	laddr := len(n.Addr)
	binary.BigEndian.PutUint32(buf[pos:], uint32(laddr))
	pos += 4
	copy(buf[pos:], []byte(n.Addr))
	pos += laddr
	return nil
}

func (n *Node) ReadFrom(buf []byte) error {
	pos := 0
	n.Type = int32(binary.BigEndian.Uint32(buf[pos:]))
	pos += 4
	n.Id = int32(binary.BigEndian.Uint32(buf[pos:]))
	pos += 4
	laddr := int(binary.BigEndian.Uint32(buf[pos:]))
	pos += 4
	n.Addr = string(buf[pos : pos+laddr])
	return nil
}
