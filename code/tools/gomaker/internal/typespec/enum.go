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

func NewEnum(tt *Type) *Enum {
	return &Enum{Type: tt, Values: make(map[string]*Value)}
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

func (d *Enum) Clone() *Enum {
	ret := NewEnum(d.Type.Clone())
	for _, vv := range d.List {
		ret.AddValue(vv.Clone())
	}
	return ret
}

func (d *Enum) AddValue(val *Value) *Enum {
	if _, ok := d.Values[val.Name]; !ok {
		d.Values[val.Name] = val
		d.List = append(d.List, val)
	}
	return d
}

func (d *Enum) Format() string {
	vals := []string{}
	for _, val := range d.List {
		vals = append(vals, fmt.Sprintf("%s %s = %d %s", val.Name, val.Type.GetName(d.Type.Selector), val.Value, val.GetComment()))
	}
	return fmt.Sprintf("type %s int32\nconst(\n%s\n)", d.Type.Name, strings.Join(vals, "\n"))
}
