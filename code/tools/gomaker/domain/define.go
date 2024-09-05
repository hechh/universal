package domain

import (
	"go/ast"
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

const (
	PACKAGE = "package"
	CONFIG  = "xlsx"
	JSON    = "json"
	HTTPKIT = "httpkit.tpl"
	PBCLASS = "pbclass.tpl"
	PROTO   = "proto.tpl"
)

// 代码生成接口
type GenFunc func(string, string, *template.Template) error

type IParser interface {
	SetFile(string)
	Visit(ast.Node) ast.Visitor
}
