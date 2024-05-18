package manager

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/types"
)

var (
	genMgr = make(map[string]domain.IMaker)
)

func Register(name string, act domain.IMaker) {
	if _, ok := genMgr[name]; ok {
		panic(fmt.Sprintf("%s has already registered", name))
	}
	genMgr[name] = act
}

func GetMaker(name string) domain.IMaker {
	return genMgr[name]
}

func Help() {
	fmt.Fprintf(flag.CommandLine.Output(), "action使用说明: \n")
	for name, item := range genMgr {
		fmt.Fprint(flag.CommandLine.Output(), item.GetHelp(name))
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
			AddConst(d.pkgName, types.ParseComment(n.Doc), n.Specs)
		} else if n.Tok == token.TYPE {
			AddType(d.pkgName, types.ParseComment(n.Doc), n.Specs)
		}
	}
	return nil
}
