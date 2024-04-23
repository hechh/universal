package typespec

import (
	"fmt"
	"go/ast"
)

type Alias struct {
	BaseFunc
	PkgName string // 所在包名
	Name    string // 引用的类型名称
	Type    *Type  // 类型
	Comment string
}

func NewAlias(pkg, Name string, node *ast.TypeSpec) *Alias {
	return &Alias{
		PkgName: pkg,
		Name:    Name,
		Comment: parseComment(node.Comment),
		Type:    NewType(node.Type),
	}
}

func (d *Alias) GetType(pkg string) string {
	if len(d.PkgName) > 0 && pkg != d.PkgName {
		return d.PkgName + "." + d.Name
	}
	return d.Name
}

func (d *Alias) String() string {
	return fmt.Sprintf("type %s %s %s", d.Name, d.Type.String(), d.Comment)
}
