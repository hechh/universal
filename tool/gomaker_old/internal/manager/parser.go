package manager

import (
	"go/ast"
	"go/token"
	"sort"
	"strings"
	"universal/tool/gomaker/domain"
	"universal/tool/gomaker/internal/typespec"

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
		if n.Tok == token.TYPE {
			return d
		}
		if n.Tok == token.CONST {
			AddValue(getValue(n, d.pkgName, d.doc))
		}
		return nil
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

func getAlias(vv *ast.TypeSpec, pkgName string, doc *ast.CommentGroup) *typespec.Alias {
	tt, token := getType(0, vv.Type, pkgName, nil)
	return &typespec.Alias{
		Token: token,
		AliasType: GetOrAddType(&typespec.Type{
			Kind:     domain.ALIAS,
			Selector: pkgName,
			Name:     vv.Name.Name,
			Doc:      getDoc(doc),
		}),
		Type:    tt,
		Comment: getDoc(vv.Comment),
	}
}

func getStruct(vv *ast.TypeSpec, pkgName string, doc *ast.CommentGroup) *typespec.Struct {
	st := vv.Type.(*ast.StructType)
	tmps := make(map[string]*typespec.Field)
	list := []*typespec.Field{}
	for _, field := range st.Fields.List {
		tt, token := getType(0, vv.Type, pkgName, doc)
		ff := &typespec.Field{
			Token:   token,
			Name:    field.Names[0].Name,
			Type:    tt,
			Tag:     getTag(field.Tag),
			Comment: getDoc(vv.Comment),
		}
		tmps[ff.Name] = ff
		list = append(list, ff)
	}
	sort.Slice(list, func(i, j int) bool {
		return strings.Compare(list[i].Name, list[j].Name) < 0
	})
	return &typespec.Struct{
		Type: GetOrAddType(&typespec.Type{
			Kind:     domain.STRUCT,
			Selector: pkgName,
			Name:     vv.Name.Name,
			Doc:      getDoc(doc),
		}),
		Fields: make(map[string]*typespec.Field),
		List:   list,
	}
}

func getEnum(vv *ast.GenDecl, pkgName string, doc *ast.CommentGroup) *typespec.Enum {
	tmps := make(map[string]*typespec.Value)
	list := []*typespec.Value{}
	var tt *typespec.Type
	for _, spec := range vv.Specs {
		vv, ok := spec.(*ast.ValueSpec)
		if !ok || vv.Type == nil {
			continue
		}
		tt, _ = getType(domain.ENUM, vv.Type, pkgName, doc)
		// 保存
		item := &typespec.Value{
			Name:    vv.Names[0].Name,
			Type:    tt,
			Value:   cast.ToInt32(vv.Values[0].(*ast.BasicLit).Value),
			Comment: getDoc(vv.Comment),
		}
		tmps[vv.Names[0].Name] = item
		list = append(list, item)
	}
	return &typespec.Enum{Type: tt, Values: tmps, List: list}
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
}

func getType(kind uint32, n ast.Expr, pkgName string, doc *ast.CommentGroup) (tt *typespec.Type, token uint32) {
	tt = &typespec.Type{Kind: kind, Doc: getDoc(doc), Selector: pkgName}
	parseType(n, tt, &token)
	tt = GetOrAddType(tt)
	return
}

func parseType(n ast.Expr, tt *typespec.Type, token *uint32) {
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
