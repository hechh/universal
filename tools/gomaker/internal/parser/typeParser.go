package parser

import (
	"go/ast"
	"go/token"
	"universal/tools/gomaker/domain"
)

type TypeParser struct {
	pkgName string         // 当前解析文件的包名
	types   domain.IParser // 类型管理
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
