package manager

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/typespec"
)

var (
	genMgr = make(map[string]domain.IParser)
)

func Register(act domain.IParser) {
	name := act.GetAction()
	if _, ok := genMgr[name]; ok {
		panic(fmt.Sprintf("%s has already registered", name))
	}
	genMgr[name] = act
}

func GetParser(name string) domain.IParser {
	return genMgr[name]
}

func Help() {
	fmt.Fprintf(flag.CommandLine.Output(), "action使用说明: \n")
	for _, item := range genMgr {
		fmt.Fprint(flag.CommandLine.Output(), item.Help())
	}
}

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
			AddConst(d.pkgName, typespec.ParseComment(n.Doc), n.Specs)
		} else if n.Tok == token.TYPE {
			AddType(d.pkgName, typespec.ParseComment(n.Doc), n.Specs)
		}
	}
	return nil
}
