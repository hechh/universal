/*
* 本代码由cfgtool工具生成，请勿手动修改
 */

package sng_match_config

import (
	"poker_server/common/config"
	"poker_server/common/pb"
	"sync/atomic"

	"github.com/golang/protobuf/proto"
)

var obj = atomic.Value{}

type SngMatchConfigData struct {
	_List      []*pb.SngMatchConfig
	_MaxGameId int32
	_GameId    map[int32]*pb.SngMatchConfig
}

// 注册函数
func init() {
	config.Register("SngMatchConfig", parse)
}

func parse(buf string) error {
	data := &pb.SngMatchConfigAry{}
	if err := proto.UnmarshalText(buf, data); err != nil {
		return err
	}

	var _MaxGameId int32
	_GameId := make(map[int32]*pb.SngMatchConfig)
	for _, item := range data.Ary {
		if _MaxGameId < item.GameId {
			_MaxGameId = item.GameId
		}
		_GameId[item.GameId] = item
	}

	obj.Store(&SngMatchConfigData{
		_List:      data.Ary,
		_MaxGameId: _MaxGameId,
		_GameId:    _GameId,
	})
	return nil
}

func SGet() *pb.SngMatchConfig {
	obj, ok := obj.Load().(*SngMatchConfigData)
	if !ok {
		return nil
	}
	return obj._List[len(obj._List)-1]
}

func LGet() (rets []*pb.SngMatchConfig) {
	obj, ok := obj.Load().(*SngMatchConfigData)
	if !ok {
		return
	}
	rets = make([]*pb.SngMatchConfig, len(obj._List))
	copy(rets, obj._List)
	return
}

func Walk(f func(*pb.SngMatchConfig) bool) {
	obj, ok := obj.Load().(*SngMatchConfigData)
	if !ok {
		return
	}
	for _, item := range obj._List {
		if !f(item) {
			return
		}
	}
}

func MGetGameIdKey(val int32) int32 {
	if obj, ok := obj.Load().(*SngMatchConfigData); ok && val > obj._MaxGameId {
		return obj._MaxGameId
	}
	return val
}

func MGetGameId(GameId int32) *pb.SngMatchConfig {
	obj, ok := obj.Load().(*SngMatchConfigData)
	if !ok {
		return nil
	}
	if val, ok := obj._GameId[GameId]; ok {
		return val
	}
	return nil
}
