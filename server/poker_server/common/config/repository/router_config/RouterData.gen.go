/*
* 本代码由cfgtool工具生成，请勿手动修改
 */

package router_config

import (
	"poker_server/common/config"
	"poker_server/common/pb"
	"sync/atomic"

	"github.com/golang/protobuf/proto"
)

var obj = atomic.Value{}

type RouterConfigData struct {
	_List                      []*pb.RouterConfig
	_MaxCmd                    uint32
	_Cmd                       map[uint32]*pb.RouterConfig
	_NodeTypeActorNameFuncName map[pb.Index3[pb.NodeType, string, string]]*pb.RouterConfig
}

// 注册函数
func init() {
	config.Register("RouterConfig", parse)
}

func parse(buf string) error {
	data := &pb.RouterConfigAry{}
	if err := proto.UnmarshalText(buf, data); err != nil {
		return err
	}

	var _MaxCmd uint32
	_Cmd := make(map[uint32]*pb.RouterConfig)
	_NodeTypeActorNameFuncName := make(map[pb.Index3[pb.NodeType, string, string]]*pb.RouterConfig)
	for _, item := range data.Ary {
		if _MaxCmd < item.Cmd {
			_MaxCmd = item.Cmd
		}
		_Cmd[item.Cmd] = item
		keyNodeTypeActorNameFuncName := pb.Index3[pb.NodeType, string, string]{item.NodeType, item.ActorName, item.FuncName}
		_NodeTypeActorNameFuncName[keyNodeTypeActorNameFuncName] = item
	}

	obj.Store(&RouterConfigData{
		_List:                      data.Ary,
		_MaxCmd:                    _MaxCmd,
		_Cmd:                       _Cmd,
		_NodeTypeActorNameFuncName: _NodeTypeActorNameFuncName,
	})
	return nil
}

func SGet() *pb.RouterConfig {
	obj, ok := obj.Load().(*RouterConfigData)
	if !ok {
		return nil
	}
	return obj._List[len(obj._List)-1]
}

func LGet() (rets []*pb.RouterConfig) {
	obj, ok := obj.Load().(*RouterConfigData)
	if !ok {
		return
	}
	rets = make([]*pb.RouterConfig, len(obj._List))
	copy(rets, obj._List)
	return
}

func Walk(f func(*pb.RouterConfig) bool) {
	obj, ok := obj.Load().(*RouterConfigData)
	if !ok {
		return
	}
	for _, item := range obj._List {
		if !f(item) {
			return
		}
	}
}

func MGetCmdKey(val uint32) uint32 {
	if obj, ok := obj.Load().(*RouterConfigData); ok && val > obj._MaxCmd {
		return obj._MaxCmd
	}
	return val
}

func MGetCmd(Cmd uint32) *pb.RouterConfig {
	obj, ok := obj.Load().(*RouterConfigData)
	if !ok {
		return nil
	}
	if val, ok := obj._Cmd[Cmd]; ok {
		return val
	}
	return nil
}

func MGetNodeTypeActorNameFuncName(NodeType pb.NodeType, ActorName string, FuncName string) *pb.RouterConfig {
	obj, ok := obj.Load().(*RouterConfigData)
	if !ok {
		return nil
	}
	if val, ok := obj._NodeTypeActorNameFuncName[pb.Index3[pb.NodeType, string, string]{NodeType, ActorName, FuncName}]; ok {
		return val
	}
	return nil
}
