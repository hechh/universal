package types

import (
	"fmt"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/common/base"
)

type Type struct {
	base.BaseFunc
	Token   int32  // 类型
	PkgName string // 引用类型所在的包
	Name    string // 引用的类型名称
	Key     *Type  // map特殊处理（key类型）
	Value   *Type  // map特殊处理（value类型）
}

type Field struct {
	base.BaseFunc
	Type    *Type  // 类型
	Name    string // 字段名字
	Comment string // 注释
	Tag     string // 标签
}

type Struct struct {
	base.BaseFunc
	Type   *Type             // 类型
	Doc    string            // 注释
	Fields map[string]*Field // 解析的字段
	List   []*Field          // 排序
}

type Value struct {
	Name    string
	Value   int32
	Comment string
}

type Enum struct {
	base.BaseFunc
	Type   *Type             // 类型
	Doc    string            // 注释规则
	Fields map[string]*Value // 解析的字段
	List   []*Value          // 排序队列
}

type Alias struct {
	base.BaseFunc
	Type      *Type  // 别名类型
	Doc       string // 规则注释
	Comment   string // 字段注释
	Reference *Type  // 引用类型
}

func (d *Type) GetType(pkg string) string {
	if len(d.PkgName) > 0 && len(pkg) > 0 && pkg != d.PkgName {
		return fmt.Sprintf("%s.%s", d.PkgName, d.Name)
	}
	return d.Name
}

func (d *Type) String(pkg string) (ret string) {
	if domain.MAP&d.Token > 0 {
		return fmt.Sprintf("map[%s]%s", d.Key.String(pkg), d.Value.String(pkg))
	}
	if domain.ARRAY&d.Token > 0 {
		ret += "[]"
	}
	if domain.POINTER&d.Token > 0 {
		ret += "*"
	}
	ret += d.GetType(pkg)
	return
}

func (d *Field) String(pkg string) string {
	return fmt.Sprintf("%s %s %s %s", d.Name, d.Type.String(pkg), d.Tag, d.Comment)
}

func (d *Struct) String(pkg string) string {
	str := "\n"
	for _, item := range d.List {
		str += (item.String(pkg) + "\n")
	}
	return fmt.Sprintf("type %s struct {%s}", d.Type.Name, str)
}
