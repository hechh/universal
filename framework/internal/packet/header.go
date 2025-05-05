package packet

import (
	"encoding/binary"
	"universal/framework/define"
)

type Header struct {
	SrcNodeType uint32
	SrcNodeId   uint32
	DstNodeType uint32
	DstNodeId   uint32
	Cmd         uint32
	Uid         uint64
	RouteId     uint64
	ActorName   string
	FuncName    string
	Table       *define.RouteInfo
}

func NewHeader() define.IHeader {
	return &Header{Table: &define.RouteInfo{}}
}

func (h *Header) GetSrcNodeType() uint32 {
	return h.SrcNodeType
}

func (h *Header) GetSrcNodeId() uint32 {
	return h.SrcNodeId
}

func (h *Header) GetDstNodeType() uint32 {
	return h.DstNodeType
}

func (h *Header) GetDstNodeId() uint32 {
	return h.DstNodeId
}

func (h *Header) GetCmd() uint32 {
	return h.Cmd
}

func (h *Header) GetUid() uint64 {
	return h.Uid
}

func (h *Header) GetRouteId() uint64 {
	return h.RouteId
}

func (h *Header) GetActorName() string {
	return h.ActorName
}

func (h *Header) GetFuncName() string {
	return h.FuncName
}

func (h *Header) GetTable() *define.RouteInfo {
	return h.Table
}

func (h *Header) GetSize() int {
	return 44 + len(h.ActorName) + len(h.FuncName) + 20
}

func (h *Header) ToBytes(buf []byte) []byte {
	pos := 0
	binary.BigEndian.PutUint32(buf[pos:], uint32(h.SrcNodeType))
	pos += 4
	binary.BigEndian.PutUint32(buf[pos:], uint32(h.SrcNodeId))
	pos += 4
	binary.BigEndian.PutUint32(buf[pos:], uint32(h.DstNodeType))
	pos += 4
	binary.BigEndian.PutUint32(buf[pos:], uint32(h.DstNodeId))
	pos += 4
	binary.BigEndian.PutUint32(buf[pos:], h.Cmd)
	pos += 4
	binary.BigEndian.PutUint64(buf[pos:], h.Uid)
	pos += 8
	binary.BigEndian.PutUint64(buf[pos:], h.RouteId)
	pos += 8
	lactor := len(h.ActorName)
	binary.BigEndian.PutUint32(buf[pos:], uint32(lactor))
	pos += 4
	copy(buf[pos:], []byte(h.ActorName))
	pos += lactor
	lfunc := len(h.FuncName)
	binary.BigEndian.PutUint32(buf[pos:], uint32(lfunc))
	pos += 4
	copy(buf[pos:], []byte(h.FuncName))
	pos += lfunc
	if h.Table == nil {
		h.Table = &define.RouteInfo{}
	}
	binary.BigEndian.PutUint32(buf[pos:], h.Table.Gate)
	pos += 4
	binary.BigEndian.PutUint32(buf[pos:], h.Table.Db)
	pos += 4
	binary.BigEndian.PutUint32(buf[pos:], h.Table.Game)
	pos += 4
	binary.BigEndian.PutUint32(buf[pos:], h.Table.Tool)
	pos += 4
	binary.BigEndian.PutUint32(buf[pos:], h.Table.Rank)
	return buf
}

func (h *Header) Parse(buf []byte) define.IHeader {
	pos := 0
	h.SrcNodeType = binary.BigEndian.Uint32(buf[pos:])
	pos += 4
	h.SrcNodeId = binary.BigEndian.Uint32(buf[pos:])
	pos += 4
	h.DstNodeType = binary.BigEndian.Uint32(buf[pos:])
	pos += 4
	h.DstNodeId = binary.BigEndian.Uint32(buf[pos:])
	pos += 4
	h.Cmd = binary.BigEndian.Uint32(buf[pos:])
	pos += 4
	h.Uid = binary.BigEndian.Uint64(buf[pos:])
	pos += 8
	h.RouteId = binary.BigEndian.Uint64(buf[pos:])
	pos += 8
	lactor := int(binary.BigEndian.Uint32(buf[pos:]))
	pos += 4
	h.ActorName = string(buf[pos : pos+lactor])
	pos += lactor
	lfunc := int(binary.BigEndian.Uint32(buf[pos:]))
	pos += 4
	h.FuncName = string(buf[pos : pos+lfunc])
	pos += lfunc
	h.Table = &define.RouteInfo{}
	h.Table.Gate = binary.BigEndian.Uint32(buf[pos:])
	pos += 4
	h.Table.Db = binary.BigEndian.Uint32(buf[pos:])
	pos += 4
	h.Table.Game = binary.BigEndian.Uint32(buf[pos:])
	pos += 4
	h.Table.Tool = binary.BigEndian.Uint32(buf[pos:])
	pos += 4
	h.Table.Rank = binary.BigEndian.Uint32(buf[pos:])
	return h
}

func (h *Header) SetCmd(cmd uint32) define.IHeader {
	h.Cmd = cmd
	return h
}

func (h *Header) SetUid(uid uint64) define.IHeader {
	h.Uid = uid
	return h
}

func (h *Header) SetRouteId(routeId uint64) define.IHeader {
	h.RouteId = routeId
	return h
}

func (d *Header) SetSrcNode(node define.INode) define.IHeader {
	d.SrcNodeType = node.GetType()
	d.SrcNodeId = node.GetId()
	return d
}

func (h *Header) SetDstNode(node define.INode) define.IHeader {
	h.DstNodeType = node.GetType()
	h.DstNodeId = node.GetId()
	return h
}

func (h *Header) SetActorName(name string) define.IHeader {
	h.ActorName = name
	return h
}

func (h *Header) SetFuncName(name string) define.IHeader {
	h.FuncName = name
	return h
}

func (h *Header) SetTable(tab *define.RouteInfo) define.IHeader {
	h.Table = tab
	return h
}
