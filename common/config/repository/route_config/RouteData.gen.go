/*
* 本代码由cfgtool工具生成，请勿手动修改
 */

package route_config

import (
	"sync/atomic"
	"universal/common/config"
	"universal/common/pb"

	"github.com/golang/protobuf/proto"
)

var obj = atomic.Value{}

type RouteConfigData struct {
	_List                        []*pb.RouteConfig
	_Cmd                         map[uint32]*pb.RouteConfig
	_ServerTypeActorNameFuncName map[pb.Index3[pb.ServerType, string, string]]*pb.RouteConfig
}

// 注册函数
func init() {
	config.Register("RouteConfig", parse)
}

func parse(buf string) error {
	data := &pb.RouteConfigAry{}
	if err := proto.UnmarshalText(buf, data); err != nil {
		return err
	}

	_Cmd := make(map[uint32]*pb.RouteConfig)
	_ServerTypeActorNameFuncName := make(map[pb.Index3[pb.ServerType, string, string]]*pb.RouteConfig)
	for _, item := range data.Ary {
		// map数据
		_Cmd[item.Cmd] = item
		// map数据
		keyServerTypeActorNameFuncName := pb.Index3[pb.ServerType, string, string]{item.ServerType, item.ActorName, item.FuncName}
		_ServerTypeActorNameFuncName[keyServerTypeActorNameFuncName] = item
	}

	obj.Store(&RouteConfigData{
		_List:                        data.Ary,
		_Cmd:                         _Cmd,
		_ServerTypeActorNameFuncName: _ServerTypeActorNameFuncName,
	})
	return nil
}

func SGet() *pb.RouteConfig {
	obj, ok := obj.Load().(*RouteConfigData)
	if !ok {
		return nil
	}
	return obj._List[0]
}

func LGet() (rets []*pb.RouteConfig) {
	obj, ok := obj.Load().(*RouteConfigData)
	if !ok {
		return
	}
	rets = make([]*pb.RouteConfig, len(obj._List))
	copy(rets, obj._List)
	return
}

func Walk(f func(*pb.RouteConfig) bool) {
	obj, ok := obj.Load().(*RouteConfigData)
	if !ok {
		return
	}
	for _, item := range obj._List {
		if !f(item) {
			return
		}
	}
}

func MGetCmd(Cmd uint32) *pb.RouteConfig {
	obj, ok := obj.Load().(*RouteConfigData)
	if !ok {
		return nil
	}

	if val, ok := obj._Cmd[Cmd]; ok {
		return val
	}
	return nil
}

func MGetServerTypeActorNameFuncName(ServerType pb.ServerType, ActorName string, FuncName string) *pb.RouteConfig {
	obj, ok := obj.Load().(*RouteConfigData)
	if !ok {
		return nil
	}

	if val, ok := obj._ServerTypeActorNameFuncName[pb.Index3[pb.ServerType, string, string]{ServerType, ActorName, FuncName}]; ok {
		return val
	}
	return nil
}
