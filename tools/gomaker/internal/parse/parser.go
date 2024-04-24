package parse

import (
	"go/ast"
	"go/parser"
	"go/token"
	"universal/framework/basic"
	"universal/tools/gomaker/internal/manager"
)

type TypeParser struct {
	pkgName string // 当前解析文件的包名
}

func NewTypeParser() *TypeParser {
	return &TypeParser{}
}

func (d *TypeParser) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.File:
		d.pkgName = n.Name.Name
		return d
	case *ast.GenDecl:
		if n.Tok == token.CONST {
			manager.AddConst(d.pkgName, n.Specs)
		} else if n.Tok == token.TYPE {
			manager.AddType(d.pkgName, n.Specs)
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
			return basic.NewUError(2, -1, err)
		}
		ast.Walk(d, f)
	}
	return nil
}
