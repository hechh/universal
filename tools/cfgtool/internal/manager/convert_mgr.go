package manager

import (
	"hego/tools/cfgtool/internal/base"

	"github.com/spf13/cast"
)

var (
	convertMgr = make(map[string]*base.Convert)
)

func GetConvFunc(name string) func(string) interface{} {
	if val, ok := convertMgr[name]; ok {
		return val.ConvFunc
	}

	// 默认枚举转换函数
	if item, ok := enumMgr[name]; ok {
		return func(str string) interface{} {
			if vv, ok := item.Values[str]; ok {
				return cast.ToInt32(vv)
			}
			return cast.ToInt32(str)
		}
	}
	return nil
}

func GetConvType(name string) string {
	if val, ok := convertMgr[name]; ok {
		return val.Name
	}
	return name
}
