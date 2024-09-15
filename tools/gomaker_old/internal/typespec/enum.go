package typespec

import (
	"fmt"
	"strings"
)

type Value struct {
	Type  *Type  // 字段类型
	Name  string // 字段名字
	Value int32  // 字段值
	Doc   string // 注释
}

func (d *Value) GetDoc() string {
	if len(d.Doc) <= 0 {
		return ""
	}
	return fmt.Sprintf("// %s", d.Doc)
}

func (d *Value) Format() string {
	return fmt.Sprintf("%s %s = %d %s", d.Name, d.Type.GetName(d.Type.PkgName), d.Value, d.GetDoc())
}

type Enum struct {
	Type   *Type             // 类型
	Values map[string]*Value // 字段
	List   []*Value          // 排序队列
	Doc    string            // 注释
}

func (d *Enum) Add(tt *Type, name string, value int32, doc string) {
	item := &Value{Type: tt, Name: name, Value: value, Doc: doc}
	d.Values[name] = item
	d.List = append(d.List, item)
	d.Type = tt
}

func (d *Enum) GetDoc() string {
	if len(d.Doc) <= 0 {
		return ""
	}
	return fmt.Sprintf("// %s", d.Doc)
}

func (d *Enum) Format() string {
	strs := []string{}
	for _, val := range d.List {
		strs = append(strs, val.Format())
	}
	return fmt.Sprintf("type %s int32\nconst(\n%s\n)", d.Type.GetName(d.Type.PkgName), strings.Join(strs, "\n"))
}
