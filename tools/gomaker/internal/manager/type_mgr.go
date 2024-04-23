package manager

import (
	"fmt"
	"go/ast"
	"universal/tools/gomaker/internal/base"
)

type TypeMgr struct {
	structs map[string]*base.Struct // 结构类型
	enums   map[string]*base.Enum   // 枚举类型
	alias   map[string]*base.Alias  // 类型别名
}

func NewTypeMgr() *TypeMgr {
	return &TypeMgr{
		structs: make(map[string]*base.Struct),
		enums:   make(map[string]*base.Enum),
		alias:   make(map[string]*base.Alias),
	}
}

func (d *TypeMgr) AddType(pkgName string, specs []ast.Spec) {
	for _, spec := range specs {
		// 判断是否有效
		node, ok := spec.(*ast.TypeSpec)
		if !ok || node == nil {
			continue
		}
		// 解析结构
		switch vv := node.Type.(type) {
		case *ast.StructType:
			d.structs[node.Name.Name] = base.NewStruct(pkgName, node.Name.Name, vv.Fields.List)
		default:
			d.alias[node.Name.Name] = base.NewAlias(pkgName, node.Name.Name, node)
		}
	}
}

func (d *TypeMgr) AddConst(pkgName string, specs []ast.Spec) {
	ee := base.NewEnum(pkgName, specs)
	d.enums[ee.GetType("")] = ee
}

func (d *TypeMgr) Print() {
	for _, st := range d.structs {
		fmt.Println(st.String())
	}
	for _, en := range d.enums {
		fmt.Println(en.String())
	}
	for _, al := range d.alias {
		fmt.Println(al.String())
	}
}
