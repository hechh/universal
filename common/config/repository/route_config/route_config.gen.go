package route_config

import (
	"hego/common/config/internal/manager"
	"hego/common/pb"
	"hego/framework/uerror"
	"sync/atomic"

	"github.com/golang/protobuf/proto"
)

var (
	data atomic.Value
)

type map0 struct {
	ID uint32
}

type group0 struct {
	ID         uint32
	ServerType pb.SERVER
}

type RouteConfigData struct {
	cfg    *pb.RouteConfig
	cfgs   []*pb.RouteConfig
	map0   map[map0]*pb.RouteConfig
	group0 map[group0][]*pb.RouteConfig
}

func init() {
	manager.Register("RouteConfig", load)
}

func load(buf []byte) error {
	ary := &pb.RouteConfigAry{}
	if err := proto.Unmarshal(buf, ary); err != nil {
		return uerror.NewUError(1, -1, "加载RouteConfig配置失败: %v", err)
	}

	dd := &RouteConfigData{
		cfg:    ary.Ary[0],
		map0:   make(map[map0]*pb.RouteConfig),
		group0: make(map[group0][]*pb.RouteConfig),
	}
	for _, item := range ary.Ary {
		dd.cfgs = append(dd.cfgs, item)
		dd.map0[map0{item.ID}] = item
		dd.group0[group0{item.ID, item.ServerType}] = append(dd.group0[group0{item.ID, item.ServerType}], item)

	}
	data.Store(dd)
	return nil
}

func getObj() *RouteConfigData {
	if obj, ok := data.Load().(*RouteConfigData); ok && obj != nil {
		return obj
	}
	return nil
}

func Get() *pb.RouteConfig {
	return getObj().cfg
}

func Gets() (rets []*pb.RouteConfig) {
	list := getObj().cfgs
	rets = make([]*pb.RouteConfig, len(list))
	copy(rets, list)
	return
}

func GetByID(ID uint32) *pb.RouteConfig {
	return getObj().map0[map0{ID}]
}

func GetsByIDServerType(ID uint32, ServerType pb.SERVER) (rets []*pb.RouteConfig) {
	list := getObj().group0[group0{ID, ServerType}]
	if len(list) > 0 {
		rets = make([]*pb.RouteConfig, len(list))
		copy(rets, list)
	}
	return
}
