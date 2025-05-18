package manager

import (
	"sort"
	"strings"

	"forevernine.com/planet/server/tool/gomaker/domain"
)

var (
	_structs   = make(map[string]*domain.AstStruct)      // pbname--->AstStruct
	_enums     = make(map[string]*domain.AstEnum)        // type--->AstEnum
	_pbconfigs = make(map[string]struct{})               // 配置协议需要特殊处理
	_rules     = make(map[string]map[string]interface{}) // rule ----> pbname ---> 解析的规则数据
)

func WalkConfig(f func(item *domain.AstStruct)) {
	strs := []string{}
	for pbname := range _pbconfigs {
		strs = append(strs, pbname)
	}
	sort.Slice(strs, func(i, j int) bool {
		return strings.Compare(strs[i], strs[j]) <= 0
	})
	for _, pbname := range strs {
		f(_structs[pbname])
	}
}

func WalkAstEnum(typ string, f func(item *domain.AstValue)) {
	vals, ok := _enums[typ]
	if !ok {
		return
	}
	for _, item := range vals.Values {
		f(item)
	}
}

func GetAstStruct(pbname string) *domain.AstStruct {
	if val, ok := _structs[pbname]; ok {
		return val
	}
	return nil
}

func GetAstEnum(typ string) *domain.AstEnum {
	if val, ok := _enums[typ]; ok {
		return val
	}
	return nil
}

func GetRules(rule string) map[string]interface{} {
	return _rules[rule]
}

func GetRule(rule, pbname string) interface{} {
	if val, ok := _rules[rule]; ok {
		return val[pbname]
	}
	return nil
}

func AddAstStruct(item *domain.AstStruct) {
	if item == nil {
		return
	}
	_structs[item.Type.Name] = item
}

func AddAstEnum(item *domain.AstEnum) {
	if item == nil {
		return
	}
	_enums[item.Type.Name] = item
}

func Update() {
	for pbname, val := range _structs {
		_, ok := _pbconfigs[pbname]
		astStruct(ok, val)
		if strings.HasSuffix(pbname, "ConfigAry") {
			delete(_pbconfigs, pbname)
		}
	}
}

func astStruct(flag bool, item *domain.AstStruct) {
	for _, val := range item.Idents {
		astType(flag, val.Type)
	}
	for _, val := range item.Arrays {
		astType(flag, val.Type)
	}
	for _, val := range item.Maps {
		astType(flag, val.KType)
		astType(flag, val.VType)
	}
}

func astType(flag bool, item *domain.AstType) {
	if val, ok := _structs[item.Name]; ok {
		item.Token |= domain.STRUCT
		if _, ok := _pbconfigs[val.Type.Name]; flag && !ok {
			_pbconfigs[val.Type.Name] = struct{}{}
			astStruct(flag, val)
		}
		return
	}
	if _, ok := _enums[item.Name]; ok {
		item.Token |= domain.ENUM
		return
	}
	item.Token |= domain.BASE
}
