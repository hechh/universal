package base

import (
	"fmt"
	"go/ast"
)

type Alias struct {
	pkgName  string // 所在包名
	name     string // 引用的类型名称
	realType *Type  // 类型
	comment  string
}

func NewAlias(pkg, name string, node *ast.TypeSpec) *Alias {
	return &Alias{
		pkgName:  pkg,
		name:     name,
		comment:  parseComment(node.Comment),
		realType: NewType(node.Type),
	}
}

func (d *Alias) GetType(pkg string) string {
	if len(d.pkgName) > 0 && pkg != d.pkgName {
		return d.pkgName + "." + d.name
	}
	return d.name
}

func (d *Alias) String() string {
	return fmt.Sprintf("type %s %s %s", d.name, d.realType.String(), d.comment)
}
