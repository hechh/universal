package domain

import "text/template"

const (
	IDENT  = 0
	ENUM   = 1
	ALIAS  = 2
	STRUCT = 3
)

const (
	POINTER = 1
	ARRAY   = 1 << 1
	MAP     = 1 << 2
)

const (
	PACKAGE = "package"
	CLIENT  = "client"
)

// 代码生成接口
type GenFunc func(string, string, map[string]*template.Template) error
