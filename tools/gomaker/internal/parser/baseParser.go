package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"universal/common/pb"
	"universal/framework/fbasic"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/base"
)

type BaseParser struct {
	name  string
	help  string
	param string
	gen   domain.GenFunc
	tpls  map[string]*template.Template
}

func NewBaseParser(f domain.GenFunc, name, param, help string) *BaseParser {
	return &BaseParser{name: name, help: help, param: param, gen: f}
}

func (d *BaseParser) GetAction() string {
	return d.name
}

func (d *BaseParser) GetHelp() string {
	var str string
	if len(d.param) > 0 {
		str = fmt.Sprintf("    -action=%s -param=%s", d.name, d.param)
	} else {
		str = fmt.Sprintf("    -action=%s", d.name)
	}
	return fmt.Sprintf("%-80s #%s\n", str, d.help)
}

func (d *BaseParser) OpenTpl(cwd string, cmd *domain.CmdLine) error {
	if len(cmd.Tpl) <= 0 {
		return fbasic.NewUError(1, pb.ErrorCode_Parameter, "-tpl")
	}
	d.tpls = make(map[string]*template.Template)
	// 便利所有模版文件
	filepath.Walk(cmd.Tpl, func(path string, info os.FileInfo, err error) error {
		if path == cmd.Tpl {
			return nil
		}
		if info.IsDir() {
			d.tpls[filepath.Base(path)] = template.Must(template.ParseGlob(path + "/*.tpl")).Funcs(template.FuncMap{"html": func(s string) string { return s }})
			return nil
		}
		return nil
	})
	return nil
}

func (d *BaseParser) ParseFile(cwd string, cmd *domain.CmdLine, extend interface{}) error {
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
		pp = base.GetAbsPath(pp, cwd)
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

func (d *BaseParser) Gen(cwd string, cmd *domain.CmdLine) error {
	return d.gen(cwd, cmd, d.tpls)
}
