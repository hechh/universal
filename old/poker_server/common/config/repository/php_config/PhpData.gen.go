/*
* 本代码由cfgtool工具生成，请勿手动修改
 */

package php_config

import (
	"poker_server/common/config"
	"poker_server/common/pb"
	"sync/atomic"

	"github.com/golang/protobuf/proto"
)

var obj = atomic.Value{}

type PhpConfigData struct {
	_List        []*pb.PhpConfig
	_EnvTypeName map[pb.Index2[pb.EnvType, string]]*pb.PhpConfig
}

// 注册函数
func init() {
	config.Register("PhpConfig", parse)
}

func parse(buf string) error {
	data := &pb.PhpConfigAry{}
	if err := proto.UnmarshalText(buf, data); err != nil {
		return err
	}

	_EnvTypeName := make(map[pb.Index2[pb.EnvType, string]]*pb.PhpConfig)
	for _, item := range data.Ary {
		keyEnvTypeName := pb.Index2[pb.EnvType, string]{item.EnvType, item.Name}
		_EnvTypeName[keyEnvTypeName] = item
	}

	obj.Store(&PhpConfigData{
		_List:        data.Ary,
		_EnvTypeName: _EnvTypeName,
	})
	return nil
}

func SGet() *pb.PhpConfig {
	obj, ok := obj.Load().(*PhpConfigData)
	if !ok {
		return nil
	}
	return obj._List[len(obj._List)-1]
}

func LGet() (rets []*pb.PhpConfig) {
	obj, ok := obj.Load().(*PhpConfigData)
	if !ok {
		return
	}
	rets = make([]*pb.PhpConfig, len(obj._List))
	copy(rets, obj._List)
	return
}

func Walk(f func(*pb.PhpConfig) bool) {
	obj, ok := obj.Load().(*PhpConfigData)
	if !ok {
		return
	}
	for _, item := range obj._List {
		if !f(item) {
			return
		}
	}
}

func MGetEnvTypeName(EnvType pb.EnvType, Name string) *pb.PhpConfig {
	obj, ok := obj.Load().(*PhpConfigData)
	if !ok {
		return nil
	}
	if val, ok := obj._EnvTypeName[pb.Index2[pb.EnvType, string]{EnvType, Name}]; ok {
		return val
	}
	return nil
}
