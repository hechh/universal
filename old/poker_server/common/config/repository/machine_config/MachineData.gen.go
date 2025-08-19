/*
* 本代码由cfgtool工具生成，请勿手动修改
 */

package machine_config

import (
	"poker_server/common/config"
	"poker_server/common/pb"
	"sync/atomic"

	"github.com/golang/protobuf/proto"
)

var obj = atomic.Value{}

type MachineConfigData struct {
	_List     []*pb.MachineConfig
	_GameType map[pb.GameType]*pb.MachineConfig
}

// 注册函数
func init() {
	config.Register("MachineConfig", parse)
}

func parse(buf string) error {
	data := &pb.MachineConfigAry{}
	if err := proto.UnmarshalText(buf, data); err != nil {
		return err
	}

	_GameType := make(map[pb.GameType]*pb.MachineConfig)
	for _, item := range data.Ary {
		_GameType[item.GameType] = item
	}

	obj.Store(&MachineConfigData{
		_List:     data.Ary,
		_GameType: _GameType,
	})
	return nil
}

func SGet() *pb.MachineConfig {
	obj, ok := obj.Load().(*MachineConfigData)
	if !ok {
		return nil
	}
	return obj._List[len(obj._List)-1]
}

func LGet() (rets []*pb.MachineConfig) {
	obj, ok := obj.Load().(*MachineConfigData)
	if !ok {
		return
	}
	rets = make([]*pb.MachineConfig, len(obj._List))
	copy(rets, obj._List)
	return
}

func Walk(f func(*pb.MachineConfig) bool) {
	obj, ok := obj.Load().(*MachineConfigData)
	if !ok {
		return
	}
	for _, item := range obj._List {
		if !f(item) {
			return
		}
	}
}

func MGetGameType(GameType pb.GameType) *pb.MachineConfig {
	obj, ok := obj.Load().(*MachineConfigData)
	if !ok {
		return nil
	}
	if val, ok := obj._GameType[GameType]; ok {
		return val
	}
	return nil
}
