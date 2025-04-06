package RouteData

import (
	"encoding/json"
	"hego/common/cfg"
	"sync/atomic"
)

var obj = atomic.Value{}

type ServerTypeRouteType struct {
	ServerType cfg.ServerType
	RouteType  cfg.RouteType
}

type RouteData struct {
	_list                []*cfg.RouteConfig
	_ID                  map[uint32]*cfg.RouteConfig
	_ServerTypeRouteType map[ServerTypeRouteType][]*cfg.RouteConfig
}

func SGet() *cfg.RouteConfig {
	if d, ok := obj.Load().(*RouteData); ok {
		return d._list[0]
	}
	return nil
}

func LGet() (rets []*cfg.RouteConfig) {
	if d, ok := obj.Load().(*RouteData); ok {
		rets = make([]*cfg.RouteConfig, len(d._list))
		copy(rets, d._list)
	}
	return
}

func MGetID(ID uint32) *cfg.RouteConfig {
	if d, ok := obj.Load().(*RouteData); ok {
		val, ok := d._ID[ID]
		if ok {
			return val
		}
	}
	return nil
}

func GGetServerTypeRouteType(ServerType cfg.ServerType, RouteType cfg.RouteType) (rets []*cfg.RouteConfig) {
	if d, ok := obj.Load().(*RouteData); ok {
		vals, ok := d._ServerTypeRouteType[ServerTypeRouteType{ServerType, RouteType}]
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
	_ID := make(map[uint32]*cfg.RouteConfig)
	_ServerTypeRouteType := make(map[ServerTypeRouteType][]*cfg.RouteConfig)
	for _, item := range ary {
		_ID[item.ID] = item
		_ServerTypeRouteType[ServerTypeRouteType{item.ServerType, item.RouteType}] = append(_ServerTypeRouteType[ServerTypeRouteType{item.ServerType, item.RouteType}], item)
	}
	obj.Store(&RouteData{
		_list:                ary,
		_ID:                  _ID,
		_ServerTypeRouteType: _ServerTypeRouteType,
	})
}
