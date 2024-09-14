package parse

import (
	"go/ast"
	"go/token"
	"sort"
	"strings"
	"universal/tools/gomaker_new/domain"
	"universal/tools/gomaker_new/internal/manager"
	"universal/tools/gomaker_new/internal/typespec"
	"universal/tools/gomaker_new/internal/util"

	"github.com/spf13/cast"
)

type Parser struct {
	pkg string
	doc *ast.CommentGroup
}

func (d *Parser) Visit(node ast.Node) ast.Visitor {
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
			d.parseEnum(n)
		}
		return nil
	case *ast.TypeSpec:
		switch n.Type.(type) {
		case *ast.StructType:
			d.parseStruct(n)
		default:
			d.parseAlias(n)
		}
	}
	return nil
}

func (d *Parser) parseStruct(n *ast.TypeSpec) {
	ret := &typespec.Struct{
		Type: manager.GetTypeReference(&typespec.Type{
			Kind:    domain.KIND_STRUCT,
			PkgName: d.pkg,
			Name:    n.Name.Name,
			Doc:     parseDoc(d.doc),
		}),
		Fields: make(map[string]*typespec.Field),
	}
	for _, field := range n.Type.(*ast.StructType).Fields.List {
		tt, token := getType(domain.KIND_IDENT, field.Type, nil, d.pkg)
		ret.Add(&typespec.Field{
			Token: token,
			Type:  tt,
			Name:  field.Names[0].Name,
			Tag:   parseTag(field.Tag),
			Doc:   parseDoc(field.Doc),
		})
	}
	util.Panic(manager.AddStruct(ret))
}

func (d *Parser) parseAlias(n *ast.TypeSpec) {
	tt := manager.GetTypeReference(&typespec.Type{
		Kind:    domain.KIND_ALIAS,
		PkgName: d.pkg,
		Name:    n.Name.Name,
	})
	real, tokens := getType(domain.KIND_ALIAS, n.Type, nil, d.pkg)
	item := &typespec.Alias{Token: tokens, Type: tt, RealType: real, Doc: parseDoc(n.Doc)}
	util.Panic(manager.AddAlias(item))
}

func (d *Parser) parseEnum(n *ast.GenDecl) {
	ret := &typespec.Enum{
		Values: make(map[string]*typespec.Value),
		Doc:    parseDoc(d.doc),
	}
	for _, spec := range n.Specs {
		vv, ok := spec.(*ast.ValueSpec)
		if !ok || vv == nil || vv.Type == nil {
			continue
		}
		value := vv.Values[0].(*ast.BasicLit).Value
		tt, _ := getType(domain.KIND_ENUM, vv.Type, nil, d.pkg)
		ret.Add(tt, vv.Names[0].Name, cast.ToInt32(value), parseDoc(vv.Doc))
	}
	if len(ret.List) > 0 {
		sort.Slice(ret.List, func(i, j int) bool { return ret.List[i].Value < ret.List[j].Value })
		util.Panic(manager.AddEnum(ret))
	}
}

func parseTag(tt *ast.BasicLit) string {
	if tt != nil {
		return tt.Value
	}
	return ""
}

func getType(kind uint32, n ast.Expr, doc *ast.CommentGroup, pkg string) (*typespec.Type, []uint32) {
	token := []uint32{}
	tt := &typespec.Type{Kind: kind, Doc: parseDoc(doc), PkgName: pkg}
	parseType(n, tt, &token)
	return manager.GetTypeReference(tt), token
}

func parseDoc(doc *ast.CommentGroup) string {
	if doc == nil {
		return ""
	}
	ll := len(doc.List)
	if ll <= 0 {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(doc.List[ll-1].Text, "//"))
}

func parseType(n ast.Expr, tt *typespec.Type, token *[]uint32) {
	switch vv := n.(type) {
	case *ast.Ident:
		tt.Name = vv.Name
		switch vv.Name {
		case "uint32", "uint64", "int64", "int32", "bool", "float32", "float64", "string", "byte":
			tt.PkgName = ""
		}
	case *ast.SelectorExpr:
		tt.PkgName = vv.X.(*ast.Ident).Name
		tt.Name = vv.Sel.Name
	case *ast.ArrayType:
		*token = append(*token, domain.TOKEN_ARRAY)
		parseType(vv.Elt, tt, token)
	case *ast.StarExpr:
		*token = append(*token, domain.TOKEN_POINTER)
		parseType(vv.X, tt, token)
	}
}
