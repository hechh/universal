package base

import (
	"go/ast"
	"go/token"
	"universal/tools/creator/domain"

	"github.com/spf13/cast"
)

type Parser struct {
	pkgName      string // 包名
	typeSpecName string // type
}

func (d *Parser) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.File:
		d.pkgName = n.Name.Name
		return d
	case *ast.GenDecl:
		if n.Tok == token.CONST {
			NewEnum(d.pkgName, d.typeSpecName, n.Specs)
			return nil
		}
		if n.Tok == token.TYPE {
			return d
		}
		return nil
	case *ast.TypeSpec:
		d.typeSpecName = n.Name.Name
		switch vv := n.Type.(type) {
		case *ast.StructType:
			NewStruct(d.pkgName, d.typeSpecName, vv.Fields.List)
		}
		return nil
	}
	return nil
}

func getComment(g *ast.CommentGroup) string {
	if g == nil {
		return ""
	}
	return g.List[0].Text
}

func NewEnum(pkg, name string, vals []ast.Spec) *domain.Enum {
	ret := &domain.Enum{
		Type:   &domain.Type{PkgName: pkg, Name: name, Token: domain.ENUM},
		Fields: make(map[string]*domain.EnumValue),
	}
	for _, vv := range vals {
		val := vv.(*ast.ValueSpec)
		ret.Fields[val.Names[0].Name] = &domain.EnumValue{
			Name:    val.Names[0].Name,
			Value:   cast.ToInt32(val.Values[0].(*ast.BasicLit).Value),
			Comment: getComment(val.Comment),
		}
	}
	return ret
}

func parseType(val ast.Expr, typ *domain.Type) {
	switch v := val.(type) {
	case *ast.SelectorExpr:
		typ.PkgName = v.X.(*ast.Ident).Name
		typ.Name = v.Sel.Name
	case *ast.StarExpr:
		typ.Token |= domain.POINTER
		parseType(v.X, typ)
	case *ast.ArrayType:
		typ.Token |= domain.ARRAY
		parseType(v.Elt, typ)
	case *ast.Ident:
		typ.Name = v.Name
	}
}

func NewStruct(pkgName, name string, fields []*ast.Field) *domain.Struct {
	ret := &domain.Struct{Type: &domain.Type{Token: domain.STRUCT, Name: name, PkgName: pkgName}}
	for _, field := range fields {
		switch vv := field.Type.(type) {
		case *ast.MapType:
			k := &domain.Type{Token: domain.MAP}
			parseType(vv.Key, k)
			v := &domain.Type{Token: domain.MAP}
			parseType(vv.Value, v)
			ret.Maps = append(ret.Maps, &domain.MapField{Name: field.Names[0].Name, KType: k, VType: v})
		default:
			v := &domain.Type{}
			parseType(field.Type, v)
			ret.Others = append(ret.Others, &domain.Field{Name: field.Names[0].Name, Type: v})
		}
	}
	return ret
}
