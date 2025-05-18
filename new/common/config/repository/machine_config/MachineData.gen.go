/*
* 本代码由xlsx工具生成，请勿手动修改
 */

package machine_config

import (
	"encoding/json"
	"sync/atomic"

	"poker_server/common/config"
	"poker_server/common/pb"
)

var obj = atomic.Value{}

type MachineConfigData struct {
	_List   []*pb.MachineConfig
	_GameId map[int32]*pb.MachineConfig
}

// 注册函数
func init() {
	config.Register("MachineConfig", parse)
}

func parse(buf string) error {
	data := &pb.MachineConfigAry{}
	if err := json.Unmarshal([]byte(buf), data); err != nil {
		return err
	}

	_GameId := make(map[int32]*pb.MachineConfig)
	for _, item := range data.Ary {
		_GameId[item.GameId] = item
	}

	obj.Store(&MachineConfigData{
		_List:   data.Ary,
		_GameId: _GameId,
	})
	return nil
}

func SGet() *pb.MachineConfig {
	obj, ok := obj.Load().(*MachineConfigData)
	if !ok {
		return nil
	}
	return obj._List[0]
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

func MGetGameId(GameId int32) *pb.MachineConfig {
	obj, ok := obj.Load().(*MachineConfigData)
	if !ok {
		return nil
	}

	if val, ok := obj._GameId[GameId]; ok {
		return val
	}
	return nil
}
