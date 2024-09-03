package test

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"testing"
	"text/template"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/parse"
	"universal/tools/gomaker/internal/util"
)

func TestAst(t *testing.T) {
	fset := token.NewFileSet()
	filename := "../../../common/pb/common.pb.go"
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	fp, _ := os.OpenFile("../../../common/pb/common.pb.ini", os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.FileMode(0644))
	ast.Fprint(fp, fset, f, nil)
}

func TestPlayer(t *testing.T) {
	fset := token.NewFileSet()
	filename := "../../../common/pb/playerStruct.pb.go"
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	fp, _ := os.OpenFile("../../../common/pb/playerStruct.pb.ini", os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.FileMode(0644))
	ast.Fprint(fp, fset, f, nil)
}

func TestParser(t *testing.T) {
	fset := token.NewFileSet()
	//filename := "./pb/common.pb.go"
	filename := "../../../common/pb/playerStruct.pb.go"
	t.Log(util.ParseFile(&parse.GoParser{}, fset, filename))
	manager.Print()
}

func TestTpl(t *testing.T) {
	a := template.Must(template.New("package.tpl").Parse("package {{.}}"))
	buf := bytes.NewBuffer(nil)
	a.ExecuteTemplate(buf, "package.tpl", "hch")
	t.Log(buf.String())
}
