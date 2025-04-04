package parser

import (
	"go/ast"
	"go/token"
	"hego/tools/gomaker/domain"
	"hego/tools/gomaker/internal/base"
	"hego/tools/gomaker/internal/manager"
	"hego/tools/gomaker/internal/typespec"
	"strings"

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
		if n.Tok == token.CONST && len(n.Specs) > 0 {
			base.Panic(manager.AddEnum(d.parseEnum(n)))
		}
		return nil
	case *ast.TypeSpec:
		switch n.Type.(type) {
		case *ast.StructType:
			base.Panic(manager.AddStruct(d.parseStruct(n)))
		default:
			base.Panic(manager.AddAlias(d.parseAlias(n)))
		}
	}
	return nil
}

func (d *Parser) parseStruct(n *ast.TypeSpec) *typespec.Struct {
	list := []*typespec.Field{}
	for i, field := range n.Type.(*ast.StructType).Fields.List {
		tt, tos := d.parseType(domain.KindTypeIdent, domain.SourceTypeGo, d.pkg, field.Type)
		list = append(list, typespec.FIELD(tt, field.Names[0].Name, i+1, d.parseTag(field.Tag), d.parseDoc(field.Comment), tos...))
	}
	ttt := manager.GetType(domain.KindTypeStruct, domain.SourceTypeGo, d.pkg, n.Name.Name, "")
	return typespec.STRUCT(ttt, d.parseDoc(n.Doc), list...)
}

func (d *Parser) parseAlias(n *ast.TypeSpec) *typespec.Alias {
	tt, tos := d.parseType(domain.KindTypeIdent, domain.SourceTypeGo, d.pkg, n.Type)
	ttt := manager.GetType(domain.KindTypeAlias, domain.SourceTypeGo, d.pkg, n.Name.Name, "")
	return typespec.ALIAS(ttt, tt, d.parseDoc(n.Doc), tos...)
}

func (d *Parser) parseEnum(n *ast.GenDecl) *typespec.Enum {
	values := []*typespec.Value{}
	for _, spec := range n.Specs {
		vv, ok := spec.(*ast.ValueSpec)
		// 过滤非枚举的常量定义
		if !ok || vv == nil || vv.Type == nil {
			continue
		}
		// 解析
		tt, _ := d.parseType(domain.KindTypeEnum, domain.SourceTypeGo, d.pkg, vv.Type)
		vall := cast.ToInt32(vv.Values[0].(*ast.BasicLit).Value)
		values = append(values, typespec.VALUE(tt, vv.Names[0].Name, vall, d.parseDoc(vv.Comment)))
	}
	if len(values) > 0 {
		return typespec.ENUM(values[0].Type, d.parseDoc(n.Doc), values...)
	}
	return nil
}

func (d *Parser) parseTag(tt *ast.BasicLit) string {
	if tt != nil {
		return tt.Value
	}
	return ""
}

func (d *Parser) parseDoc(doc *ast.CommentGroup) string {
	if doc == nil {
		return ""
	}
	ll := len(doc.List)
	if ll <= 0 {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(doc.List[ll-1].Text, "//"))
}

func (d *Parser) parseType(kind, source int32, pkg string, n ast.Expr) (*typespec.Type, []rune) {
	tt := typespec.TYPE(kind, source, pkg, "", "")
	token := []rune{}
	parseAstType(n, tt, &token)
	return manager.GetTypeReference(tt), token
}

func parseAstType(n ast.Expr, tt *typespec.Type, token *[]int32) {
	switch vv := n.(type) {
	case *ast.Ident:
		tt.Name = vv.Name
		if _, ok := domain.BasicTypes[vv.Name]; ok {
			tt.Pkg = ""
			tt.Kind = domain.KindTypeIdent
		}
	case *ast.SelectorExpr:
		tt.Pkg = vv.X.(*ast.Ident).Name
		tt.Name = vv.Sel.Name
	case *ast.ArrayType:
		if token != nil {
			*token = append(*token, domain.TokenTypeArray)
		}
		parseAstType(vv.Elt, tt, token)
	case *ast.StarExpr:
		if token != nil {
			*token = append(*token, domain.TokenTypePointer)
		}
		parseAstType(vv.X, tt, token)
	}
}
