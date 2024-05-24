package manager

import (
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/common/types"
)

var (
	structs = make(map[string]*types.Struct) // 结构类型
	enums   = make(map[string]*types.Enum)   // 枚举类型
	alias   = make(map[string]*types.Alias)  // 类型别名
)

func GetMapStruct() map[string]*types.Struct {
	return structs
}

func GetStruct(name string) *types.Struct {
	return structs[name]
}

func GetMapEnum() map[string]*types.Enum {
	return enums
}

func GetEnum(name string) *types.Enum {
	return enums[name]
}

func GetAlias(name string) *types.Alias {
	return alias[name]
}

func AddType(data interface{}) {
	switch vv := data.(type) {
	case *types.Struct:
		structs[vv.Type.Name] = vv
	case *types.Alias:
		alias[vv.Type.Name] = vv
	case *types.Enum:
		enums[vv.Type.Name] = vv
	}
}

func Finish() {
	typeF := func(t *types.Type) {
		if vv, ok := structs[t.Name]; ok {
			t.Token |= domain.STRUCT
			t.PkgName = vv.Type.PkgName
		}
		if vv, ok := enums[t.Name]; ok {
			t.Token |= domain.ENUM
			t.PkgName = vv.Type.Name
		}
		if vv, ok := alias[t.Name]; ok {
			t.Token |= domain.ALIAS
			t.PkgName = vv.Type.Name
		}
	}
	for _, st := range structs {
		for _, field := range st.Fields {
			typeF(field.Type)
		}
	}
	for _, al := range alias {
		typeF(al.Reference)
	}
}

/*
func Print() {
	for _, st := range structs {
		fmt.Println(st.String())
	}
	for _, en := range enums {
		fmt.Println(en.String())
	}
	for _, al := range alias {
		fmt.Println(al.String())
	}
}
*/
