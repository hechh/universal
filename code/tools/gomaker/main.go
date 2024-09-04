package main

import (
	"flag"
	"fmt"
	"go/token"
	"os"
	"strings"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/parse"
	"universal/tools/gomaker/internal/util"
	"universal/tools/gomaker/repository/client"
	"universal/tools/gomaker/repository/xlsx"
)

func init() {
	flag.Usage = func() {
		flag.PrintDefaults()
		manager.Help()
	}

	manager.Register(domain.HTTPKIT, client.HttpKitGenerator, "生成client代码")
	manager.Register(domain.PBCLASS, client.OmitEmptyGenerator, "生成client代码")
	manager.Register(domain.PROTO, client.ProtoGenerator, "生成client代码")
	manager.Register(domain.JSON, xlsx.JsonGenerator, "xlsx配置文件生成json文件")
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	ParseArgs(cwd).Check().Handle()
}

type Args struct {
	action  string // 模式
	tplpath string // tpl模板文件目录
	srcpath string // 解析原文件目录
	dstpath string // 生成go文件目录
	param   string // action模式参数
}

func ParseArgs(cwd string) *Args {
	var action, tpl, src, dst, param string
	flag.StringVar(&action, "action", "", "操作模式")
	flag.StringVar(&param, "param", "", "参数")
	flag.StringVar(&tpl, "tpl", "", "模板文件目录")
	flag.StringVar(&src, "src", "", "原文件目录")
	flag.StringVar(&dst, "dst", "", "生成文件目录")
	flag.Parse()
	return &Args{
		action:  action,
		param:   param,
		tplpath: util.GetAbsPath(cwd, tpl),
		srcpath: util.GetAbsPath(cwd, src),
		dstpath: util.GetAbsPath(cwd, dst),
	}
}

func (d *Args) Check() *Args {
	if !manager.IsAction(d.action) {
		panic(fmt.Errorf("%s不支持", d.action))
	}
	return d
}

func (d *Args) Handle() {
	switch d.action {
	case domain.JSON:
		d.handleJson()
	default:
		d.handleGo()
	}
}

func (d *Args) handleJson() {
	// 解析表格
	if strings.HasSuffix(d.srcpath, ".xlsx") {
		if err := util.ParseXlsx(parse.NewXlsxParser(), d.srcpath); err != nil {
			panic(err)
		}
	} else {
		if err := util.ParseDirXlsx(parse.NewXlsxParser(), d.srcpath); err != nil {
			panic(err)
		}
	}

	// 生成json文件
	if err := manager.Generator(d.action, d.dstpath, d.param, nil); err != nil {
		panic(err)
	}
}

func (d *Args) handleGo() {
	// 加载模板文件
	tplFile, err := util.OpenTemplate(d.tplpath)
	if err != nil {
		panic(err)
	}

	// 解析文件
	fset := token.NewFileSet()
	if strings.HasSuffix(d.srcpath, ".go") {
		if err := util.ParseFile(&parse.GoParser{}, fset, d.srcpath); err != nil {
			panic(err)
		}
	} else {
		if err := util.ParseDir(&parse.GoParser{}, fset, d.srcpath); err != nil {
			panic(err)
		}
	}

	// 生成文件
	if err := manager.Generator(d.action, d.dstpath, d.param, tplFile); err != nil {
		panic(err)
	}
}
