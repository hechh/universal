/*
* 本代码由cfgtool工具生成，请勿手动修改
 */

package texas_test_config

import (
	"poker_server/common/config"
	"poker_server/common/pb"
	"sync/atomic"

	"github.com/golang/protobuf/proto"
)

var obj = atomic.Value{}

type TexasTestConfigData struct {
	_List     []*pb.TexasTestConfig
	_MaxRound uint32
	_Round    map[uint32]*pb.TexasTestConfig
}

// 注册函数
func init() {
	config.Register("TexasTestConfig", parse)
}

func parse(buf string) error {
	data := &pb.TexasTestConfigAry{}
	if err := proto.UnmarshalText(buf, data); err != nil {
		return err
	}

	var _MaxRound uint32
	_Round := make(map[uint32]*pb.TexasTestConfig)
	for _, item := range data.Ary {
		if _MaxRound < item.Round {
			_MaxRound = item.Round
		}
		_Round[item.Round] = item
	}

	obj.Store(&TexasTestConfigData{
		_List:     data.Ary,
		_MaxRound: _MaxRound,
		_Round:    _Round,
	})
	return nil
}

func SGet() *pb.TexasTestConfig {
	obj, ok := obj.Load().(*TexasTestConfigData)
	if !ok {
		return nil
	}
	return obj._List[len(obj._List)-1]
}

func LGet() (rets []*pb.TexasTestConfig) {
	obj, ok := obj.Load().(*TexasTestConfigData)
	if !ok {
		return
	}
	rets = make([]*pb.TexasTestConfig, len(obj._List))
	copy(rets, obj._List)
	return
}

func Walk(f func(*pb.TexasTestConfig) bool) {
	obj, ok := obj.Load().(*TexasTestConfigData)
	if !ok {
		return
	}
	for _, item := range obj._List {
		if !f(item) {
			return
		}
	}
}

func MGetRoundKey(val uint32) uint32 {
	if obj, ok := obj.Load().(*TexasTestConfigData); ok && val > obj._MaxRound {
		return obj._MaxRound
	}
	return val
}

func MGetRound(Round uint32) *pb.TexasTestConfig {
	obj, ok := obj.Load().(*TexasTestConfigData)
	if !ok {
		return nil
	}
	if val, ok := obj._Round[Round]; ok {
		return val
	}
	return nil
}
