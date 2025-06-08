/*
* 本代码由cfgtool工具生成，请勿手动修改
 */

package router_config

import (
	"sync/atomic"
	"universal/common/config"
	"universal/common/pb"

	"github.com/golang/protobuf/proto"
)

var obj atomic.Pointer[RouterConfigData]

type RouterConfigData struct {
	_List                      []*pb.RouterConfig
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

	_Cmd := make(map[uint32]*pb.RouterConfig)
	_NodeTypeActorNameFuncName := make(map[pb.Index3[pb.NodeType, string, string]]*pb.RouterConfig)
	for _, item := range data.Ary {
		// map数据
		_Cmd[item.Cmd] = item
		// map数据
		keyNodeTypeActorNameFuncName := pb.Index3[pb.NodeType, string, string]{item.NodeType, item.ActorName, item.FuncName}
		_NodeTypeActorNameFuncName[keyNodeTypeActorNameFuncName] = item
	}

	obj.Store(&RouterConfigData{
		_List:                      data.Ary,
		_Cmd:                       _Cmd,
		_NodeTypeActorNameFuncName: _NodeTypeActorNameFuncName,
	})
	return nil
}

func SGet() *pb.RouterConfig {
	if obj := obj.Load(); obj != nil {
		return obj._List[0]
	}
	return nil
}

func LGet() (rets []*pb.RouterConfig) {
	if obj := obj.Load(); obj != nil {
		rets = make([]*pb.RouterConfig, len(obj._List))
		copy(rets, obj._List)
	}
	return
}

func Walk(f func(*pb.RouterConfig) bool) {
	if obj := obj.Load(); obj != nil {
		for _, item := range obj._List {
			if !f(item) {
				return
			}
		}
	}
	return
}

func MGetCmd(Cmd uint32) *pb.RouterConfig {
	if obj := obj.Load(); obj != nil {
		if val, ok := obj._Cmd[Cmd]; ok {
			return val
		}
	}
	return nil
}

func MGetNodeTypeActorNameFuncName(NodeType pb.NodeType, ActorName string, FuncName string) *pb.RouterConfig {
	if obj := obj.Load(); obj != nil {
		if val, ok := obj._NodeTypeActorNameFuncName[pb.Index3[pb.NodeType, string, string]{NodeType, ActorName, FuncName}]; ok {
			return val
		}
	}
	return nil
}
