package typespec

import (
	"fmt"
	"sort"
	"strings"
	"universal/tools/gomaker_new/domain"
)

type Type struct {
	Kind int32  // 数据类型
	Pkg  string // 包名
	Name string // 字段名
}

func TYPE(k int32, pkg, name string) *Type {
	return &Type{Kind: k, Pkg: pkg, Name: name}
}

func (d *Type) GetType(pkg string) string {
	if len(d.Pkg) <= 0 || d.Pkg == pkg {
		return d.Name
	}
	return fmt.Sprintf("%s.%s", d.Pkg, d.Name)
}

type Alias struct {
	Type      *Type   // 别名类型
	RealToken []int32 // 结构类型
	RealType  *Type   // 引用类型
	Class     string  // 分类
	Doc       string  // 注释
}

func ALIAS(t, r *Type, class, doc string, ts ...int32) *Alias {
	return &Alias{Type: t, RealType: r, Class: class, Doc: doc, RealToken: ts}
}

func getToken(pkg string, tt *Type, ts ...int32) string {
	ls := []string{}
	for _, tname := range ts {
		switch tname {
		case domain.TokenTypePointer:
			ls = append(ls, "*")
		case domain.TokenTypeArray:
			ls = append(ls, "[]")
		}
	}
	return fmt.Sprintf("%s%s", strings.Join(ls, ""), tt.GetType(pkg))
}

func (d *Alias) String() string {
	return fmt.Sprintf("%s\ntype %s %s", getDoc(d.Doc), d.Type.Name, getToken(d.Type.Pkg, d.RealType, d.RealToken...))
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
	Type   *Type             // 类型
	Values map[string]*Value // 字段
	List   []*Value          // 排序队列
	Class  string            // 分类
	Doc    string            // 注释
}

func ENUM(t *Type, class, doc string, vs ...*Value) *Enum {
	tmp := make(map[string]*Value)
	for _, v := range vs {
		tmp[v.Name] = v
	}
	sort.Slice(vs, func(i, j int) bool { return vs[i].Value < vs[j].Value })
	return &Enum{Type: t, Values: tmp, List: vs, Doc: doc, Class: class}
}

func (d *Enum) Set(class, doc string) *Enum {
	d.Class = class
	d.Doc = doc
	return d
}

func (d *Enum) AddValue(name string, val int32, doc string) {
	item := &Value{Type: d.Type, Name: name, Value: val, Doc: doc}
	d.List = append(d.List, item)
	d.Values[name] = item
}

func (d *Enum) Proto() string {
	tmps := []string{}
	for _, val := range d.List {
		tmps = append(tmps, fmt.Sprintf("\t%s\t=\t%d;\t%s", val.Name, val.Value, getDoc(val.Doc)))
	}
	return fmt.Sprintf("%s\nenum %s {\n%s\n}", getDoc(d.Doc), d.Type.Name, strings.Join(tmps, "\n"))
}

func (d *Enum) String() string {
	tmps := []string{}
	for _, val := range d.List {
		tmps = append(tmps, fmt.Sprintf("\t%s_%s\t%s = %d\t%s", d.Type.Name, val.Name, d.Type.Name, val.Value, getDoc(val.Doc)))
	}
	return fmt.Sprintf("%s\ntype %s int32\nconst (\n%s\n)", getDoc(d.Doc), d.Type.Name, strings.Join(tmps, "\n"))
}

type Field struct {
	Token []int32 // 结构类型
	Type  *Type   // 类型
	Name  string  // 字段名字
	Index int     // 下标
	Tag   string  // 标签
	Doc   string  // 注释
}

func (d *Field) GetToken() string {
	return getToken(d.Type.Pkg, d.Type, d.Token...)
}

type Struct struct {
	Type   *Type             // 类型
	Fields map[string]*Field // 字段
	List   []*Field          // 排序队列
	Class  string            // 分类
	Doc    string            // 注释
}

func STRUCT(t *Type, class, doc string, vs ...*Field) *Struct {
	tmp := make(map[string]*Field)
	for _, val := range vs {
		tmp[val.Name] = val
	}
	return &Struct{Type: t, Fields: tmp, List: vs, Doc: doc, Class: class}
}

func (d *Struct) Proto() string {
	tmps := []string{}
	for _, val := range d.List {
		tmps = append(tmps, fmt.Sprintf("\t%s\t%s\t=\t%d;\t%s", val.Type.GetType(d.Type.Name), val.Name, val.Index+1, getDoc(val.Doc)))
	}
	return fmt.Sprintf("%s\nmessage %s {\n%s\n}", getDoc(d.Doc), d.Type.Name, strings.Join(tmps, "\n"))
}

func (d *Struct) String() string {
	tmps := []string{}
	for _, val := range d.List {
		tmps = append(tmps, fmt.Sprintf("\t%s\t%s\t%s\t%s", val.Name, getToken(d.Type.Pkg, val.Type, val.Token...), val.Tag, getDoc(val.Doc)))
	}
	return fmt.Sprintf("%s\ntype %s struct {\n%s\n}", getDoc(d.Doc), d.Type.Name, strings.Join(tmps, "\n"))
}

func getDoc(doc string) string {
	if len(doc) > 0 {
		return fmt.Sprintf("// %s", doc)
	}
	return doc
}
