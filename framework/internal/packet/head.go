package packet

import (
	"encoding/binary"
	"universal/framework/domain"
)

type Head struct {
	srcNodeType int32
	srcNodeId   int32
	dstNodeType int32
	dstNodeId   int32
	routeId     uint64
	uid         uint64
	actorName   string
	funcName    string
}

func (d *Head) GetSrcNodeType() int32 {
	return d.srcNodeType
}

func (d *Head) GetSrcNodeId() int32 {
	return d.srcNodeId
}

func (d *Head) GetDstNodeType() int32 {
	return d.dstNodeType
}

func (d *Head) GetDstNodeId() int32 {
	return d.dstNodeId
}

func (d *Head) GetRouteId() uint64 {
	return d.routeId
}

func (d *Head) GetUid() uint64 {
	return d.uid
}

func (d *Head) GetActorName() string {
	return d.actorName
}

func (d *Head) GetFuncName() string {
	return d.funcName
}

// --------设置方法------
func (d *Head) SetSrcNodeType(srcNodeType int32) domain.IHead {
	d.srcNodeType = srcNodeType
	return d
}

func (d *Head) SetSrcNodeId(srcNodeId int32) domain.IHead {
	d.srcNodeId = srcNodeId
	return d
}

func (d *Head) SetDstNodeType(dstNodeType int32) domain.IHead {
	d.dstNodeType = dstNodeType
	return d
}

func (d *Head) SetDstNodeId(dstNodeId int32) domain.IHead {
	d.dstNodeId = dstNodeId
	return d
}

func (d *Head) SetRouteId(routeId uint64) domain.IHead {
	d.routeId = routeId
	return d
}

func (d *Head) SetUid(uid uint64) domain.IHead {
	d.uid = uid
	return d
}

func (d *Head) SetActorName(actorName string) domain.IHead {
	d.actorName = actorName
	return d
}

func (d *Head) SetFuncName(funcName string) domain.IHead {
	d.funcName = funcName
	return d
}

func (d *Head) GetSize() int {
	return 4 + 4 + 4 + 4 + 8 + 8 + len(d.actorName) + len(d.funcName) + 4 + 4
}

func (d *Head) WriteTo(buf []byte) error {
	pos := 0
	binary.BigEndian.PutUint32(buf[pos:], uint32(d.srcNodeType))
	pos += 4
	binary.BigEndian.PutUint32(buf[pos:], uint32(d.srcNodeId))
	pos += 4
	binary.BigEndian.PutUint32(buf[pos:], uint32(d.dstNodeType))
	pos += 4
	binary.BigEndian.PutUint32(buf[pos:], uint32(d.dstNodeId))
	pos += 4
	binary.BigEndian.PutUint64(buf[pos:], d.routeId)
	pos += 8
	binary.BigEndian.PutUint64(buf[pos:], d.uid)
	pos += 8
	lactor := len(d.actorName)
	binary.BigEndian.PutUint32(buf[pos:], uint32(lactor))
	pos += 4
	copy(buf[pos:], d.actorName)
	pos += lactor
	lfunc := len(d.funcName)
	binary.BigEndian.PutUint32(buf[pos:], uint32(lfunc))
	pos += 4
	copy(buf[pos:], d.funcName)
	return nil
}

func (d *Head) ReadFrom(buf []byte) error {
	pos := 0
	d.srcNodeType = int32(binary.BigEndian.Uint32(buf[pos:]))
	pos += 4
	d.srcNodeId = int32(binary.BigEndian.Uint32(buf[pos:]))
	pos += 4
	d.dstNodeType = int32(binary.BigEndian.Uint32(buf[pos:]))
	pos += 4
	d.dstNodeId = int32(binary.BigEndian.Uint32(buf[pos:]))
	pos += 4
	d.routeId = binary.BigEndian.Uint64(buf[pos:])
	pos += 8
	d.uid = binary.BigEndian.Uint64(buf[pos:])
	pos += 8
	lactor := int(binary.BigEndian.Uint32(buf[pos:]))
	pos += 4
	d.actorName = string(buf[pos : pos+lactor])
	pos += lactor
	lfunc := int(binary.BigEndian.Uint32(buf[pos:]))
	pos += 4
	d.funcName = string(buf[pos : pos+lfunc])
	return nil
}
