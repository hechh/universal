package service

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"

	"forevernine.com/planet/server/tool/gomaker/internal/base"
	"forevernine.com/planet/server/tool/gomaker/internal/manager"
)

type Analysis struct{}

func (d *Analysis) Visit(v ast.Node) ast.Visitor {
	switch n := v.(type) {
	case *ast.Package:
		return d
	case *ast.File:
		return d
	case *ast.GenDecl:
		switch n.Tok {
		case token.CONST:
			manager.AddAstEnum(base.ParseAstEnum(n))
			return nil
		case token.TYPE:
			// 解析规则
			item, ok := n.Specs[0].(*ast.TypeSpec)
			if !ok || item == nil {
				return nil
			}
			manager.AddConfigAry(item.Name.Name)
			manager.AddAstStruct(base.ParseAstStruct(item))
			if n.Doc == nil {
				return nil
			}
			for _, cc := range n.Doc.List {
				manager.ParseRule(item.Name.Name, base.TrimSpace(cc.Text))
			}
		}
		return nil
	}
	return nil
}

func ParseFiles(files ...string) {
	fset := token.NewFileSet()
	analysis := &Analysis{}
	for _, file := range files {
		f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		// ast文件
		/*
			fp, _ := os.OpenFile("ast.ini", os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.FileMode(0644))
			ast.Fprint(fp, fset, f, nil)
		*/

		ast.Walk(analysis, f)
	}
	// 更新
	manager.Update()
}

func GenCode(path, action string) {
	rules := manager.GetRuleTypes()
	if len(action) > 0 {
		rules = manager.GetRuleType(action)
	}
	buf := bytes.NewBuffer(nil)
	for _, rule := range rules {
		manager.GenCode(rule, path, buf)
	}
}
