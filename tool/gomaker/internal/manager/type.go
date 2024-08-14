package manager

import (
	"fmt"
	"universal/tool/gomaker/domain"
)

var (
	types   = make(map[string]*domain.Type)
	values  = make(map[string]map[int32]*domain.Value)
	structs = make(map[string]*domain.Struct)
	alias   = make(map[string]*domain.Alias)
)

func AddValue(vv *domain.Value) {
	// 存储类型
	key := fmt.Sprintf("%s.%s", vv.Type.Selector, vv.Type.Name)
	if _, ok := types[key]; !ok {
		types[key] = vv.Type
	}
	// 存储数据
	if _, ok := values[key]; !ok {
		values[key] = make(map[int32]*domain.Value)
	}
	values[key][vv.Value] = vv
}

func AddStruct(vv *domain.Struct) {
	// 存储类型
	key := fmt.Sprintf("%s.%s", vv.Type.Selector, vv.Type.Name)
	if _, ok := types[key]; !ok {
		types[key] = vv.Type
	}
	// field类型存储
	for _, field := range vv.List {
		fkey := fmt.Sprintf("%s.%s", field.Type.Selector, field.Type.Name)
		if _, ok := types[fkey]; !ok {
			types[fkey] = field.Type
		}
	}
	// 存储struct数据
	structs[key] = vv
}

func AddAlias(vv *domain.Alias) {
	// 存储类型
	key := fmt.Sprintf("%s.%s", vv.AliasType.Selector, vv.AliasType.Name)
	if _, ok := types[key]; !ok {
		types[key] = vv.AliasType
	}
	rkey := fmt.Sprintf("%s.%s", vv.Type.Selector, vv.Type.Name)
	if _, ok := types[rkey]; !ok {
		types[rkey] = vv.Type
	}
	// 存储别名
	alias[key] = vv
}
