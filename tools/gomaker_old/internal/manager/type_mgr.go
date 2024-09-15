package manager

import (
	"sort"
	"strings"
	"universal/framework/uerror"
	"universal/tools/gomaker/internal/typespec"
)

var (
	types   = make(map[string]*typespec.Type)
	alias   = make(map[string]*typespec.Alias)
	enums   = make(map[string]*typespec.Enum)
	evals   = make(map[string]*typespec.Value)
	structs = make(map[string]*typespec.Struct)
)

func InitEvals() {
	for _, item := range enums {
		for _, val := range item.List {
			evals[val.Doc] = val
		}
	}
}

func AddEnum(vv *typespec.Enum) error {
	name := vv.Type.GetPkgType()
	if _, ok := enums[name]; ok {
		return uerror.NewUError(2, -1, "枚举类型(%s)已经存在", name)
	}
	enums[name] = vv
	return nil
}

func GetEnum(name string) *typespec.Enum {
	return enums[name]
}

func GetOrNewEnum(tt *typespec.Type) *typespec.Enum {
	key := tt.GetPkgType()
	if _, ok := enums[key]; !ok {
		enums[key] = &typespec.Enum{Type: tt, Values: make(map[string]*typespec.Value)}
	}
	return enums[key]
}

func GetEnumList() (rets []*typespec.Enum) {
	for _, val := range enums {
		rets = append(rets, val)
	}
	sort.Slice(rets, func(i, j int) bool {
		return strings.Compare(rets[i].Type.Name, rets[j].Type.Name) < 0
	})
	return
}

func AddStruct(vv *typespec.Struct) error {
	name := vv.Type.GetPkgType()
	if _, ok := structs[name]; ok {
		return uerror.NewUError(2, -1, "结构体(%s)已经存在", name)
	}
	structs[name] = vv
	return nil
}

func GetStruct(name string) *typespec.Struct {
	return structs[name]
}

func GetStructList() (rets []*typespec.Struct) {
	for _, val := range structs {
		rets = append(rets, val)
	}
	sort.Slice(rets, func(i, j int) bool {
		return strings.Compare(rets[i].Type.Name, rets[j].Type.Name) < 0
	})
	return
}

func AddAlias(vv *typespec.Alias) error {
	name := vv.Type.GetPkgType()
	if _, ok := alias[name]; ok {
		return uerror.NewUError(2, -1, "别名(%s)已经存在", name)
	}
	alias[name] = vv
	return nil
}

func GetAlias(name string) *typespec.Alias {
	return alias[name]
}

func GetAliasList() (rets []*typespec.Alias) {
	for _, val := range alias {
		rets = append(rets, val)
	}
	sort.Slice(rets, func(i, j int) bool {
		return strings.Compare(rets[i].Type.Name, rets[j].Type.Name) < 0
	})
	return
}

func GetTypeReference(tt *typespec.Type) *typespec.Type {
	key := tt.GetPkgType()
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
