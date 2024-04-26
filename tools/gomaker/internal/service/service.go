package service

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"universal/framework/fbasic"
	"universal/tools/gomaker/internal/base"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/typespec"
)

type TypeParser struct {
	pkgName string // 当前解析文件的包名
}

func (d *TypeParser) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.File:
		d.pkgName = n.Name.Name
		return d
	case *ast.GenDecl:
		if n.Tok == token.CONST {
			manager.AddConst(d.pkgName, typespec.ParseComment(n.Doc), n.Specs)
		} else if n.Tok == token.TYPE {
			manager.AddType(d.pkgName, typespec.ParseComment(n.Doc), n.Specs)
		}
	}
	return nil
}

func ParseFiles(src string, CwdPath string) error {
	if len(src) <= 0 {
		return nil
	}
	//解析文件
	fset := token.NewFileSet()
	d := &TypeParser{}
	for _, pp := range strings.Split(src, ",") {
		pp = base.GetAbsPath(pp, CwdPath)
		if !strings.HasSuffix(pp, ".go") {
			pp = filepath.Join(pp, "*.go")
		}
		// 读取所有文件
		files, err := filepath.Glob(pp)
		if err != nil {
			return fbasic.NewUError(2, -1, err)
		}
		// 解析文件
		for _, file := range files {
			f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
			if err != nil {
				return fbasic.NewUError(2, -1, err)
			}
			ast.Walk(d, f)
		}
	}
	// 修复pkgname
	manager.Finished()
	return nil
}
