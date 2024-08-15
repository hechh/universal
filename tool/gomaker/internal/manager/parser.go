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
	pkg string
	doc *ast.CommentGroup
}

func (d *TypeParser) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.File:
		d.pkg = n.Name.Name
		return d
	case *ast.GenDecl:
		d.doc = n.Doc
		if n.Tok == token.TYPE {
			return d
		}
		if n.Tok == token.CONST {
			AddEnum(d.GetEnum(n))
		}
		return nil
	case *ast.TypeSpec:
		switch n.Type.(type) {
		case *ast.StructType:
			AddStruct(d.GetStruct(n))
		default:
			AddAlias(d.GetAlias(n))
		}
	}
	return nil
}

func (d *TypeParser) GetAlias(vv *ast.TypeSpec) *typespec.Alias {
	tt, token := getType(0, vv.Type, nil, d.pkg)
	return &typespec.Alias{
		Token: token,
		Type: GetOrAddType(&typespec.Type{
			Kind:     domain.ALIAS,
			Selector: d.pkg,
			Name:     vv.Name.Name,
			Doc:      getDoc(d.doc),
		}),
		RealType: tt,
		Comment:  getDoc(vv.Comment),
	}
}

func (d *TypeParser) GetStruct(vv *ast.TypeSpec) *typespec.Struct {
	ret := &typespec.Struct{
		Fields: make(map[string]*typespec.Field),
		Type: GetOrAddType(&typespec.Type{
			Kind:     domain.STRUCT,
			Selector: d.pkg,
			Name:     vv.Name.Name,
			Doc:      getDoc(d.doc),
		}),
	}
	for _, field := range vv.Type.(*ast.StructType).Fields.List {
		tt, token := getType(0, field.Type, nil, d.pkg)
		ff := &typespec.Field{
			Token:   token,
			Name:    field.Names[0].Name,
			Type:    tt,
			Tag:     getTag(field.Tag),
			Comment: getDoc(field.Comment),
		}
		ret.Fields[ff.Name] = ff
		ret.List = append(ret.List, ff)
	}
	sort.Slice(ret.List, func(i, j int) bool {
		return strings.Compare(ret.List[i].Name, ret.List[j].Name) < 0
	})
	return ret
}

func (d *TypeParser) GetEnum(vv *ast.GenDecl) *typespec.Enum {
	ret := &typespec.Enum{Values: make(map[string]*typespec.Value)}
	for _, spec := range vv.Specs {
		vv, ok := spec.(*ast.ValueSpec)
		if !ok || vv == nil || vv.Type == nil {
			continue
		}
		tt, _ := getType(domain.ENUM, vv.Type, nil, d.pkg)
		val := &typespec.Value{
			Name:    vv.Names[0].Name,
			Type:    tt,
			Value:   cast.ToInt32(vv.Values[0].(*ast.BasicLit).Value),
			Comment: getDoc(vv.Comment),
		}
		ret.List = append(ret.List, val)
		ret.Values[val.Name] = val
		ret.Type = tt
	}
	if len(ret.List) <= 0 {
		return nil
	}
	sort.Slice(ret.List, func(i, j int) bool {
		return ret.List[i].Value < ret.List[j].Value
	})
	return ret
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

func getType(kind uint32, n ast.Expr, doc *ast.CommentGroup, pkg string) (*typespec.Type, []uint32) {
	token := []uint32{}
	tt := &typespec.Type{Kind: kind, Doc: getDoc(doc), Selector: pkg}
	parseType(n, tt, &token)
	return GetOrAddType(tt), token
}

func parseType(n ast.Expr, tt *typespec.Type, token *[]uint32) {
	switch vv := n.(type) {
	case *ast.Ident:
		tt.Name = vv.Name
	case *ast.SelectorExpr:
		tt.Selector = vv.X.(*ast.Ident).Name
		tt.Name = vv.Sel.Name
	case *ast.ArrayType:
		*token = append(*token, domain.ARRAY)
		parseType(vv.Elt, tt, token)
	case *ast.StarExpr:
		*token = append(*token, domain.POINTER)
		parseType(vv.X, tt, token)
	}
}
