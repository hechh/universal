package manager

import (
	"universal/tools/gomaker/domain"

	"github.com/spf13/cast"
)

var (
	convs = make(map[string]*typeInfo) // 配置类型转成golang类型
	evals = make(map[string]int32)     // 配置枚举值准换
)

type typeInfo struct {
	cfgType string
	goType  string
	pbType  string
	conv    domain.ConvFunc
}

func InitEvals() {
	for _, item := range enums {
		for _, val := range item.List {
			evals[val.Doc] = val.Value
		}
	}
}

func AddConv(cfg, g, p string, f domain.ConvFunc) {
	convs[cfg] = &typeInfo{cfgType: cfg, goType: g, pbType: p, conv: f}
}

func GetGoType(cfg string) string {
	if val, ok := convs[cfg]; ok {
		return val.goType
	}
	return cfg
}

func GetPbType(cfg string) string {
	return convs[cfg].pbType
}

func GetConv(typ string) domain.ConvFunc {
	if val, ok := convs[typ]; ok {
		return val.conv
	}
	return defaultEnumConv
}

func Cast(typ, str string) interface{} {
	if val, ok := convs[typ]; ok {
		return val.conv(str)
	}
	return defaultEnumConv(str)
}

// 默认枚举值转换函数
func defaultEnumConv(str string) interface{} {
	if val, ok := evals[str]; ok {
		return val
	}
	return cast.ToInt32(str)
}
