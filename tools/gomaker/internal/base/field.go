package base

import "fmt"

type Field struct {
	name      string // 字段名字
	comment   string // 注释
	fieldType *Type  // 类型
}

func NewField(name, comment string, typ *Type) *Field {
	return &Field{name, comment, typ}
}

func (d *Field) GetName() string {
	return d.name
}

func (d *Field) GetComment() string {
	return d.comment
}

func (d *Field) GetType(pkg string) string {
	return d.fieldType.GetType(pkg)
}

func (d *Field) String() string {
	return fmt.Sprintf("\t%s %s %s", d.name, d.GetType(""), d.comment)
}
