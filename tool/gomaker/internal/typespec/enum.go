package typespec

import (
	"fmt"
	"strings"
)

type Value struct {
	Name    string // 字段名字
	Type    *Type  // 字段类型
	Value   int32  // 字段值
	Comment string // 注释
}

type Enum struct {
	Type   *Type             // 类型
	Values map[string]*Value // 字段
	List   []*Value          // 排序队列
}

func (d *Enum) Format() string {
	vals := []string{}
	for _, val := range d.List {
		vals = append(vals, fmt.Sprintf("%s %s = %d // %s", val.Name, val.Type.Format(d.Type.Selector), val.Value, val.Comment))
	}
	return fmt.Sprintf("const(\n%s\n)", strings.Join(vals, "\n"))
}
