package RouteData

import (
	"encoding/json"
	"hego/common/cfg"
	"sync/atomic"
)

var obj = atomic.Value{}

type RouteData struct {
	listData   []*cfg.RouteConfig
	mapData0   map[uint32]*cfg.RouteConfig
	groupData0 map[group0][]*cfg.RouteConfig
}

type group0 struct {
	ServerType cfg.ServerType
	RouteType  cfg.RouteType
}

func SGet() *cfg.RouteConfig {
	if d, ok := obj.Load().(*RouteData); ok {
		return d.listData[0]
	}
	return nil
}

func LGet() (rets []*cfg.RouteConfig) {
	if d, ok := obj.Load().(*RouteData); ok {
		rets = make([]*cfg.RouteConfig, len(d.listData))
		copy(rets, d.listData)
	}
	return
}

func MGet0(ID uint32) *cfg.RouteConfig {
	if d, ok := obj.Load().(*RouteData); ok {
		val, ok := d.mapData0[ID]
		if ok {
			return val
		}
	}
	return nil
}

func GGet0(ServerType cfg.ServerType, RouteType cfg.RouteType) (rets []*cfg.RouteConfig) {
	if d, ok := obj.Load().(*RouteData); ok {
		vals, ok := d.groupData0[group0{ServerType, RouteType}]
		if ok {
			rets = make([]*cfg.RouteConfig, len(vals))
			copy(rets, vals)
		}
	}
	return
}

func Parse(buf []byte) {
	ary := []*cfg.RouteConfig{}
	if err := json.Unmarshal(buf, &ary); err != nil {
		panic(err)
	}
	data := &RouteData{
		mapData0:   make(map[uint32]*cfg.RouteConfig),
		groupData0: make(map[group0][]*cfg.RouteConfig),
	}
	for _, item := range ary {
		data.listData = append(data.listData, item)
		data.mapData0[item.ID] = item
		key := group0{item.ServerType, item.RouteType}
		data.groupData0[key] = append(data.groupData0[key], item)
	}
	obj.Store(data)
}
