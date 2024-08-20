package typespec

import (
	"fmt"
	"sort"
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

func (d *Value) GetComment() string {
	if len(d.Comment) > 0 {
		return fmt.Sprintf("// %s", d.Comment)
	}
	return ""
}

func (d *Value) Clone() *Value {
	return &Value{d.Name, d.Type.Clone(), d.Value, d.Comment}
}

func (d *Enum) Format() string {
	vals := []string{}
	for _, val := range d.List {
		vals = append(vals, fmt.Sprintf("%s %s = %d %s", val.Name, val.Type.GetName(d.Type.Selector), val.Value, val.GetComment()))
	}
	return fmt.Sprintf("type %s int32\nconst(\n%s\n)", d.Type.Name, strings.Join(vals, "\n"))
}

func (d *Enum) Clone() *Enum {
	tmps := make(map[string]*Value)
	list := []*Value{}
	for _, vv := range d.List {
		item := vv.Clone()
		list = append(list, item)
		tmps[item.Name] = item
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Value < list[j].Value
	})
	return &Enum{d.Type.Clone(), tmps, list}
}
