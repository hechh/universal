/*
* 本代码由xlsx工具生成，请勿手动修改
 */

package open_api_config

import (
	"encoding/json"
	"sync/atomic"

	"poker_server/common/config"
	"poker_server/common/pb"
)

var obj = atomic.Value{}

type OpenApiConfigData struct {
	_List                        []*pb.OpenApiConfig
	_Cmd                         map[uint32]*pb.OpenApiConfig
	_ServerTypeActorNameFuncName map[pb.Index3[pb.ServerType, string, string]]*pb.OpenApiConfig
}

// 注册函数
func init() {
	config.Register("OpenApiConfig", parse)
}

func parse(buf string) error {
	data := &pb.OpenApiConfigAry{}
	if err := json.Unmarshal([]byte(buf), data); err != nil {
		return err
	}

	_Cmd := make(map[uint32]*pb.OpenApiConfig)
	_ServerTypeActorNameFuncName := make(map[pb.Index3[pb.ServerType, string, string]]*pb.OpenApiConfig)
	for _, item := range data.Ary {
		_Cmd[item.Cmd] = item
		_ServerTypeActorNameFuncName[pb.Index3[pb.ServerType, string, string]{item.ServerType, item.ActorName, item.FuncName}] = item
	}

	obj.Store(&OpenApiConfigData{
		_List:                        data.Ary,
		_Cmd:                         _Cmd,
		_ServerTypeActorNameFuncName: _ServerTypeActorNameFuncName,
	})
	return nil
}

func SGet() *pb.OpenApiConfig {
	obj, ok := obj.Load().(*OpenApiConfigData)
	if !ok {
		return nil
	}
	return obj._List[0]
}

func LGet() (rets []*pb.OpenApiConfig) {
	obj, ok := obj.Load().(*OpenApiConfigData)
	if !ok {
		return
	}
	rets = make([]*pb.OpenApiConfig, len(obj._List))
	copy(rets, obj._List)
	return
}

func Walk(f func(*pb.OpenApiConfig) bool) {
	obj, ok := obj.Load().(*OpenApiConfigData)
	if !ok {
		return
	}
	for _, item := range obj._List {
		if !f(item) {
			return
		}
	}
}

func MGetCmd(Cmd uint32) *pb.OpenApiConfig {
	obj, ok := obj.Load().(*OpenApiConfigData)
	if !ok {
		return nil
	}

	if val, ok := obj._Cmd[Cmd]; ok {
		return val
	}
	return nil
}

func MGetServerTypeActorNameFuncName(ServerType pb.ServerType, ActorName string, FuncName string) *pb.OpenApiConfig {
	obj, ok := obj.Load().(*OpenApiConfigData)
	if !ok {
		return nil
	}

	if val, ok := obj._ServerTypeActorNameFuncName[pb.Index3[pb.ServerType, string, string]{ServerType, ActorName, FuncName}]; ok {
		return val
	}
	return nil
}
