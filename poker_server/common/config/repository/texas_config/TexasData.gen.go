/*
* 本代码由cfgtool工具生成，请勿手动修改
 */

package texas_config

import (
	"poker_server/common/config"
	"poker_server/common/pb"
	"sync/atomic"

	"github.com/golang/protobuf/proto"
)

var obj = atomic.Value{}

type TexasConfigData struct {
	_List                      []*pb.TexasConfig
	_GameTypeMatchTypeCoinType map[pb.Index3[pb.GameType, pb.MatchType, pb.CoinType]][]*pb.TexasConfig
	_MaxID                     int32
	_ID                        map[int32]*pb.TexasConfig
	_MatchType                 map[pb.MatchType][]*pb.TexasConfig
}

// 注册函数
func init() {
	config.Register("TexasConfig", parse)
}

func parse(buf string) error {
	data := &pb.TexasConfigAry{}
	if err := proto.UnmarshalText(buf, data); err != nil {
		return err
	}

	_GameTypeMatchTypeCoinType := make(map[pb.Index3[pb.GameType, pb.MatchType, pb.CoinType]][]*pb.TexasConfig)
	var _MaxID int32
	_ID := make(map[int32]*pb.TexasConfig)
	_MatchType := make(map[pb.MatchType][]*pb.TexasConfig)
	for _, item := range data.Ary {
		keyGameTypeMatchTypeCoinType := pb.Index3[pb.GameType, pb.MatchType, pb.CoinType]{item.GameType, item.MatchType, item.CoinType}
		_GameTypeMatchTypeCoinType[keyGameTypeMatchTypeCoinType] = append(_GameTypeMatchTypeCoinType[keyGameTypeMatchTypeCoinType], item)
		if _MaxID < item.ID {
			_MaxID = item.ID
		}
		_ID[item.ID] = item
		_MatchType[item.MatchType] = append(_MatchType[item.MatchType], item)
	}

	obj.Store(&TexasConfigData{
		_List:                      data.Ary,
		_GameTypeMatchTypeCoinType: _GameTypeMatchTypeCoinType,
		_MaxID:                     _MaxID,
		_ID:                        _ID,
		_MatchType:                 _MatchType,
	})
	return nil
}

func SGet() *pb.TexasConfig {
	obj, ok := obj.Load().(*TexasConfigData)
	if !ok {
		return nil
	}
	return obj._List[len(obj._List)-1]
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

func GGetGameTypeMatchTypeCoinType(GameType pb.GameType, MatchType pb.MatchType, CoinType pb.CoinType) (rets []*pb.TexasConfig) {
	obj, ok := obj.Load().(*TexasConfigData)
	if !ok {
		return
	}
	if vals, ok := obj._GameTypeMatchTypeCoinType[pb.Index3[pb.GameType, pb.MatchType, pb.CoinType]{GameType, MatchType, CoinType}]; ok {
		rets = make([]*pb.TexasConfig, len(vals))
		copy(rets, vals)
		return
	}
	return
}

func MGetIDKey(val int32) int32 {
	if obj, ok := obj.Load().(*TexasConfigData); ok && val > obj._MaxID {
		return obj._MaxID
	}
	return val
}

func MGetID(ID int32) *pb.TexasConfig {
	obj, ok := obj.Load().(*TexasConfigData)
	if !ok {
		return nil
	}
	if val, ok := obj._ID[ID]; ok {
		return val
	}
	return nil
}

func GGetMatchType(MatchType pb.MatchType) (rets []*pb.TexasConfig) {
	obj, ok := obj.Load().(*TexasConfigData)
	if !ok {
		return
	}
	if vals, ok := obj._MatchType[MatchType]; ok {
		rets = make([]*pb.TexasConfig, len(vals))
		copy(rets, vals)
		return
	}
	return
}
