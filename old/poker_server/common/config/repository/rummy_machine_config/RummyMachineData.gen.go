/*
* 本代码由cfgtool工具生成，请勿手动修改
 */

package rummy_machine_config

import (
	"poker_server/common/config"
	"poker_server/common/pb"
	"sync/atomic"

	"github.com/golang/protobuf/proto"
)

var obj = atomic.Value{}

type RummyMachineConfigData struct {
	_List     []*pb.RummyMachineConfig
	_GameType map[pb.GameType]*pb.RummyMachineConfig
}

// 注册函数
func init() {
	config.Register("RummyMachineConfig", parse)
}

func parse(buf string) error {
	data := &pb.RummyMachineConfigAry{}
	if err := proto.UnmarshalText(buf, data); err != nil {
		return err
	}

	_GameType := make(map[pb.GameType]*pb.RummyMachineConfig)
	for _, item := range data.Ary {
		_GameType[item.GameType] = item
	}

	obj.Store(&RummyMachineConfigData{
		_List:     data.Ary,
		_GameType: _GameType,
	})
	return nil
}

func SGet() *pb.RummyMachineConfig {
	obj, ok := obj.Load().(*RummyMachineConfigData)
	if !ok {
		return nil
	}
	return obj._List[len(obj._List)-1]
}

func LGet() (rets []*pb.RummyMachineConfig) {
	obj, ok := obj.Load().(*RummyMachineConfigData)
	if !ok {
		return
	}
	rets = make([]*pb.RummyMachineConfig, len(obj._List))
	copy(rets, obj._List)
	return
}

func Walk(f func(*pb.RummyMachineConfig) bool) {
	obj, ok := obj.Load().(*RummyMachineConfigData)
	if !ok {
		return
	}
	for _, item := range obj._List {
		if !f(item) {
			return
		}
	}
}

func MGetGameType(GameType pb.GameType) *pb.RummyMachineConfig {
	obj, ok := obj.Load().(*RummyMachineConfigData)
	if !ok {
		return nil
	}
	if val, ok := obj._GameType[GameType]; ok {
		return val
	}
	return nil
}
