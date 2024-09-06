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

func NewValue(tt *Type, name string, val int32, doc string) *Value {
	return &Value{tt, name, val, doc}
}

func (d *Value) GetTypeName() string {
	return d.Type.GetName(d.Type.PkgName)
}

func (d *Value) GetDoc() string {
	if len(d.Doc) <= 0 {
		return ""
	}
	return fmt.Sprintf("// %s", d.Doc)
}

func (d *Value) Format() string {
	return fmt.Sprintf("%s %s = %d %s", d.Name, d.GetTypeName(), d.Value, d.GetDoc())
}

type Enum struct {
	FileName string            // 定义所在文件名
	Type     *Type             // 类型
	Values   map[string]*Value // 字段
	List     []*Value          // 排序队列
	Doc      string            // 注释
}

func NewEnum(filename string, tt *Type) *Enum {
	return &Enum{Type: tt, FileName: filename, Values: make(map[string]*Value)}
}

func (d *Enum) GetTypeName() string {
	return d.Type.GetName(d.Type.PkgName)
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
	return fmt.Sprintf("type %s int32\nconst(\n%s\n)", d.GetTypeName(), strings.Join(strs, "\n"))
}

func (d *Enum) Add(val *Value) *Enum {
	if _, ok := d.Values[val.Name]; !ok {
		d.Values[val.Name] = val
		d.List = append(d.List, val)
	}
	if d.Type == nil {
		d.Type = val.Type
	}
	return d
}
