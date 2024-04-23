package base

import (
	"fmt"
	"go/ast"

	"github.com/spf13/cast"
)

type Value struct {
	name    string
	value   int32
	comment string
}

type Enum struct {
	pkgName string            // 所在包名
	name    string            // 引用的类型名称
	fields  map[string]*Value // 解析的字段
}

func NewEnum(pkg string, specs []ast.Spec) *Enum {
	item := &Enum{pkgName: pkg, fields: make(map[string]*Value)}
	for _, node := range specs {
		// 判断是否有效
		vv, ok := node.(*ast.ValueSpec)
		if !ok || vv == nil || vv.Type == nil {
			continue
		}
		// 解析枚举类型
		item.name = vv.Type.(*ast.Ident).Name
		// 保存字段
		item.fields[vv.Names[0].Name] = &Value{
			name:    vv.Names[0].Name,
			comment: parseComment(vv.Comment),
			value:   cast.ToInt32(vv.Values[0].(*ast.BasicLit).Value),
		}
	}
	return item
}

func (d *Enum) GetType(pkg string) string {
	if len(d.pkgName) > 0 && pkg != d.pkgName {
		return d.pkgName + "." + d.name
	}
	return d.name
}

func (d *Enum) String() string {
	str := ""
	for _, v := range d.fields {
		str += fmt.Sprintf("\t %s %s = %d // %s\n", v.name, d.name, v.value, v.comment)
	}
	return fmt.Sprintf("const (\n %s )", str)
}
