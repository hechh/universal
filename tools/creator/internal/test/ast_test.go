package test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	//analysis := &Analysis{}
	fset := token.NewFileSet()
	for _, file := range []string{"error.pb.go", "packet.pb.go"} {
		f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		// ast文件
		fp, _ := os.OpenFile(file+".ini", os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.FileMode(0644))
		ast.Fprint(fp, fset, f, nil)

		//ast.Walk(analysis, f)
	}
}
