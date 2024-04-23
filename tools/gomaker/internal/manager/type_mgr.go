package manager

import (
	"fmt"
	"go/ast"
	"universal/tools/gomaker/internal/typespec"
)

var (
	typeMgr = NewTypeMgr()
)

type TypeMgr struct {
	structs map[string]*typespec.Struct // 结构类型
	enums   map[string]*typespec.Enum   // 枚举类型
	alias   map[string]*typespec.Alias  // 类型别名
}

func GetTypeMgr() *TypeMgr {
	return typeMgr
}

func NewTypeMgr() *TypeMgr {
	return &TypeMgr{
		structs: make(map[string]*typespec.Struct),
		enums:   make(map[string]*typespec.Enum),
		alias:   make(map[string]*typespec.Alias),
	}
}

func (d *TypeMgr) GetAlias(name string) *typespec.Alias {
	return d.alias[name]
}

func (d *TypeMgr) GetEnum(name string) *typespec.Enum {
	return d.enums[name]
}

func (d *TypeMgr) GetStruct(name string) *typespec.Struct {
	return d.structs[name]
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
			d.structs[node.Name.Name] = typespec.NewStruct(pkgName, node.Name.Name, vv.Fields.List)
		default:
			d.alias[node.Name.Name] = typespec.NewAlias(pkgName, node.Name.Name, node)
		}
	}
}

func (d *TypeMgr) AddConst(pkgName string, specs []ast.Spec) {
	ee := typespec.NewEnum(pkgName, specs)
	if _, ok := d.alias[ee.Name]; ok {
		d.enums[ee.Name] = ee
	}
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
