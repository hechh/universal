package manager

import (
	"go/ast"
	"go/token"
	"strings"
	"universal/tool/gomaker/domain"

	"github.com/spf13/cast"
)

type TypeParser struct {
	pkgName string
	doc     *ast.CommentGroup
}

func (d *TypeParser) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.File:
		d.pkgName = n.Name.Name
		return d
	case *ast.GenDecl:
		d.doc = n.Doc
		if n.Tok != token.CONST && n.Tok != token.TYPE {
			return nil
		}
		return d
	case *ast.ValueSpec:
		switch n.Values[0].(type) {
		case *ast.BasicLit:
			AddValue(getValue(n, d.pkgName, d.doc))
		}
	case *ast.TypeSpec:
		switch n.Type.(type) {
		case *ast.StructType:
			AddStruct(getStruct(n, d.pkgName, d.doc))
		default:
			AddAlias(getAlias(n, d.pkgName, d.doc))
		}
	}
	return nil
}

func getAlias(vv *ast.TypeSpec, pkgName string, doc *ast.CommentGroup) *domain.Alias {
	tt, token := getType(vv.Type, pkgName, nil)
	return &domain.Alias{
		Token: token,
		AliasType: &domain.Type{
			Kind:     domain.ALIAS,
			Selector: pkgName,
			Name:     vv.Name.Name,
			Doc:      getDoc(doc),
		},
		Type:    tt,
		Comment: getDoc(vv.Comment),
	}
}

func getStruct(vv *ast.TypeSpec, pkgName string, doc *ast.CommentGroup) *domain.Struct {
	st := vv.Type.(*ast.StructType)
	tmps := make(map[string]*domain.Field)
	list := []*domain.Field{}
	for _, field := range st.Fields.List {
		tt, token := getType(vv.Type, pkgName, doc)
		ff := &domain.Field{
			Token:   token,
			Name:    field.Names[0].Name,
			Type:    tt,
			Tag:     getTag(field.Tag),
			Comment: getDoc(vv.Comment),
		}
		tmps[ff.Name] = ff
		list = append(list, ff)
	}
	return &domain.Struct{
		Type: &domain.Type{
			Kind:     domain.STRUCT,
			Selector: pkgName,
			Name:     vv.Name.Name,
			Doc:      getDoc(doc),
		},
		Fields: make(map[string]*domain.Field),
		List:   list,
	}
}

func getValue(vv *ast.ValueSpec, pkgName string, doc *ast.CommentGroup) *domain.Value {
	tt, _ := getType(vv.Type, pkgName, doc)
	return &domain.Value{
		Name:    vv.Names[0].Name,
		Type:    tt,
		Value:   cast.ToInt32(vv.Values[0].(*ast.BasicLit).Value),
		Comment: getDoc(vv.Comment),
	}
}

func getTag(tt *ast.BasicLit) string {
	if tt != nil {
		return tt.Value
	}
	return ""
}

func getDoc(doc *ast.CommentGroup) string {
	if doc == nil {
		return ""
	}
	ll := len(doc.List)
	if ll <= 0 {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(doc.List[ll-1].Text, "//"))
	//return doc.List[ll-1].Text
}

func getType(n ast.Expr, pkgName string, doc *ast.CommentGroup) (tt *domain.Type, token uint32) {
	tt = &domain.Type{Doc: getDoc(doc), Selector: pkgName}
	parseType(n, tt, &token)
	tt = GetOrAddType(tt)
	return
}

func parseType(n ast.Expr, tt *domain.Type, token *uint32) {
	switch vv := n.(type) {
	case *ast.Ident:
		tt.Name = vv.Name
	case *ast.SelectorExpr:
		tt.Selector = vv.X.(*ast.Ident).Name
		tt.Name = vv.Sel.Name
	case *ast.ArrayType:
		*token |= domain.ARRAY
		parseType(vv.Elt, tt, token)
	case *ast.StarExpr:
		*token |= domain.POINTER
		parseType(vv.X, tt, token)
	}
}
