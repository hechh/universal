package typespec

import (
	"fmt"
	"sort"
	"strings"
	"universal/tools/gomaker/domain"
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

func (d *Field) GetToken() string {
	strToken := ""
	for _, val := range d.Token {
		switch val {
		case domain.POINTER:
			strToken += "*"
		case domain.ARRAY:
			strToken += "[]"
		}
	}
	return strToken
}

func (d *Field) GetComment() string {
	if len(d.Comment) > 0 {
		return fmt.Sprintf("// %s", d.Comment)
	}
	return ""
}

func (d *Field) Clone() *Field {
	tmps := make([]uint32, len(d.Token))
	copy(tmps, d.Token)
	return &Field{tmps, d.Name, d.Type, d.Tag, d.Comment}
}

func (d *Struct) Format() string {
	vals := []string{}
	for _, val := range d.List {
		vals = append(vals, fmt.Sprintf("%s %s%s %s %s", val.Name, val.GetToken(), val.Type.GetName(d.Type.Selector), val.Tag, val.GetComment()))
	}
	return fmt.Sprintf("type %s struct {\n%s\n}", d.Type.Name, strings.Join(vals, "\n"))
}

func (d *Struct) Clone() *Struct {
	tmps := make(map[string]*Field)
	list := []*Field{}
	for _, ff := range d.List {
		item := ff.Clone()
		list = append(list, item)
		tmps[item.Name] = item
	}
	sort.Slice(list, func(i, j int) bool {
		return strings.Compare(list[i].Name, list[j].Name) < 0
	})
	return &Struct{Type: d.Type.Clone(), Fields: tmps, List: list}
}
