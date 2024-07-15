package types

import (
	"go/ast"
	"go/token"
)

type AddFunc func(interface{})

type TypeParser struct {
	pkgName string  // 当前解析文件的包名
	addFun  AddFunc // 增加函数
}

func NewTypeParser(f AddFunc) *TypeParser {
	return &TypeParser{addFun: f}
}

func (d *TypeParser) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.File:
		d.pkgName = n.Name.Name
		return d
	case *ast.GenDecl:
		doc := ParseComment(n.Doc)
		if n.Tok == token.CONST {
			d.addFun(ParseEnum(d.pkgName, doc, n.Specs...))
			return nil
		}
		if n.Tok != token.TYPE {
			return nil
		}
		for _, spec := range n.Specs {
			// 判断是否有效
			node, ok := spec.(*ast.TypeSpec)
			if !ok || node == nil {
				continue
			}
			// 解析结构
			switch vv := node.Type.(type) {
			case *ast.StructType:
				d.addFun(ParseStruct(d.pkgName, node.Name.Name, doc, vv.Fields.List))
			default:
				d.addFun(ParseAlias(d.pkgName, node.Name.Name, doc, node))
			}
		}
		return nil
	}
	return nil
}
