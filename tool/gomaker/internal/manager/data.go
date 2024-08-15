package manager

import (
	"fmt"
	"universal/tool/gomaker/internal/typespec"
)

var (
	types   = make(map[string]*typespec.Type)
	enums   = make(map[string]*typespec.Enum)
	structs = make(map[string]*typespec.Struct)
	alias   = make(map[string]*typespec.Alias)
)

func Print() {
	/*
		for key, val := range alias {
			fmt.Println(val.Format())
			if vv, ok := enums[key]; ok {
				fmt.Println(vv.Format())
			}
		}
	*/
	for _, val := range structs {
		fmt.Println(val.Format())
	}
}

func GetOrAddType(tt *typespec.Type) *typespec.Type {
	key := fmt.Sprintf("%s.%s", tt.Selector, tt.Name)
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

func AddEnum(vv *typespec.Enum) {
	if vv != nil {
		enums[fmt.Sprintf("%s.%s", vv.Type.Selector, vv.Type.Name)] = vv
	}
}

func AddStruct(vv *typespec.Struct) {
	if vv != nil {
		structs[fmt.Sprintf("%s.%s", vv.Type.Selector, vv.Type.Name)] = vv
	}
}

func AddAlias(vv *typespec.Alias) {
	if vv != nil {
		alias[fmt.Sprintf("%s.%s", vv.Type.Selector, vv.Type.Name)] = vv
	}
}
