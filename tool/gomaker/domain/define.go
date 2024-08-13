package domain

import "text/template"

// 代码生成接口
type GenFunc func(string, string, map[string]*template.Template) error

// --------------------基础数据类型---------------------
type Value struct {
	Name    string
	Value   int32
	Comment string
}

type Type struct {
	Token   int32  // 基础类型
	PkgName string // 解析文件的包名
	RefName string // 引用字段的包名
	Name    string // 字段名称
}

type Enum struct {
	Type   *Type             // 类型
	Doc    string            // 注释
	Fields map[string]*Value // 字段
	List   []*Value          // 排序队列
}

type Field struct {
	Type    *Type  // 类型
	Name    string // 字段名字
	Tag     string // 标签
	Comment string // 注释
}

type Struct struct {
	Type   *Type             // 类型
	Doc    string            // 注释
	Fields map[string]*Field // 字段
	List   []*Field          // 排序队列
}

type Alias struct {
	Type *Type  // 类型
	Name string // 别名
	Doc  string // 注释
}

const (
	ENUM   = 1
	STRUCT = 2
)
