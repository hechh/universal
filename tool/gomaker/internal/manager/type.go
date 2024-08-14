package manager

import (
	"encoding/json"
	"fmt"
	"universal/tool/gomaker/domain"
)

var (
	types   = make(map[string]*domain.Type)
	values  = make(map[string]map[int32]*domain.Value)
	structs = make(map[string]*domain.Struct)
	alias   = make(map[string]*domain.Alias)
)

func Print() string {
	buf, _ := json.Marshal(&values)
	return string(buf)
}

func GetOrAddType(tt *domain.Type) *domain.Type {
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

func AddValue(vv *domain.Value) {
	// 存储数据
	key := fmt.Sprintf("%s.%s", vv.Type.Selector, vv.Type.Name)
	if _, ok := values[key]; !ok {
		values[key] = make(map[int32]*domain.Value)
	}
	values[key][vv.Value] = vv
}

func AddStruct(vv *domain.Struct) {
	structs[fmt.Sprintf("%s.%s", vv.Type.Selector, vv.Type.Name)] = vv
}

func AddAlias(vv *domain.Alias) {
	alias[fmt.Sprintf("%s.%s", vv.AliasType.Selector, vv.AliasType.Name)] = vv
}
