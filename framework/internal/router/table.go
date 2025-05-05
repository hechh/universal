package router

import (
	"encoding/binary"
	"universal/framework/define"
)

type Table struct {
	Gate uint32
	Db   uint32
	Game uint32
	Tool uint32
	Rank uint32
}

func NewTable() define.ITable {
	return &Table{}
}

func (r *Table) GetSize() int {
	return 20
}

func (r *Table) ToBytes(buf []byte) {
	pos := 0
	binary.BigEndian.PutUint32(buf[pos:], r.Gate)
	pos += 4
	binary.BigEndian.PutUint32(buf[pos:], r.Db)
	pos += 4
	binary.BigEndian.PutUint32(buf[pos:], r.Game)
	pos += 4
	binary.BigEndian.PutUint32(buf[pos:], r.Tool)
	pos += 4
	binary.BigEndian.PutUint32(buf[pos:], r.Rank)
}

func (r *Table) Parse(buf []byte) define.ITable {
	pos := 0
	r.Gate = binary.BigEndian.Uint32(buf[pos:])
	pos += 4
	r.Db = binary.BigEndian.Uint32(buf[pos:])
	pos += 4
	r.Game = binary.BigEndian.Uint32(buf[pos:])
	pos += 4
	r.Tool = binary.BigEndian.Uint32(buf[pos:])
	pos += 4
	r.Rank = binary.BigEndian.Uint32(buf[pos:])
	return r
}

func (r *Table) Get(nodeType uint32) uint32 {
	switch nodeType {
	case uint32(define.NodeTypeGate):
		return r.Gate
	case uint32(define.NodeTypeDb):
		return r.Db
	case uint32(define.NodeTypeGame):
		return r.Game
	case uint32(define.NodeTypeTool):
		return r.Tool
	case uint32(define.NodeTypeRank):
		return r.Rank
	}
	return 0
}

func (r *Table) Set(nodeType, nodeId uint32) {
	switch nodeType {
	case uint32(define.NodeTypeGate):
		r.Gate = nodeId
	case uint32(define.NodeTypeDb):
		r.Db = nodeId
	case uint32(define.NodeTypeGame):
		r.Game = nodeId
	case uint32(define.NodeTypeTool):
		r.Tool = nodeId
	case uint32(define.NodeTypeRank):
		r.Rank = nodeId
	}
}
