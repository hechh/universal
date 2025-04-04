package manager

import (
	"hego/framework/uerror"
	"hego/tools/gomaker/internal/typespec"
	"sort"
	"strings"
)

var (
	types   = make(map[string]*typespec.Type)
	alias   = make(map[string]*typespec.Alias)
	enums   = make(map[string]*typespec.Enum)
	structs = make(map[string]*typespec.Struct)
)

// 查询类型
func GetType(k int32, pkg, name, class string) *typespec.Type {
	return GetTypeReference(typespec.TYPE(k, pkg, name, class))
}

func GetTypeReference(tt *typespec.Type) *typespec.Type {
	key := tt.GetPkgType()
	val, ok := types[key]
	if !ok {
		types[key] = tt
		return tt
	}
	if tt.Kind > 0 {
		val.Kind = tt.Kind
	}
	if len(tt.Class) > 0 {
		val.Class = tt.Class
	}
	return val
}

// ---------------------添加别名------------------------
func AddAlias(t *typespec.Alias) error {
	key := t.Type.GetPkgType()
	if _, ok := alias[key]; ok {
		return uerror.NewUError(2, -1, "别名(%s)已经存在", key)
	}
	alias[key] = t
	return nil
}

// 查询别名
func GetAlias(key string) *typespec.Alias {
	return alias[key]
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

// ---------------添加枚举-----------------------
func AddEnum(vv *typespec.Enum) error {
	if vv != nil {
		key := vv.Type.GetPkgType()
		if _, ok := enums[key]; ok {
			return uerror.NewUError(2, -1, "枚举类型(%s)已经存在", key)
		}
		enums[key] = vv
	}
	return nil
}

func LoadEnum(tt *typespec.Type) *typespec.Enum {
	key := tt.GetPkgType()
	if _, ok := enums[key]; !ok {
		enums[key] = typespec.ENUM(tt, "")
	}
	return enums[key]
}

func GetEnum(key string) *typespec.Enum {
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

// -----------------------struct数据结构----------------------
func AddStruct(vv *typespec.Struct) error {
	key := vv.Type.GetPkgType()
	if _, ok := structs[key]; ok {
		return uerror.NewUError(2, -1, "结构体(%s)已经存在", key)
	}
	structs[key] = vv
	return nil
}

func GetStruct(key string) *typespec.Struct {
	return structs[key]
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
