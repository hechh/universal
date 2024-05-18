package types

import (
	"fmt"
	"go/ast"
	"universal/tools/gomaker/internal/base"
)

type Field struct {
	base.BaseFunc
	Name    string // 字段名字
	Comment string // 注释
	Type    *Type  // 类型
	Tag     string // 标签
}

func NewField(Name, Comment string, typ *Type, tag *ast.BasicLit) *Field {
	ret := &Field{
		Name:    Name,
		Comment: Comment,
		Type:    typ,
	}
	if tag != nil {
		ret.Tag = tag.Value
	}
	return ret
}

func (d *Field) GetTypeString(pkg string) string {
	return d.Type.GetTypeString(pkg)
}

func (d *Field) String() string {
	return fmt.Sprintf("\t%s %s %s", d.Name, d.GetTypeString(""), d.Comment)
}
