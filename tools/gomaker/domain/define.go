package domain

import "go/ast"

// token类型
const (
	IDENT   = 0x01
	POINTER = 0x02
	ARRAY   = 0x04
	MAP     = 0x08
)

type IParser interface {
	AddType(pkgName string, specs []ast.Spec)
	AddConst(pkgName string, specs []ast.Spec)
}

type Type struct {
	Token   int32  // 类型
	Name    string // 引用的类型名称
	PkgName string // 引用类型所在的包
	Value   *Type  // map特殊处理（value类型）
}

type Field struct {
	Name    string
	Comment string
	Type    *Type
}

type Value struct {
	Value   int32
	Comment string
}

type Struct struct {
	PkgName string            // 所在包名
	Name    string            // 引用的类型名称
	Fields  map[string]*Field // 解析的字段
}
type Enum struct {
	PkgName string            // 所在包名
	Name    string            // 引用的类型名称
	Fields  map[string]*Value // 解析的字段
}

type Alias struct {
	PkgName string // 所在包名
	Name    string // 引用的类型名称
	Type    *Type  // 类型
	Comment string
}
