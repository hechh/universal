package types

import (
	"fmt"
	"go/ast"
	"sort"
	"strings"
	"unicode"
	"universal/tools/gomaker/internal/base"
)

type Struct struct {
	base.BaseFunc
	PkgName string            // 所在包名
	Name    string            // 引用的类型名称
	Doc     string            // 注释
	Fields  map[string]*Field // 解析的字段
	List    []*Field          // 排序
}

func ParseComment(g *ast.CommentGroup) string {
	if g == nil || len(g.List) <= 0 {
		return ""
	}
	return g.List[0].Text
}

func NewStruct(pkg, Name, doc string, Fields []*ast.Field) *Struct {
	item := &Struct{
		PkgName: pkg,
		Name:    Name,
		Doc:     doc,
		Fields:  make(map[string]*Field),
		List:    make([]*Field, 0),
	}
	for _, field := range Fields {
		// 过滤不对外开放接
		if unicode.IsLower(rune(field.Names[0].Name[0])) {
			continue
		}
		ff := NewField(field.Names[0].Name, ParseComment(field.Comment), NewType(field.Type), field.Tag)
		item.Fields[field.Names[0].Name] = ff
		item.List = append(item.List, ff)
	}
	sort.Slice(item.List, func(i, j int) bool {
		return strings.Compare(item.List[i].Name, item.List[j].Name) < 0
	})
	return item
}

func (d *Struct) GetTypeString(pkg string) string {
	if len(d.PkgName) > 0 || pkg != d.PkgName {
		return d.PkgName + "." + d.Name
	}
	return d.Name
}

func (d *Struct) String() string {
	str := ""
	for _, f := range d.Fields {
		str += f.String()
		str += "\n"
	}
	return fmt.Sprintf("type %s struct {\n %s }", d.Name, str)
}
