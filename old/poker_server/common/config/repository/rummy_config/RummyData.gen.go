/*
* 本代码由cfgtool工具生成，请勿手动修改
 */

package rummy_config

import (
	"poker_server/common/config"
	"poker_server/common/pb"
	"sync/atomic"

	"github.com/golang/protobuf/proto"
)

var obj = atomic.Value{}

type RummyConfigData struct {
	_List                     []*pb.RummyConfig
	_GameTypeRoomTypeCoinType map[pb.Index3[pb.GameType, pb.RoomType, pb.CoinType]]*pb.RummyConfig
	_GameTypeCoinType         map[pb.Index2[pb.GameType, pb.CoinType]][]*pb.RummyConfig
	_MaxID                    uint32
	_ID                       map[uint32]*pb.RummyConfig
}

// 注册函数
func init() {
	config.Register("RummyConfig", parse)
}

func parse(buf string) error {
	data := &pb.RummyConfigAry{}
	if err := proto.UnmarshalText(buf, data); err != nil {
		return err
	}

	_GameTypeRoomTypeCoinType := make(map[pb.Index3[pb.GameType, pb.RoomType, pb.CoinType]]*pb.RummyConfig)
	_GameTypeCoinType := make(map[pb.Index2[pb.GameType, pb.CoinType]][]*pb.RummyConfig)
	var _MaxID uint32
	_ID := make(map[uint32]*pb.RummyConfig)
	for _, item := range data.Ary {
		keyGameTypeRoomTypeCoinType := pb.Index3[pb.GameType, pb.RoomType, pb.CoinType]{item.GameType, item.RoomType, item.CoinType}
		_GameTypeRoomTypeCoinType[keyGameTypeRoomTypeCoinType] = item
		keyGameTypeCoinType := pb.Index2[pb.GameType, pb.CoinType]{item.GameType, item.CoinType}
		_GameTypeCoinType[keyGameTypeCoinType] = append(_GameTypeCoinType[keyGameTypeCoinType], item)
		if _MaxID < item.ID {
			_MaxID = item.ID
		}
		_ID[item.ID] = item
	}

	obj.Store(&RummyConfigData{
		_List:                     data.Ary,
		_GameTypeRoomTypeCoinType: _GameTypeRoomTypeCoinType,
		_GameTypeCoinType:         _GameTypeCoinType,
		_MaxID:                    _MaxID,
		_ID:                       _ID,
	})
	return nil
}

func SGet() *pb.RummyConfig {
	obj, ok := obj.Load().(*RummyConfigData)
	if !ok {
		return nil
	}
	return obj._List[len(obj._List)-1]
}

func LGet() (rets []*pb.RummyConfig) {
	obj, ok := obj.Load().(*RummyConfigData)
	if !ok {
		return
	}
	rets = make([]*pb.RummyConfig, len(obj._List))
	copy(rets, obj._List)
	return
}

func Walk(f func(*pb.RummyConfig) bool) {
	obj, ok := obj.Load().(*RummyConfigData)
	if !ok {
		return
	}
	for _, item := range obj._List {
		if !f(item) {
			return
		}
	}
}

func MGetGameTypeRoomTypeCoinType(GameType pb.GameType, RoomType pb.RoomType, CoinType pb.CoinType) *pb.RummyConfig {
	obj, ok := obj.Load().(*RummyConfigData)
	if !ok {
		return nil
	}
	if val, ok := obj._GameTypeRoomTypeCoinType[pb.Index3[pb.GameType, pb.RoomType, pb.CoinType]{GameType, RoomType, CoinType}]; ok {
		return val
	}
	return nil
}

func GGetGameTypeCoinType(GameType pb.GameType, CoinType pb.CoinType) (rets []*pb.RummyConfig) {
	obj, ok := obj.Load().(*RummyConfigData)
	if !ok {
		return
	}
	if vals, ok := obj._GameTypeCoinType[pb.Index2[pb.GameType, pb.CoinType]{GameType, CoinType}]; ok {
		rets = make([]*pb.RummyConfig, len(vals))
		copy(rets, vals)
		return
	}
	return
}

func MGetIDKey(val uint32) uint32 {
	if obj, ok := obj.Load().(*RummyConfigData); ok && val > obj._MaxID {
		return obj._MaxID
	}
	return val
}

func MGetID(ID uint32) *pb.RummyConfig {
	obj, ok := obj.Load().(*RummyConfigData)
	if !ok {
		return nil
	}
	if val, ok := obj._ID[ID]; ok {
		return val
	}
	return nil
}
