/*
* 本代码由xlsx工具生成，请勿手动修改
 */

package texas_config

import (
	"encoding/json"
	"sync/atomic"

	"poker_server/common/config"
	"poker_server/common/pb"
)

var obj = atomic.Value{}

type TexasConfigData struct {
	_List              []*pb.TexasConfig
	_RoomStageCoinType map[pb.Index2[int32, int32]]*pb.TexasConfig
}

// 注册函数
func init() {
	config.Register("TexasConfig", parse)
}

func parse(buf string) error {
	data := &pb.TexasConfigAry{}
	if err := json.Unmarshal([]byte(buf), data); err != nil {
		return err
	}

	_RoomStageCoinType := make(map[pb.Index2[int32, int32]]*pb.TexasConfig)
	for _, item := range data.Ary {
		_RoomStageCoinType[pb.Index2[int32, int32]{item.RoomStage, item.CoinType}] = item
	}

	obj.Store(&TexasConfigData{
		_List:              data.Ary,
		_RoomStageCoinType: _RoomStageCoinType,
	})
	return nil
}

func SGet() *pb.TexasConfig {
	obj, ok := obj.Load().(*TexasConfigData)
	if !ok {
		return nil
	}
	return obj._List[0]
}

func LGet() (rets []*pb.TexasConfig) {
	obj, ok := obj.Load().(*TexasConfigData)
	if !ok {
		return
	}
	rets = make([]*pb.TexasConfig, len(obj._List))
	copy(rets, obj._List)
	return
}

func Walk(f func(*pb.TexasConfig) bool) {
	obj, ok := obj.Load().(*TexasConfigData)
	if !ok {
		return
	}
	for _, item := range obj._List {
		if !f(item) {
			return
		}
	}
}

func MGetRoomStageCoinType(RoomStage int32, CoinType int32) *pb.TexasConfig {
	obj, ok := obj.Load().(*TexasConfigData)
	if !ok {
		return nil
	}

	if val, ok := obj._RoomStageCoinType[pb.Index2[int32, int32]{RoomStage, CoinType}]; ok {
		return val
	}
	return nil
}
