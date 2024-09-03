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

func (d *Value) GetType() string {
	return d.Type.GetName(d.Type.PkgName)
}

type Enum struct {
	Type     *Type             // 类型
	FileName string            // 文件名
	Values   map[string]*Value // 字段
	List     []*Value          // 排序队列
	Doc      string
}

func NewEnum(tt *Type, filename string) *Enum {
	return &Enum{Type: tt, FileName: filename, Values: make(map[string]*Value)}
}

func (d *Enum) GetDoc() string {
	if len(d.Doc) <= 0 {
		return ""
	}
	return fmt.Sprintf("// %s", d.Doc)
}

func (d *Enum) GetType() string {
	return d.Type.GetName(d.Type.PkgName)
}

func (d *Enum) Add(val *Value) *Enum {
	if _, ok := d.Values[val.Name]; !ok {
		d.Values[val.Name] = val
		d.List = append(d.List, val)
	}
	return d
}

func (d *Enum) Format() string {
	strs := []string{}
	for _, val := range d.List {
		strs = append(strs, fmt.Sprintf("%s %s = %d %s", val.Name, val.GetType(), val.Value, val.GetDoc()))
	}
	return fmt.Sprintf("type %s int32\nconst(\n%s\n)", d.GetType(), strings.Join(strs, "\n"))
}
