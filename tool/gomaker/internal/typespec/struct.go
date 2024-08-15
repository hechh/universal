package typespec

import (
	"fmt"
	"strings"
)

// 字段结构
type Field struct {
	Token   []uint32 // 数据类型
	Name    string   // 字段名字
	Type    *Type    // 类型
	Tag     string   // 标签
	Comment string   // 注释
}

type Struct struct {
	Type   *Type             // 类型
	Fields map[string]*Field // 字段
	List   []*Field          // 排序队列
}

func (d *Struct) Format() string {
	vals := []string{}
	for _, val := range d.List {
		vals = append(vals, fmt.Sprintf("%s %s // %s %s", val.Name, val.Type.Format(d.Type.Selector), val.Tag, val.Comment))
	}
	return fmt.Sprintf("type %s struct {\n%s\n}", d.Type.Name, strings.Join(vals, "\n"))
}
