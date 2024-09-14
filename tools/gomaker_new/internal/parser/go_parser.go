package parser

import (
	"go/ast"
	"go/token"
	"strings"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/typespec"
	"universal/tools/gomaker/internal/util"

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
			util.Panic(manager.AddEnum(d.parseEnum(n)))
		}
		return nil
	case *ast.TypeSpec:
		switch n.Type.(type) {
		case *ast.StructType:
			util.Panic(manager.AddStruct(d.parseStruct(n)))
		default:
			util.Panic(manager.AddAlias(d.parseAlias(n)))
		}
	}
	return nil
}

func (d *Parser) parseStruct(n *ast.TypeSpec) *typespec.Struct {
	list := []*typespec.Field{}
	for _, field := range n.Type.(*ast.StructType).Fields.List {
		list = append(list, &typespec.Field{
			Type: d.parseType(0, d.pkg, field.Type),
			Name: field.Names[0].Name,
			Tag:  d.parseTag(field.Tag),
			Doc:  d.parseDoc(field.Comment),
		})
	}
	return typespec.STRUCT(manager.GetType(domain.KindTypeStruct, d.pkg, n.Name.Name), "", d.parseDoc(n.Doc), list...)
}

func (d *Parser) parseAlias(n *ast.TypeSpec) *typespec.Alias {
	return typespec.ALIAS(
		manager.GetType(domain.KindTypeAlias, d.pkg, n.Name.Name),
		d.parseType(0, d.pkg, n.Type),
		"",
		d.parseDoc(n.Doc),
	)
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
		values = append(values, typespec.VALUE(
			d.parseType(domain.KindTypeEnum, d.pkg, vv.Type),
			vv.Names[0].Name,
			cast.ToInt32(vv.Values[0].(*ast.BasicLit).Value),
			d.parseDoc(vv.Comment),
		))
	}
	return typespec.ENUM(values[0].Type, "", d.parseDoc(n.Doc), values...)
}

func (d *Parser) parseType(k int32, pkg string, n ast.Expr) *typespec.Type {
	tt := typespec.TYPE(k, pkg, "")
	parseAstType(n, tt)
	return manager.GetTypeReference(tt)
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

func parseAstType(n ast.Expr, tt *typespec.Type) {
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
		tt.Token = append(tt.Token, domain.TokenTypeArray)
		parseAstType(vv.Elt, tt)
	case *ast.StarExpr:
		tt.Token = append(tt.Token, domain.TokenTypePointer)
		parseAstType(vv.X, tt)
	}
}
