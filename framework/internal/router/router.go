package router

import (
	"universal/framework/domain"
	"universal/framework/global"
)

type Route struct {
	gate  int32 // 网关
	game  int32 // 游戏
	db    int32 // 数据服
	room  int32 // 房间服
	match int32 // 匹配服
}

func NewRouter() domain.IRouter {
	return &Route{}
}

func (d *Route) GetSize() int {
	return 4 * 5
}

func (d *Route) WriteTo(buf []byte) error {
	pos := 0
	buf[pos] = byte(d.gate)
	pos += 4
	buf[pos] = byte(d.game)
	pos += 4
	buf[pos] = byte(d.db)
	pos += 4
	buf[pos] = byte(d.room)
	pos += 4
	buf[pos] = byte(d.match)
	return nil
}

func (d *Route) ReadFrom(buf []byte) error {
	pos := 0
	d.gate = int32(buf[pos])
	pos += 4
	d.game = int32(buf[pos])
	pos += 4
	d.db = int32(buf[pos])
	pos += 4
	d.room = int32(buf[pos])
	pos += 4
	d.match = int32(buf[pos])
	return nil
}

func (d *Route) Get(nodeType int32) int32 {
	switch nodeType {
	case int32(global.NodeTypeGate):
		return d.gate
	case int32(global.NodeTypeGame):
		return d.game
	case int32(global.NodeTypeDb):
		return d.db
	case int32(global.NodeTypeRoom):
		return d.room
	case int32(global.NodeTypeMatch):
		return d.match
	}
	return d.gate
}

func (d *Route) Set(nodeType, nodeId int32) {
	switch nodeType {
	case int32(global.NodeTypeGate):
		d.gate = nodeId
	case int32(global.NodeTypeGame):
		d.game = nodeId
	case int32(global.NodeTypeDb):
		d.db = nodeId
	case int32(global.NodeTypeRoom):
		d.room = nodeId
	case int32(global.NodeTypeMatch):
		d.match = nodeId
	}
}
