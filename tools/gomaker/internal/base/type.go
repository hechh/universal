package base

import (
	"fmt"
	"go/ast"
	"universal/tools/gomaker/domain"
)

type Type struct {
	token   int32  // 类型
	name    string // 引用的类型名称
	pkgName string // 引用类型所在的包
	key     *Type  // map特殊处理（key类型）
	value   *Type  // map特殊处理（value类型）
}

func NewType(val ast.Expr) *Type {
	item := &Type{}
	parseType(val, item)
	return item
}

func parseType(val ast.Expr, typ *Type) {
	switch v := val.(type) {
	case *ast.MapType:
		typ.AddToken(domain.MAP)
		typ.key = &Type{}
		typ.value = &Type{}
		parseType(v.Key, typ.key)
		parseType(v.Value, typ.value)
	case *ast.SelectorExpr:
		typ.pkgName = v.X.(*ast.Ident).Name
		typ.name = v.Sel.Name
	case *ast.StarExpr:
		typ.AddToken(domain.POINTER)
		parseType(v.X, typ)
	case *ast.ArrayType:
		typ.AddToken(domain.ARRAY)
		parseType(v.Elt, typ)
	case *ast.Ident:
		typ.AddToken(domain.IDENT)
		typ.name = v.Name
	}
}

func (d *Type) AddToken(val int32) {
	d.token <<= 4
	d.token |= val & 0x0f
}

func (d *Type) GetName(pkg string) string {
	if len(d.pkgName) > 0 && pkg != d.pkgName {
		return d.pkgName + "." + d.name
	}
	return d.name
}

func (d *Type) GetType(pkg string) (str string) {
	for i := 7; i >= 0; i-- {
		switch (d.token >> (i * 4)) & 0x0f {
		case domain.IDENT:
			str += d.GetName(pkg)
		case domain.POINTER:
			str += "*"
		case domain.ARRAY:
			str += "[]"
		case domain.MAP:
			str += fmt.Sprintf("map[%s]%s", d.key.GetType(pkg), d.value.GetType(pkg))
		}
	}
	return
}

func (d *Type) String() string {
	return d.GetType("")
}
