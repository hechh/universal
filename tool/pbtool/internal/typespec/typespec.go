package typespec

import (
	"fmt"
	"go/ast"
	"strings"
	"universal/tool/pbtool/domain"
)

type IdentExpr struct {
	tokens []int32
	kind   int32
	name   string
	pkg    string
}

func NewIdentExpr(fac domain.IFactory, pkg string, vv *ast.Ident) *IdentExpr {
	return &IdentExpr{
		kind: fac.GetKind(vv.Name),
		name: vv.Name,
		pkg:  pkg,
	}
}

func (i *IdentExpr) Push(tok int32) {
	i.tokens = append(i.tokens, tok)
}

func (i *IdentExpr) GetKind() int32 {
	return i.kind
}

func (i *IdentExpr) GetName() string {
	return i.name
}

func (i *IdentExpr) GetPkg() string {
	return i.pkg
}

func (i *IdentExpr) FullName(curpkg string) string {
	if len(curpkg) > 0 && i.pkg != curpkg {
		return i.pkg + "." + i.name
	}
	return i.name
}

type SelectorExpr struct {
	tokens []int32
	kind   int32
	name   string
	pkg    string
}

func NewSelectorExpr(fac domain.IFactory, pkg string, vv *ast.SelectorExpr) *SelectorExpr {
	return &SelectorExpr{
		kind: domain.KIND_IDENT,
		name: vv.Sel.Name,
		pkg:  vv.X.(*ast.Ident).Name,
	}
}

func (i *SelectorExpr) Push(tok int32) {
	i.tokens = append(i.tokens, tok)
}

func (i *SelectorExpr) GetKind() int32 {
	return i.kind
}

func (i *SelectorExpr) GetName() string {
	return i.name
}

func (i *SelectorExpr) GetPkg() string {
	return i.pkg
}

func (i *SelectorExpr) FullName(curpkg string) string {
	return i.pkg + "." + i.name
}

type StarExpr struct {
	domain.IType
}

func NewStarExpr(fac domain.IFactory, pkg string, vv *ast.StarExpr) *StarExpr {
	var item domain.IType
	switch nn := vv.X.(type) {
	case *ast.Ident:
		item = NewIdentExpr(fac, pkg, nn)
	case *ast.SelectorExpr:
		item = NewSelectorExpr(fac, pkg, nn)
	default:
		return nil
	}
	item.Push(domain.TOKEN_POINTER)
	return &StarExpr{IType: item}
}

func (i *StarExpr) FullName(curpkg string) string {
	return "*" + i.IType.FullName(curpkg)
}

type ArrayExpr struct {
	domain.IType
}

func NewArrayExpr(fac domain.IFactory, pkg string, vv *ast.ArrayType) *ArrayExpr {
	var item domain.IType
	switch nn := vv.Elt.(type) {
	case *ast.Ident:
		item = NewIdentExpr(fac, pkg, nn)
	case *ast.StarExpr:
		item = NewStarExpr(fac, pkg, nn)
	case *ast.SelectorExpr:
		item = NewSelectorExpr(fac, pkg, nn)
	default:
		return nil
	}
	item.Push(domain.TOKEN_ARRAY)
	return &ArrayExpr{IType: item}
}

func (i *ArrayExpr) FullName(curpkg string) string {
	return "[]" + i.IType.FullName(curpkg)
}

type MapExpr struct {
	domain.IType
	key domain.IType
}

func NewMapExpr(fac domain.IFactory, pkg string, vv *ast.MapType) *MapExpr {
	key := NewIdentExpr(fac, pkg, vv.Key.(*ast.Ident))
	var item domain.IType
	switch nn := vv.Value.(type) {
	case *ast.Ident:
		item = NewIdentExpr(fac, pkg, nn)
	case *ast.StarExpr:
		item = NewStarExpr(fac, pkg, nn)
	case *ast.SelectorExpr:
		item = NewSelectorExpr(fac, pkg, nn)
	case *ast.ArrayType:
		item = NewArrayExpr(fac, pkg, nn)
	default:
		return nil
	}
	item.Push(domain.TOKEN_MAP)
	return &MapExpr{key: key, IType: item}
}

func (i *MapExpr) FullName(curpkg string) string {
	return fmt.Sprintf("map[%s]%s", i.key.FullName(curpkg), i.IType.FullName(curpkg))
}

func ParseType(fac domain.IFactory, pkg string, expr ast.Expr) domain.IType {
	switch vv := expr.(type) {
	case *ast.Ident:
		return NewIdentExpr(fac, pkg, vv)
	case *ast.SelectorExpr:
		return NewSelectorExpr(fac, pkg, vv)
	case *ast.StarExpr:
		return NewStarExpr(fac, pkg, vv)
	case *ast.ArrayType:
		return NewArrayExpr(fac, pkg, vv)
	case *ast.MapType:
		return NewMapExpr(fac, pkg, vv)
	}
	return nil
}

type Attribute struct {
	domain.IType
	name string // 字段名字
}

func NewAttribute(name string, tt domain.IType) *Attribute {
	return &Attribute{IType: tt, name: name}
}

func (a *Attribute) GetName() string {
	return a.name
}

func (a *Attribute) String() string {
	return fmt.Sprintf("\t%s\t%s", a.name, a.FullName(""))
}

type Class struct {
	domain.IType
	fields map[string]domain.IAttribute
	list   []domain.IAttribute
}

func NewClass(tt domain.IType) *Class {
	return &Class{
		IType:  tt,
		fields: make(map[string]domain.IAttribute),
	}
}

func (c *Class) Add(attr domain.IAttribute) {
	if _, ok := c.fields[attr.GetName()]; !ok {
		c.fields[attr.GetName()] = attr
		c.list = append(c.list, attr)
	}
}

func (c *Class) Get(name string) domain.IAttribute {
	return c.fields[name]
}

func (c *Class) GetAll() []domain.IAttribute {
	return c.list
}

func (a *Class) String() string {
	strs := []string{}
	for _, field := range a.list {
		strs = append(strs, field.String())
	}
	return fmt.Sprintf("type\t%s\tstruct\t{\n%s\n}\n", a.GetName(), strings.Join(strs, "\n"))
}
