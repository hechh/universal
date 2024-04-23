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
