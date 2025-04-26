package parse

import (
	"go/ast"
	"go/parser"
	"go/token"
	"universal/library/baselib/uerror"
)

// 解析go文件
func ParseFiles(v ast.Visitor, files ...string) error {
	fset := token.NewFileSet()
	for _, filename := range files {
		// 解析语法树
		fs, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
		if err != nil {
			return uerror.New(1, -1, "filename: %v, error: %v", filename, err)
		}

		// 遍历语法树
		ast.Walk(v, fs)
	}
	return nil
}
