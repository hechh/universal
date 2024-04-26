package typespec

import (
	"fmt"
	"go/ast"
	"universal/tools/gomaker/domain"
)

type Type struct {
	BaseFunc
	Token   int32  // 类型
	PkgName string // 引用类型所在的包
	Name    string // 引用的类型名称
	Key     *Type  // map特殊处理（key类型）
	Value   *Type  // map特殊处理（value类型）
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
		typ.Key = &Type{}
		typ.Value = &Type{}
		parseType(v.Key, typ.Key)
		parseType(v.Value, typ.Value)
	case *ast.SelectorExpr:
		typ.AddToken(domain.IDENT)
		typ.PkgName = v.X.(*ast.Ident).Name
		typ.Name = v.Sel.Name
	case *ast.StarExpr:
		typ.AddToken(domain.POINTER)
		parseType(v.X, typ)
	case *ast.ArrayType:
		typ.AddToken(domain.ARRAY)
		parseType(v.Elt, typ)
	case *ast.Ident:
		typ.AddToken(domain.IDENT)
		typ.Name = v.Name
	}
}

func (d *Type) AddToken(val int32) {
	d.Token <<= 4
	d.Token |= val & 0x0f
}

func (d *Type) GetName(pkg string) string {
	if len(d.PkgName) > 0 && pkg != d.PkgName {
		return d.PkgName + "." + d.Name
	}
	return d.Name
}

func (d *Type) GetTypeString(pkg string) (str string) {
	for i := 7; i >= 0; i-- {
		switch (d.Token >> (i * 4)) & 0x0f {
		case domain.IDENT:
			str += d.GetName(pkg)
		case domain.POINTER:
			str += "*"
		case domain.ARRAY:
			str += "[]"
		case domain.MAP:
			str += fmt.Sprintf("map[%s]%s", d.Key.GetTypeString(pkg), d.Value.GetTypeString(pkg))
		}
	}
	return
}

func (d *Type) String() string {
	return d.GetTypeString("")
}

// 判单是否为基础数据类型
func (d *Type) IsMap() bool {
	for i := 7; i >= 0; i-- {
		if val := (d.Token >> (i * 4)) & 0x0f; val > 0 {
			return (val & domain.MAP) == domain.MAP
		}
	}
	return false
}

func (d *Type) IsArray() bool {
	for i := 7; i >= 0; i-- {
		if val := (d.Token >> (i * 4)) & 0x0f; val > 0 {
			return (val & domain.ARRAY) == domain.ARRAY
		}
	}
	return false
}

func (d *Type) IsCustom() bool {
	return len(d.PkgName) > 0
}
