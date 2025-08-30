/*
* 本代码由cfgtool工具生成，请勿手动修改
 */

package star_comp_rank_config

import (
	"sync/atomic"
	"universal/common/config"
	"universal/common/pb"

	"github.com/golang/protobuf/proto"
)

var obj = atomic.Value{}

type StarCompRankConfigData struct {
	_List           []*pb.StarCompRankConfig
	_PrizeTypeLevel map[pb.Index2[int32, int32]]*pb.StarCompRankConfig
	_PrizeType      map[int32][]*pb.StarCompRankConfig
}

// 注册函数
func init() {
	config.Register("StarCompRankConfig", parse)
}

func parse(buf string) error {
	data := &pb.StarCompRankConfigAry{}
	if err := proto.UnmarshalText(buf, data); err != nil {
		return err
	}

	_PrizeTypeLevel := make(map[pb.Index2[int32, int32]]*pb.StarCompRankConfig)
	_PrizeType := make(map[int32][]*pb.StarCompRankConfig)
	for _, item := range data.Ary {
		keyPrizeTypeLevel := pb.Index2[int32, int32]{item.PrizeType, item.Level}
		_PrizeTypeLevel[keyPrizeTypeLevel] = item
		_PrizeType[item.PrizeType] = append(_PrizeType[item.PrizeType], item)
	}

	obj.Store(&StarCompRankConfigData{
		_List:           data.Ary,
		_PrizeTypeLevel: _PrizeTypeLevel,
		_PrizeType:      _PrizeType,
	})
	return nil
}

func SGet() *pb.StarCompRankConfig {
	obj, ok := obj.Load().(*StarCompRankConfigData)
	if !ok {
		return nil
	}
	return obj._List[len(obj._List)-1]
}

func LGet() (rets []*pb.StarCompRankConfig) {
	obj, ok := obj.Load().(*StarCompRankConfigData)
	if !ok {
		return
	}
	rets = make([]*pb.StarCompRankConfig, len(obj._List))
	copy(rets, obj._List)
	return
}

func Walk(f func(*pb.StarCompRankConfig) bool) {
	obj, ok := obj.Load().(*StarCompRankConfigData)
	if !ok {
		return
	}
	for _, item := range obj._List {
		if !f(item) {
			return
		}
	}
}

func MGetPrizeTypeLevel(PrizeType int32, Level int32) *pb.StarCompRankConfig {
	obj, ok := obj.Load().(*StarCompRankConfigData)
	if !ok {
		return nil
	}
	if val, ok := obj._PrizeTypeLevel[pb.Index2[int32, int32]{PrizeType, Level}]; ok {
		return val
	}
	return nil
}

func GGetPrizeType(PrizeType int32) (rets []*pb.StarCompRankConfig) {
	obj, ok := obj.Load().(*StarCompRankConfigData)
	if !ok {
		return
	}
	if vals, ok := obj._PrizeType[PrizeType]; ok {
		rets = make([]*pb.StarCompRankConfig, len(vals))
		copy(rets, vals)
		return
	}
	return
}
