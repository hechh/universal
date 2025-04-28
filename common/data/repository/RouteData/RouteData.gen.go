package RouteData

import (
	"encoding/json"
	"sync/atomic"
	"universal/common/config/cfg"
	"universal/common/config/internal/manager"
	"universal/library/baselib/uerror"
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

func Parse(buf []byte) error {
	ary := []*cfg.RouteConfig{}
	if err := json.Unmarshal(buf, &ary); err != nil {
		return uerror.New(1, -1, err.Error())
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
	return nil
}

func init() {
	manager.Register("route", Parse)
}
