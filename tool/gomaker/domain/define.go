package domain

import "text/template"

// 基础类型
const (
	IDENT  = 1
	STRUCT = 2
	ALIAS  = 3
)

const (
	POINTER = 1
	ARRAY   = 1 << 1
	MAP     = 1 << 2
)

// 数据类型
type Type struct {
	Kind     uint32 // 基础类型
	Selector string // 引用的包名
	Name     string // 字段名称
	Doc      string // 注释
}

type Value struct {
	Name    string // 字段名字
	Type    *Type  // 字段类型
	Value   int32  // 字段值
	Comment string // 注释
}

// 字段结构
type Field struct {
	Token   uint32 // 数据类型
	Name    string // 字段名字
	Type    *Type  // 类型
	Tag     string // 标签
	Comment string // 注释
}

type Alias struct {
	Token     uint32 // 数据类型
	AliasType *Type  // 别名类型
	Type      *Type  // 真实类型
	Comment   string // 注释
}

type Struct struct {
	Type   *Type             // 类型
	Fields map[string]*Field // 字段
	List   []*Field          // 排序队列
}

// 代码生成接口
type GenFunc func(string, string, map[string]*template.Template) error
