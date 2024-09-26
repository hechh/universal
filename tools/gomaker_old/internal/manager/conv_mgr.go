package manager

import (
	"strings"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/typespec"

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

func ToCast(typ, str string) interface{} {
	if val, ok := convs[typ]; ok {
		return val.conv(str)
	}
	return defaultEnumConv(str)
}

func Cast(ff *typespec.Field, str string) interface{} {
	ts := []int32{}
	for _, tt := range ff.Token {
		if tt == domain.TokenTypeArray {
			ts = append(ts, tt)
		}
	}
	switch len(ts) {
	case 1:
		rets := []interface{}{}
		for _, val := range strings.Split(str, ",") {
			rets = append(rets, ToCast(ff.Type.Name, val))
		}
		return rets
	case 2:
		rets := []interface{}{}
		for _, vals := range strings.Split(str, "|") {
			tts := []interface{}{}
			for _, val := range strings.Split(vals, ",") {
				tts = append(tts, ToCast(ff.Type.Name, val))
			}
			rets = append(rets, tts)
		}
		return rets
	}
	return ToCast(ff.Type.Name, str)
}

// 默认枚举值转换函数
func defaultEnumConv(str string) interface{} {
	if val, ok := evals[str]; ok {
		return val
	}
	return cast.ToInt32(str)
}
