package base

import (
	"fmt"
	"go/ast"
)

type Struct struct {
	pkgName string            // 所在包名
	name    string            // 引用的类型名称
	fields  map[string]*Field // 解析的字段
}

func parseComment(g *ast.CommentGroup) string {
	if g == nil || len(g.List) <= 0 {
		return ""
	}
	return g.List[0].Text
}

func NewStruct(pkg, name string, fields []*ast.Field) *Struct {
	item := &Struct{pkg, name, make(map[string]*Field)}
	for _, field := range fields {
		item.fields[field.Names[0].Name] = NewField(field.Names[0].Name, parseComment(field.Comment), NewType(field.Type))
	}
	return item
}

func (d *Struct) GetType(pkg string) string {
	if len(d.pkgName) > 0 || pkg != d.pkgName {
		return d.pkgName + "." + d.name
	}
	return d.name
}

func (d *Struct) String() string {
	str := ""
	for _, f := range d.fields {
		str += f.String()
		str += "\n"
	}
	return fmt.Sprintf("type %s struct {\n %s }", d.name, str)
}
