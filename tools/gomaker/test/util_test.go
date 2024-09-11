package test

/*
import (
	"bytes"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"testing"
	"text/template"
	"unicode"
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
	//filename := "./pb/common.pb.go"
	filename := "../../../common/pb/playerStruct.pb.go"
	t.Log(util.ParseFiles(&parse.GoParser{}, filename))
}

func TestTpl(t *testing.T) {
	a := template.Must(template.New("package.tpl").Parse("package {{.}}"))
	buf := bytes.NewBuffer(nil)
	a.ExecuteTemplate(buf, "package.tpl", "hch")
	t.Log(buf.String())
}

func TestJson(t *testing.T) {
	aa := []map[string]interface{}{
		{
			"tet":   123,
			"print": 1234,
		},
	}

	buf, _ := json.Marshal(aa)
	t.Log(string(buf))
}

func TestToUpper(t *testing.T) {
	aa := "addItem"
	buf := []byte(aa)
	buf[0] = byte(unicode.ToUpper(rune(buf[0])))
	t.Log(string(buf))
}
*/
