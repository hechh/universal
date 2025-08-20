package parse

import (
	"go/ast"
	"go/token"
	"universal/library/templ"
	"universal/tool/pbtool/domain"
	"universal/tool/pbtool/internal/typespec"
)

type GoParser struct {
	domain.IFactory
	pkg    string
	filter map[string]struct{}
}

func NewGoParser(fac domain.IFactory, fs ...string) *GoParser {
	tmps := map[string]struct{}{}
	for _, ff := range fs {
		tmps[ff] = struct{}{}
	}
	return &GoParser{IFactory: fac, filter: tmps}
}

func (g *GoParser) Visit(n ast.Node) ast.Visitor {
	switch vv := n.(type) {
	case *ast.File:
		g.pkg = vv.Name.Name
		return g
	case *ast.GenDecl:
		return templ.Or[*GoParser](vv.Tok != token.TYPE, nil, g)
	case *ast.TypeSpec:
		cls := typespec.NewClass(typespec.NewIdentExpr(g, g.pkg, vv.Name))
		switch nn := vv.Type.(type) {
		case *ast.StructType:
			for _, field := range nn.Fields.List {
				ft := typespec.ParseType(g, g.pkg, field.Type)
				for _, item := range field.Names {
					cls.Add(typespec.NewAttribute(item.Name, ft))
				}
			}
		}
	}
	return nil
}
