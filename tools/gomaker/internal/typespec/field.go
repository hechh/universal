package typespec

import "fmt"

type Field struct {
	BaseFunc
	Name    string // 字段名字
	Comment string // 注释
	Type    *Type  // 类型
}

func NewField(Name, Comment string, typ *Type) *Field {
	return &Field{
		Name:    Name,
		Comment: Comment,
		Type:    typ,
	}
}

func (d *Field) GetType(pkg string) string {
	return d.Type.GetType(pkg)
}

func (d *Field) String() string {
	return fmt.Sprintf("\t%s %s %s", d.Name, d.GetType(""), d.Comment)
}
