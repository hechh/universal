package visitor

import (
	"go/ast"
	"go/parser"
	"go/token"
	"universal/tools/gomaker/domain"
)

type TypeParser struct {
	pkgName string         // 当前解析文件的包名
	types   domain.IParser // 类型管理
}

func NewTypeParser(t domain.IParser) *TypeParser {
	return &TypeParser{types: t}
}

func (d *TypeParser) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.File:
		d.pkgName = n.Name.Name
		return d
	case *ast.GenDecl:
		if n.Tok == token.CONST {
			d.types.AddConst(d.pkgName, n.Specs)
		} else if n.Tok == token.TYPE {
			d.types.AddType(d.pkgName, n.Specs)
		}
	}
	return nil
}

func (d *TypeParser) ParseFiles(files ...string) error {
	//解析文件
	fset := token.NewFileSet()
	for _, file := range files {
		f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		ast.Walk(d, f)
	}
	return nil
}
