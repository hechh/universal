package maker

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"universal/common/pb"
	"universal/framework/fbasic"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/base"
)

type GenFunc func(*domain.CmdLine, *base.Templates) error

type BaseMaker struct {
	help  string
	param string
	gen   GenFunc
	tpls  *base.Templates
}

func NewBaseMaker(f GenFunc, param, help string) *BaseMaker {
	return &BaseMaker{help: help, param: param, gen: f}
}

func (d *BaseMaker) GetHelp(name string) string {
	var str string
	if len(d.param) > 0 {
		str = fmt.Sprintf("    -action=%s -param=%s", name, d.param)
	} else {
		str = fmt.Sprintf("    -action=%s", name)
	}
	return fmt.Sprintf("%-70s #%s\n", str, d.help)
}

func (d *BaseMaker) OpenTpl(cmd *domain.CmdLine) error {
	d.tpls = base.NewTemplates(cmd.Tpl)
	return nil
}

func (d *BaseMaker) ParseFile(cmd *domain.CmdLine, extend interface{}) error {
	vistor, ok := extend.(ast.Visitor)
	if !ok || vistor == nil {
		return fbasic.NewUError(1, pb.ErrorCode_Parameter, "ast.Visitor")
	}
	if len(cmd.Src) <= 0 {
		return fbasic.NewUError(1, pb.ErrorCode_Parameter, "-src")
	}
	//解析文件
	fset := token.NewFileSet()
	for _, pp := range strings.Split(cmd.Src, ",") {
		pp = base.GetAbsPath(pp, base.GetCwd())
		if !strings.HasSuffix(pp, ".go") {
			pp = filepath.Join(pp, "*.go")
		}
		// 读取所有文件
		files, err := filepath.Glob(pp)
		if err != nil {
			return fbasic.NewUError(1, -1, err)
		}
		// 解析文件
		for _, file := range files {
			f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
			if err != nil {
				return fbasic.NewUError(2, -1, err)
			}
			ast.Walk(vistor, f)
		}
	}
	return nil
}

func (d *BaseMaker) Gen(cmd *domain.CmdLine) error {
	return d.gen(cmd, d.tpls)
}
