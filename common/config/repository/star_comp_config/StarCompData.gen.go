/*
* 本代码由cfgtool工具生成，请勿手动修改
 */

package star_comp_config

import (
	"poker_server/common/config"
	"poker_server/common/pb"
	"sync/atomic"

	"github.com/golang/protobuf/proto"
)

var obj = atomic.Value{}

type StarCompConfigData struct {
	_List  []*pb.StarCompConfig
	_MaxID int32
	_ID    map[int32]*pb.StarCompConfig
}

// 注册函数
func init() {
	config.Register("StarCompConfig", parse)
}

func parse(buf string) error {
	data := &pb.StarCompConfigAry{}
	if err := proto.UnmarshalText(buf, data); err != nil {
		return err
	}

	var _MaxID int32
	_ID := make(map[int32]*pb.StarCompConfig)
	for _, item := range data.Ary {
		if _MaxID < item.ID {
			_MaxID = item.ID
		}
		_ID[item.ID] = item
	}

	obj.Store(&StarCompConfigData{
		_List:  data.Ary,
		_MaxID: _MaxID,
		_ID:    _ID,
	})
	return nil
}

func SGet() *pb.StarCompConfig {
	obj, ok := obj.Load().(*StarCompConfigData)
	if !ok {
		return nil
	}
	return obj._List[len(obj._List)-1]
}

func LGet() (rets []*pb.StarCompConfig) {
	obj, ok := obj.Load().(*StarCompConfigData)
	if !ok {
		return
	}
	rets = make([]*pb.StarCompConfig, len(obj._List))
	copy(rets, obj._List)
	return
}

func Walk(f func(*pb.StarCompConfig) bool) {
	obj, ok := obj.Load().(*StarCompConfigData)
	if !ok {
		return
	}
	for _, item := range obj._List {
		if !f(item) {
			return
		}
	}
}

func MGetIDKey(val int32) int32 {
	if obj, ok := obj.Load().(*StarCompConfigData); ok && val > obj._MaxID {
		return obj._MaxID
	}
	return val
}

func MGetID(ID int32) *pb.StarCompConfig {
	obj, ok := obj.Load().(*StarCompConfigData)
	if !ok {
		return nil
	}
	if val, ok := obj._ID[ID]; ok {
		return val
	}
	return nil
}
