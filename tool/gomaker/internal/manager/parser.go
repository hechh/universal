package manager

import (
	"fmt"
	"go/ast"
	"go/token"
	"universal/tool/gomaker/domain"

	"github.com/spf13/cast"
)

var (
	types   = make(map[string]*domain.Type)
	values  = make(map[string]map[int32]*domain.Value)
	structs = make(map[string]*domain.Struct)
	alias   = make(map[string]*domain.Alias)
)

func AddValue(vv *domain.Value) {
	// 存储类型
	key := fmt.Sprintf("%s.%s", vv.Type.Selector, vv.Type.Name)
	if _, ok := types[key]; !ok {
		types[key] = vv.Type
	}
	// 存储数据
	if _, ok := values[key]; !ok {
		values[key] = make(map[int32]*domain.Value)
	}
	values[key][vv.Value] = vv
}

func AddStruct(vv *domain.Struct) {
	// 存储类型
	key := fmt.Sprintf("%s.%s", vv.Type.Selector, vv.Type.Name)
	if _, ok := types[key]; !ok {
		types[key] = vv.Type
	}
	// field类型存储
	for _, field := range vv.List {
		fkey := fmt.Sprintf("%s.%s", field.Type.Selector, field.Type.Name)
		if _, ok := types[fkey]; !ok {
			types[fkey] = field.Type
		}
	}
	// 存储struct数据
	structs[key] = vv
}

func AddAlias(vv *domain.Alias) {
	// 存储类型
	key := fmt.Sprintf("%s.%s", vv.AliasType.Selector, vv.AliasType.Name)
	if _, ok := types[key]; !ok {
		types[key] = vv.AliasType
	}
	rkey := fmt.Sprintf("%s.%s", vv.Type.Selector, vv.Type.Name)
	if _, ok := types[rkey]; !ok {
		types[rkey] = vv.Type
	}
	// 存储别名
	alias[key] = vv
}

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
	case *ast.ValueSpec:
		AddValue(d.getValue(n))
	case *ast.TypeSpec:
		switch vv := n.Type.(type) {
		case *ast.StructType:
			AddStruct(d.getStruct(n, vv))
		default:
		}
	}
	return nil
}

func (d *TypeParser) getAlias(vv *ast.TypeSpec) *domain.Alias {
	tt, token := getType(vv.Type, nil, d.pkgName)
	return &domain.Alias{
		Token: token,
		AliasType: &domain.Type{
			Kind:     domain.ALIAS,
			Selector: d.pkgName,
			Name:     vv.Name.Name,
			Doc:      getDoc(d.doc),
		},
		Type:    tt,
		Comment: getDoc(vv.Comment),
	}
}

func (d *TypeParser) getStruct(vv *ast.TypeSpec, st *ast.StructType) *domain.Struct {
	tmps := make(map[string]*domain.Field)
	list := []*domain.Field{}
	for _, field := range st.Fields.List {
		tt, token := getType(vv.Type, d.doc, d.pkgName)
		ff := &domain.Field{
			Token:   token,
			Name:    field.Names[0].Name,
			Type:    tt,
			Tag:     field.Tag.Value,
			Comment: getDoc(vv.Comment),
		}
		tmps[ff.Name] = ff
		list = append(list, ff)
	}
	return &domain.Struct{
		Type: &domain.Type{
			Kind:     domain.STRUCT,
			Selector: d.pkgName,
			Name:     vv.Name.Name,
			Doc:      getDoc(d.doc),
		},
		Fields: make(map[string]*domain.Field),
		List:   list,
	}
}

func (d *TypeParser) getValue(vv *ast.ValueSpec) *domain.Value {
	tt, _ := getType(vv.Type, d.doc, d.pkgName)
	return &domain.Value{
		Name:    vv.Names[0].Name,
		Type:    tt,
		Value:   cast.ToInt32(vv.Values[0].(*ast.BasicLit).Value),
		Comment: getDoc(vv.Comment),
	}
}

func getDoc(doc *ast.CommentGroup) string {
	if doc == nil || len(doc.List) <= 0 {
		return ""
	}
	return doc.List[0].Text
}

func getType(n ast.Expr, doc *ast.CommentGroup, pkgName string) (tt *domain.Type, token uint32) {
	tt = &domain.Type{Doc: getDoc(doc), Selector: pkgName}
	parseType(n, tt, &token)
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
