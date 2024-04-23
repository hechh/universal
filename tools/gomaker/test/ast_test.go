package test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"
	"universal/tools/gomaker/internal/parse"
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

func TestMM(t *testing.T) {
	par := parse.NewTypeParser()

	// 构建指定目录下所有文件的匹配模式
	files, err := filepath.Glob("./*.pb.go")
	if err != nil {
		panic(err)
	}
	if err = par.ParseFiles(files...); err != nil {
		panic(err)
	}
	//manager.GetTypeMgr().Print()
	// 生成文件
	//if err := uerrors.Gen(); err != nil {
	//panic(err)
	//}
	//tmgr.Print()
}
