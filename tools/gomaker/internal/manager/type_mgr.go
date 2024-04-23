package manager

import (
	"go/ast"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/base"

	"github.com/spf13/cast"
)

type TypeMgr struct {
	structs map[string]*domain.Struct // 结构类型
	enums   map[string]*domain.Enum   // 枚举类型
	alias   map[string]*domain.Alias  // 类型别名
}

func NewTypeMgr() *TypeMgr {
	return &TypeMgr{
		structs: make(map[string]*domain.Struct),
		enums:   make(map[string]*domain.Enum),
		alias:   make(map[string]*domain.Alias),
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
			item := &domain.Struct{
				PkgName: pkgName,
				Name:    node.Name.Name,
				Fields:  make(map[string]*domain.Field),
			}
			base.ParseStruct(vv.Fields.List, item)
			d.structs[node.Name.Name] = item
		default:
			item := &domain.Alias{
				PkgName: pkgName,
				Name:    node.Name.Name,
				Type:    &domain.Type{},
				Comment: base.ParseComment(node.Comment),
			}
			base.ParseType(vv, item.Type)
			d.alias[node.Name.Name] = item
		}
	}
}

func (d *TypeMgr) AddConst(pkgName string, specs []ast.Spec) {
	for _, n := range specs {
		// 判断是否有效
		vv, ok := n.(*ast.ValueSpec)
		if !ok || vv == nil {
			continue
		}
		// 解析枚举类型
		name := vv.Type.(*ast.Ident).Name
		if _, ok := d.enums[name]; !ok {
			d.enums[name] = &domain.Enum{
				PkgName: pkgName,
				Name:    name,
				Fields:  make(map[string]*domain.Value),
			}
		}
		// 保存字段
		d.enums[name].Fields[vv.Names[0].Name] = &domain.Value{
			Comment: base.ParseComment(vv.Comment),
			Value:   cast.ToInt32(vv.Values[0].(*ast.BasicLit).Value),
		}
	}
}
