package manager

import (
	"universal/tools/gomaker/domain"

	"github.com/spf13/cast"
)

var (
	trans = make(map[string]string)          // 配置类型——>proto类型
	convs = make(map[string]domain.ConvFunc) // 配置字段值转型
	evals = make(map[string]int32)           // 配置枚举值准换
)

func InitEvals() {
	for _, item := range enums {
		for _, val := range item.List {
			evals[val.Doc] = val.Value
		}
	}
}

func AddTrans(key, val string) {
	trans[key] = val
}

func AddConv(key string, f domain.ConvFunc) {
	convs[key] = f
}

func GetProtoType(typ string) (string, bool) {
	val, ok := trans[typ]
	return val, ok
}

// 默认枚举值转换函数
func defaultEnumConv(str string) interface{} {
	if val, ok := evals[str]; ok {
		return val
	}
	return cast.ToInt32(str)
}

func Cast(typ, str string) interface{} {
	if val, ok := convs[typ]; ok {
		return val(str)
	}
	return defaultEnumConv(str)
}
