/*
* 本代码由cfgtool工具生成，请勿手动修改
 */

package star_comp_config

import (
	"sync/atomic"
	"universal/common/config"
	"universal/common/pb"

	"github.com/golang/protobuf/proto"
)

var obj = atomic.Pointer[StarCompConfigData]{}

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
	if data := obj.Load(); data != nil {
		return data._List[len(data._List)-1]
	}
	return nil
}

func LGet() (rets []*pb.StarCompConfig) {
	if data := obj.Load(); data != nil {
		rets = make([]*pb.StarCompConfig, len(data._List))
		copy(rets, data._List)
	}
	return
}

func Walk(f func(*pb.StarCompConfig) bool) {
	if data := obj.Load(); data != nil {
		for _, item := range data._List {
			if !f(item) {
				return
			}
		}
	}
}

func MGetIDKey(val int32) int32 {
	if data := obj.Load(); data != nil && val > data._MaxID {
		return data._MaxID
	}
	return val
}

func MGetID(ID int32) *pb.StarCompConfig {
	data := obj.Load()
	if data == nil {
		return nil
	}
	if val, ok := data._ID[ID]; ok {
		return val
	}
	return nil
}
