package domain

import (
	"text/template"
)

const (
	KIND_IDENT    = 0
	KIND_ENUM     = 1
	KIND_ALIAS    = 2
	KIND_STRUCT   = 3
	TOKEN_POINTER = 1
	TOKEN_ARRAY   = 1 << 1
	TOKEN_MAP     = 1 << 2
)

// 代码生成
type GenFunc func(dst string, tpls *template.Template, extra ...string) error
type ConvFunc func(string) interface{}
