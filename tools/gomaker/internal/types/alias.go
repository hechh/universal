package types

import (
	"fmt"
	"go/ast"
	"universal/tools/gomaker/internal/base"
)

type Alias struct {
	base.BaseFunc
	PkgName string // 所在包名
	Name    string // 引用的类型名称
	Doc     string // 规则注释
	Comment string // 字段注释
	Type    *Type  // 类型
}

func NewAlias(pkg, Name, doc string, node *ast.TypeSpec) *Alias {
	return &Alias{
		PkgName: pkg,
		Name:    Name,
		Doc:     doc,
		Comment: ParseComment(node.Comment),
		Type:    NewType(node.Type),
	}
}

func (d *Alias) GetTypeString(pkg string) string {
	if len(d.PkgName) > 0 && pkg != d.PkgName {
		return d.PkgName + "." + d.Name
	}
	return d.Name
}

func (d *Alias) String() string {
	return fmt.Sprintf("type %s %s %s", d.Name, d.Type.String(), d.Comment)
}
