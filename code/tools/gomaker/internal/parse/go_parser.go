package parse

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"sort"
	"strings"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/typespec"

	"github.com/spf13/cast"
)

type GoParser struct {
	filename string
	pkg      string
	doc      *ast.CommentGroup
}

func (d *GoParser) SetFile(filename string) {
	d.filename = filepath.Base(filename)
}

func (d *GoParser) Close() {}

func (d *GoParser) Visit(node ast.Node) ast.Visitor {
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
			manager.AddEnum(d.GetEnum(n))
		}
		return nil
	case *ast.TypeSpec:
		switch n.Type.(type) {
		case *ast.StructType:
			manager.AddStruct(d.GetStruct(n))
		default:
			manager.AddAlias(d.GetAlias(n))
		}
	}
	return nil
}

func (d *GoParser) GetAlias(vv *ast.TypeSpec) *typespec.Alias {
	tt, token := getType(0, vv.Type, nil, d.pkg)
	return &typespec.Alias{
		Type: manager.GetOrAddType(&typespec.Type{
			Kind:    domain.KIND_ALIAS,
			PkgName: d.pkg,
			Name:    vv.Name.Name,
			Doc:     getDoc(d.doc),
		}),
		FileName: d.filename,
		Token:    token,
		RealType: tt,
		Doc:      getDoc(vv.Doc),
	}
}

func (d *GoParser) GetStruct(vv *ast.TypeSpec) *typespec.Struct {
	ret := typespec.NewStruct(d.filename, manager.GetOrAddType(&typespec.Type{
		Kind:    domain.KIND_STRUCT,
		PkgName: d.pkg,
		Name:    vv.Name.Name,
		Doc:     getDoc(d.doc),
	}))
	for i, field := range vv.Type.(*ast.StructType).Fields.List {
		tt, token := getType(0, field.Type, nil, d.pkg)
		ret.Add(&typespec.Field{
			Index: i,
			Token: token,
			Name:  field.Names[0].Name,
			Type:  tt,
			Tag:   getTag(field.Tag),
			Doc:   getDoc(field.Doc),
		})
	}
	sort.Slice(ret.List, func(i, j int) bool {
		return strings.Compare(ret.List[i].Name, ret.List[j].Name) < 0
	})
	return ret
}

func (d *GoParser) GetEnum(vv *ast.GenDecl) *typespec.Enum {
	ret := typespec.NewEnum(d.filename, nil)
	for _, spec := range vv.Specs {
		vv, ok := spec.(*ast.ValueSpec)
		if !ok || vv == nil || vv.Type == nil {
			continue
		}
		ret.Type, _ = getType(domain.KIND_ENUM, vv.Type, nil, d.pkg)
		ret.Add(&typespec.Value{
			Name:  vv.Names[0].Name,
			Type:  ret.Type,
			Value: cast.ToInt32(vv.Values[0].(*ast.BasicLit).Value),
			Doc:   getDoc(vv.Doc),
		})
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
	tt := &typespec.Type{Kind: kind, Doc: getDoc(doc), PkgName: pkg}
	parseType(n, tt, &token)
	return manager.GetOrAddType(tt), token
}

func parseType(n ast.Expr, tt *typespec.Type, token *[]uint32) {
	switch vv := n.(type) {
	case *ast.Ident:
		tt.Name = vv.Name
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
