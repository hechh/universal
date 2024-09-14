package typespec

import (
	"fmt"
	"strings"
	"unicode"
	"universal/tools/gomaker_new/domain"
)

type Field struct {
	Token []uint32 // 数据类型
	Type  *Type    // 类型
	Name  string   // 字段名字
	Tag   string   // 标签
	Doc   string   // 注释
	Index int      // 下标
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

func (d *Field) Format(tt *Type) string {
	return fmt.Sprintf("%s %s%s %s %s", d.Name, d.GetToken(), d.Type.GetName(tt.PkgName), d.Tag, d.GetDoc())
}

type Struct struct {
	Type   *Type             // 类型
	Fields map[string]*Field // 字段
	List   []*Field          // 排序队列
	Doc    string
}

func (d *Struct) Clone() *Struct {
	ret := &Struct{Type: d.Type, Doc: d.Doc}
	for _, item := range d.List {
		if unicode.IsLower(rune(item.Name[0])) {
			continue
		}
		ret.List = append(ret.List, item)
	}
	return ret
}

func (d *Struct) Add(ff *Field) {
	if _, ok := d.Fields[ff.Name]; !ok {
		d.Fields[ff.Name] = ff
		d.List = append(d.List, ff)
	}
}

func (d *Struct) GetDoc() string {
	if len(d.Doc) <= 0 {
		return ""
	}
	return fmt.Sprintf("// %s", d.Doc)
}

func (d *Struct) Format() string {
	vals := []string{}
	for _, val := range d.List {
		vals = append(vals, val.Format(d.Type))
	}
	return fmt.Sprintf("%s\ntype %s struct {\n%s\n}", d.GetDoc(), d.Type.Name, strings.Join(vals, "\n"))
}
