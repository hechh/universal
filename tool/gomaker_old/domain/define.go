package domain

import "text/template"

// 基础类型
const (
	IDENT  = 1
	ENUM   = 2
	STRUCT = 3
	ALIAS  = 4
)

const (
	POINTER = 1
	ARRAY   = 1 << 1
	MAP     = 1 << 2
)

// 代码生成接口
type GenFunc func(string, string, map[string]*template.Template) error

const (
	PACKAGE = "package"
)
