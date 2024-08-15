package manager

import (
	"encoding/json"
	"fmt"
	"universal/tool/gomaker/domain"
	"universal/tool/gomaker/internal/typespec"
)

var (
	types   = make(map[string]*typespec.Type)
	values  = make(map[string]map[int32]*typespec.Value)
	structs = make(map[string]*typespec.Struct)
	alias   = make(map[string]*typespec.Alias)
	apis    = make(map[string]*ApiInfo)
)

type ApiInfo struct {
	help string
	gen  domain.GenFunc
}

func Print() string {
	buf, _ := json.Marshal(&values)
	return string(buf)
}

func Register(act string, g domain.GenFunc, help string) {
	apis[act] = &ApiInfo{help: help, gen: g}
}

func GetApi(act string) *ApiInfo {
	return apis[act]
}

func GetOrAddType(tt *typespec.Type) *typespec.Type {
	key := fmt.Sprintf("%s.%s", tt.Selector, tt.Name)
	if val, ok := types[key]; !ok {
		types[key] = tt
	} else {
		if len(tt.Doc) <= 0 {
			val.Doc = tt.Doc
		}
		if tt.Kind > 0 {
			val.Kind = tt.Kind
		}
	}
	return types[key]
}

func AddValue(vv *typespec.Value) {
	// 存储数据
	key := fmt.Sprintf("%s.%s", vv.Type.Selector, vv.Type.Name)
	if _, ok := values[key]; !ok {
		values[key] = make(map[int32]*typespec.Value)
	}
	values[key][vv.Value] = vv
}

func AddStruct(vv *typespec.Struct) {
	structs[fmt.Sprintf("%s.%s", vv.Type.Selector, vv.Type.Name)] = vv
}

func AddAlias(vv *typespec.Alias) {
	alias[fmt.Sprintf("%s.%s", vv.AliasType.Selector, vv.AliasType.Name)] = vv
}
