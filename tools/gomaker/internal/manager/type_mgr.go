package manager

import (
	"fmt"
	"go/ast"
	"universal/tools/gomaker/internal/typespec"
)

var (
	structs = make(map[string]*typespec.Struct) // 结构类型
	enums   = make(map[string]*typespec.Enum)   // 枚举类型
	alias   = make(map[string]*typespec.Alias)  // 类型别名
)

func initPkgName(typ *typespec.Type) {
	if typ.Key != nil {
		initPkgName(typ.Key)
		initPkgName(typ.Value)
	} else {
		if val, ok := alias[typ.Name]; ok {
			typ.PkgName = val.PkgName
			return
		}
		if val, ok := structs[typ.Name]; ok {
			typ.PkgName = val.PkgName
			return
		}
	}
}

func Finished() {
	for _, aa := range alias {
		initPkgName(aa.Type)
	}
	for _, st := range structs {
		for _, ff := range st.List {
			initPkgName(ff.Type)
		}
	}
}

func GetMapStruct() map[string]*typespec.Struct {
	return structs
}

func GetMapEnum() map[string]*typespec.Enum {
	return enums
}

func GetAlias(name string) *typespec.Alias {
	return alias[name]
}

func GetEnum(name string) *typespec.Enum {
	return enums[name]
}

func GetStruct(name string) *typespec.Struct {
	return structs[name]
}

func AddType(pkgName string, specs []ast.Spec) {
	for _, spec := range specs {
		// 判断是否有效
		node, ok := spec.(*ast.TypeSpec)
		if !ok || node == nil {
			continue
		}
		// 解析结构
		switch vv := node.Type.(type) {
		case *ast.StructType:
			structs[node.Name.Name] = typespec.NewStruct(pkgName, node.Name.Name, vv.Fields.List)
		default:
			alias[node.Name.Name] = typespec.NewAlias(pkgName, node.Name.Name, node)
		}
	}
}

func AddConst(pkgName string, specs []ast.Spec) {
	ee := typespec.NewEnum(pkgName, specs)
	if _, ok := alias[ee.Name]; ok {
		enums[ee.Name] = ee
	}
}

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
