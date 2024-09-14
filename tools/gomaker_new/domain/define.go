package domain

import "text/template"

// 配置表规则
const (
	RuleTypeEnum     = "E:"
	RuleTypeProto    = "@gomaker:proto"
	RuleTypeBytes    = "@gomaker:bytes"
	GenTable         = "生成表"
	DefaultPkg       = "pb"
	DefaultEnumClass = "other"
)

// ast解析类型
const (
	KindTypeIdent    = 0
	KindTypeEnum     = 1
	KindTypeAlias    = 2
	KindTypeStruct   = 3
	TokenTypeNone    = 0
	TokenTypePointer = 1
	TokenTypeArray   = 2
	TokenTypeMap     = 3
)

// 基础数据类型
const (
	UINT32  = "uint32"
	UINT64  = "uint64"
	INT32   = "int32"
	INT64   = "int64"
	BOOL    = "bool"
	FLOAT32 = "float32"
	FLOAT64 = "float64"
	STRING  = "string"
	BYTE    = "byte"
	RUNE    = "rune"
)

var (
	BasicTypes = map[string]struct{}{
		UINT32:  struct{}{},
		UINT64:  struct{}{},
		INT32:   struct{}{},
		INT64:   struct{}{},
		BOOL:    struct{}{},
		FLOAT32: struct{}{},
		FLOAT64: struct{}{},
		STRING:  struct{}{},
		BYTE:    struct{}{},
		RUNE:    struct{}{},
	}
)

// 代码生成
type GenFunc func(dst string, tpls *template.Template, extra ...string) error

// xlsx转bytes
type ConvFunc func(string) interface{}
