package typespec

import (
	"fmt"
	"sort"
	"strings"
	"universal/tools/pbtool/domain"
)

type Type struct {
	Kind   int32  // 数据类型
	Pkg    string // 包名
	Name   string // 字段名
	Class  string // 分类
	Source int32  // 来源
}

func TYPE(kind, source int32, pkg, name, class string) *Type {
	return &Type{Kind: kind, Pkg: pkg, Name: name, Class: class, Source: source}
}

func (d *Type) GetPkgType() string {
	if len(d.Pkg) <= 0 {
		return d.Name
	}
	return fmt.Sprintf("%s.%s", d.Pkg, d.Name)
}

func (d *Type) GetType(pkg string) string {
	if len(d.Pkg) <= 0 || pkg == d.Pkg {
		return d.Name
	}
	return fmt.Sprintf("%s.%s", d.Pkg, d.Name)
}

type Alias struct {
	Type      *Type  // 别名类型
	RealType  *Type  // 引用类型
	RealToken []rune // 结构类型
	Doc       string // 注释
}

func ALIAS(t, r *Type, doc string, ts ...int32) *Alias {
	return &Alias{Type: t, RealType: r, Doc: doc, RealToken: ts}
}

func (d *Alias) GetToken(pkg string) string {
	ls := []string{}
	for _, tname := range d.RealToken {
		switch tname {
		case domain.TokenTypePointer:
			ls = append(ls, "*")
		case domain.TokenTypeArray:
			ls = append(ls, "[]")
		}
	}
	return fmt.Sprintf("%s%s", strings.Join(ls, ""), d.Type.GetType(pkg))
}

type Value struct {
	Type  *Type  // 字段类型
	Name  string // 字段名字
	Value int32  // 字段值
	Doc   string // 注释
}

func VALUE(t *Type, name string, v int32, doc string) *Value {
	return &Value{Type: t, Name: name, Value: v, Doc: doc}
}

type Enum struct {
	Type    *Type             // 类型
	Members map[string]*Value // 字段
	List    []*Value          // 排序队列
	Doc     string            // 注释
}

func ENUM(t *Type, doc string, vs ...*Value) *Enum {
	tmp := make(map[string]*Value)
	for _, v := range vs {
		tmp[v.Name] = v
	}
	sort.Slice(vs, func(i, j int) bool { return vs[i].Value < vs[j].Value })
	return &Enum{Type: t, Members: tmp, List: vs, Doc: doc}
}

func (d *Enum) Set(doc string) *Enum {
	d.Doc = doc
	return d
}

func (d *Enum) Add(name string, val int32, doc string) {
	item := &Value{Type: d.Type, Name: name, Value: val, Doc: doc}
	d.List = append(d.List, item)
	d.Members[name] = item
}

func (d *Enum) AddValue(item *Value) {
	d.List = append(d.List, item)
	d.Members[item.Name] = item
}

type Field struct {
	Token []rune // 结构类型
	Type  *Type  // 类型
	Name  string // 字段名字
	Index int    // 下标
	Tag   string // 标签
	Doc   string // 注释
}

func FIELD(t *Type, name string, index int, tag, doc string, ts ...rune) *Field {
	return &Field{Type: t, Name: name, Index: index, Tag: tag, Doc: doc, Token: ts}
}

func (d *Field) GetToken(pkg string) string {
	ls := []string{}
	for _, tname := range d.Token {
		switch tname {
		case domain.TokenTypePointer:
			ls = append(ls, "*")
		case domain.TokenTypeArray:
			ls = append(ls, "[]")
		}
	}
	return fmt.Sprintf("%s%s", strings.Join(ls, ""), d.Type.GetType(pkg))
}

type Struct struct {
	Type    *Type             // 类型
	Members map[string]*Field // 字段
	List    []*Field          // 排序队列
	Doc     string            // 注释
}

func STRUCT(t *Type, doc string, vs ...*Field) *Struct {
	tmp := make(map[string]*Field)
	for _, val := range vs {
		tmp[val.Name] = val
	}
	return &Struct{Type: t, Members: tmp, List: vs, Doc: doc}
}

func (d *Struct) Add(t *Type, name string, index int, tag, doc string, ts ...int32) {
	val := &Field{Type: t, Name: name, Index: index, Tag: tag, Doc: doc, Token: ts}
	d.List = append(d.List, val)
	d.Members[val.Name] = val
}
