package manager

import (
	"encoding/json"
	"sort"
	"strings"
	"universal/tools/gomaker/internal/typespec"
)

var (
	types   = make(map[string]*typespec.Type)
	alias   = make(map[string]*typespec.Alias)
	enums   = make(map[string]*typespec.Enum)
	structs = make(map[string]*typespec.Struct)
)

func GetOrAddType(tt *typespec.Type) *typespec.Type {
	key := tt.GetName("")
	if val, ok := types[key]; !ok {
		types[key] = tt
	} else {
		if len(tt.Doc) > 0 {
			val.Doc = tt.Doc
		}
		if tt.Kind > 0 {
			val.Kind = tt.Kind
		}
	}
	return types[key]
}

// -----------alias----------
func AddAlias(vv *typespec.Alias) {
	if vv != nil {
		alias[vv.Type.GetName("")] = vv
	}
}

func GetAliasList() (rets []*typespec.Alias) {
	tmps := []*typespec.Alias{}
	for _, val := range alias {
		tmps = append(tmps, val)
	}
	// 深度拷贝
	buf, _ := json.Marshal(&tmps)
	json.Unmarshal(buf, &rets)
	sort.Slice(rets, func(i, j int) bool {
		return strings.Compare(rets[i].Type.Name, rets[j].Type.Name) < 0
	})
	return
}

// -------struct-----------
func AddStruct(vv *typespec.Struct) {
	if vv != nil {
		structs[vv.Type.GetName("")] = vv
	}
}

func GetStruct(name string) *typespec.Struct {
	return structs[name]
}

func GetStructList() (rets []*typespec.Struct) {
	tmps := []*typespec.Struct{}
	for _, val := range structs {
		tmps = append(tmps, val)
	}
	// 深度拷贝
	buf, _ := json.Marshal(&tmps)
	json.Unmarshal(buf, &rets)
	sort.Slice(rets, func(i, j int) bool {
		return strings.Compare(rets[i].Type.Name, rets[j].Type.Name) < 0
	})
	return
}

// -------枚举类型---------
func AddValue(filename string, vv *typespec.Value) {
	if vv != nil {
		name := vv.Type.GetName("")
		if eval, ok := enums[name]; !ok {
			enums[name] = typespec.NewEnum(filename, vv.Type).Add(vv)
		} else {
			eval.Add(vv)
		}
	}
}

func AddEnum(vv *typespec.Enum) {
	if vv != nil {
		name := vv.Type.GetName("")
		if eval, ok := enums[name]; !ok {
			enums[name] = vv
		} else {
			for _, item := range vv.Values {
				eval.Add(item)
			}
		}
	}
}

func GetEnumList() (rets []*typespec.Enum) {
	tmps := []*typespec.Enum{}
	for _, val := range enums {
		tmps = append(tmps, val)
	}
	// 深度拷贝
	buf, _ := json.Marshal(&tmps)
	json.Unmarshal(buf, &rets)
	sort.Slice(rets, func(i, j int) bool {
		return strings.Compare(rets[i].Type.Name, rets[j].Type.Name) < 0
	})
	return
}
