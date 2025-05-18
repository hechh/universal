/*
* 本代码由xlsx工具生成，请勿手动修改
 */

package texas_test_config

import (
	"encoding/json"
	"sync/atomic"

	"poker_server/common/config"
	"poker_server/common/pb"
)

var obj = atomic.Value{}

type TexasTestConfigData struct {
	_List  []*pb.TexasTestConfig
	_Round map[uint32]*pb.TexasTestConfig
}

// 注册函数
func init() {
	config.Register("TexasTestConfig", parse)
}

func parse(buf string) error {
	data := &pb.TexasTestConfigAry{}
	if err := json.Unmarshal([]byte(buf), data); err != nil {
		return err
	}

	_Round := make(map[uint32]*pb.TexasTestConfig)
	for _, item := range data.Ary {
		_Round[item.Round] = item
	}

	obj.Store(&TexasTestConfigData{
		_List:  data.Ary,
		_Round: _Round,
	})
	return nil
}

func SGet() *pb.TexasTestConfig {
	obj, ok := obj.Load().(*TexasTestConfigData)
	if !ok {
		return nil
	}
	return obj._List[0]
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
