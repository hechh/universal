package route

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

type RouteCfg struct {
	cfgs map[uint32]*pb.RouteConfig
}

func init() {
	manager.Register("RouteConfig", loadRouteConfig)
}

func loadRouteConfig(buf []byte) error {
	// 序列化
	list := &pb.RouteConfigAry{}
	if err := proto.Unmarshal(buf, list); err != nil {
		return uerror.NewUError(1, -1, "RouteConfig加载失败: %v", err)
	}
	// 加载配置数据
	cfgs := make(map[uint32]*pb.RouteConfig)
	for _, item := range list.Ary {
		cfgs[item.ID] = item
	}
	// 替换数据
	data.Store(&RouteCfg{cfgs: cfgs})
	return nil
}

func Get(apiID uint32) *pb.RouteConfig {
	if obj, ok := data.Load().(*RouteCfg); ok && obj != nil {
		return obj.cfgs[apiID]
	}
	return nil
}
