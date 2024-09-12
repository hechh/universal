package main

import (
	"flag"
	"fmt"
	"os"
	"universal/framework/basic"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/parse"
	"universal/tools/gomaker/internal/util"
	"universal/tools/gomaker/repository/client"
	"universal/tools/gomaker/repository/config"
)

func init() {
	flag.Usage = func() {
		flag.PrintDefaults()
		manager.Help()
	}

	manager.Register(domain.HTTPKIT, client.HttpKitGenerator, "生成client代码")
	manager.Register(domain.PBCLASS, client.OmitEmptyGenerator, "生成client代码")
	manager.Register(domain.PROTO, client.ProtoGenerator, "生成client代码")
	manager.Register(domain.XLSX, config.XlsxGenerator, "生成配置代码")
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	ParseArgs(cwd).Handle()
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

func (d *Args) Handle() {
	switch d.action {
	case domain.JSON:
		d.handleJson()
	case domain.CONFIG, domain.XLSX:
		d.handleCfg()
	default:
		d.handleGo()
	}
}

func (d *Args) handleCfg() {
	if !manager.IsAction(d.action) {
		panic(fmt.Errorf("%s不支持", d.action))
	}
	// 搜索所有配置
	files, err := basic.Glob(d.srcpath, "*.xlsx", "", true)
	if err != nil {
		panic(err)
	}
	// 解析所有配置
	pp := parse.NewCfgParser()
	if err := util.ParseXlsxs(pp, util.SortSheet, files...); err != nil {
		panic(err)
	}
	// 加载模板文件
	tplFile, err := util.OpenTemplate(d.tplpath, "*.tpl", true)
	if err != nil {
		panic(err)
	}
	// 生成文件
	if err := manager.Generator(d.action, d.dstpath, d.param, tplFile); err != nil {
		panic(err)
	}
}

func (d *Args) handleJson() {
	// 搜索所有配置
	files, err := basic.Glob(d.srcpath, "*.xlsx", "", true)
	if err != nil {
		panic(err)
	}
	// 解析所有配置
	pp := parse.NewJsonParser(d.dstpath)
	if err := util.ParseXlsxs(pp, util.SortSheet, util.SortXlsx(files)...); err != nil {
		panic(err)
	}
}

func (d *Args) handleGo() {
	if !manager.IsAction(d.action) {
		panic(fmt.Errorf("%s不支持", d.action))
	}
	// 加载模板文件
	tplFile, err := util.OpenTemplate(d.tplpath, "*.tpl", true)
	if err != nil {
		panic(err)
	}
	// 加载所有go文件
	files, err := basic.Glob(d.srcpath, "*.go", "", true)
	if err != nil {
		panic(err)
	}
	// 解析文件
	if err := util.ParseFiles(&parse.GoParser{}, files...); err != nil {
		panic(err)
	}
	// 生成文件
	if err := manager.Generator(d.action, d.dstpath, d.param, tplFile); err != nil {
		panic(err)
	}
}
