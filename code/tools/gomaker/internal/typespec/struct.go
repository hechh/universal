package typespec

import (
	"fmt"
	"strings"
	"universal/tools/gomaker/domain"
)

type Field struct {
	Type  *Type    // 类型
	Name  string   // 字段名字
	Index int      // 字段下标
	Token []uint32 // 数据类型
	Tag   string   // 标签
	Doc   string   // 注释
}

func (d *Field) GetToken() string {
	strToken := ""
	for _, val := range d.Token {
		switch val {
		case domain.TOKEN_POINTER:
			strToken += "*"
		case domain.TOKEN_ARRAY:
			strToken += "[]"
		}
	}
	return strToken
}

func (d *Field) GetDoc() string {
	if len(d.Doc) <= 0 {
		return ""
	}
	return fmt.Sprintf("// %s", d.Doc)
}

type Struct struct {
	Type     *Type             // 类型
	FileName string            // 文件名
	Fields   map[string]*Field // 字段
	List     []*Field          // 排序队列
	Doc      string
}

func NewStruct(tt *Type, filename string) *Struct {
	return &Struct{Type: tt, FileName: filename, Fields: make(map[string]*Field)}
}

func (d *Struct) GetDoc() string {
	if len(d.Doc) <= 0 {
		return ""
	}
	return fmt.Sprintf("// %s", d.Doc)
}

func (d *Struct) Add(ff *Field) *Struct {
	if _, ok := d.Fields[ff.Name]; !ok {
		d.Fields[ff.Name] = ff
		d.List = append(d.List, ff)
	}
	return d
}

func (d *Struct) Format() string {
	vals := []string{}
	for _, val := range d.List {
		vals = append(vals, fmt.Sprintf("%s %s%s %s %s", val.Name, val.GetToken(), val.Type.GetName(d.Type.PkgName), val.Tag, val.GetDoc()))
	}
	return fmt.Sprintf("%s\ntype %s struct {\n%s\n}", d.GetDoc(), d.Type.Name, strings.Join(vals, "\n"))
}
