package typespec

import "fmt"

// 数据类型
type Type struct {
	Kind     uint32 // 基础类型
	Selector string // 引用的包名
	Name     string // 字段名称
	Doc      string // 注释
}

func (d *Type) Format(pkg string) string {
	if len(d.Selector) <= 0 || d.Selector == pkg {
		return d.Name
	}
	return fmt.Sprintf("%s.%s", d.Selector, d.Name)
}
