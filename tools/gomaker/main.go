package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"universal/framework/basic"
	"universal/tools/gomaker/generator"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/parse"
	"universal/tools/gomaker/internal/util"
)

func init() {
	flag.Usage = func() {
		flag.PrintDefaults()
		manager.Help()
	}
	generator.Init()
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	ParseArgs(cwd).Handle()
}

type Args struct {
	action   string // 模式
	tplpath  string // tpl模板文件目录
	xlsxpath string // xlsx文件目录
	srcpath  string // 解析原文件目录
	dstpath  string // 生成go文件目录
	param    string // action模式参数
}

func ParseArgs(cwd string) *Args {
	var action, tpl, src, dst, xlsx, param string
	flag.StringVar(&action, "action", "", "操作模式")
	flag.StringVar(&xlsx, "xlsx", "", "xlsx文件目录")
	flag.StringVar(&param, "param", "", "参数")
	flag.StringVar(&tpl, "tpl", "", "模板文件目录")
	flag.StringVar(&src, "src", "", "原文件目录")
	flag.StringVar(&dst, "dst", "", "生成文件目录")
	flag.Parse()
	return &Args{
		action:   action,
		param:    param,
		tplpath:  filepath.Clean(filepath.Join(cwd, tpl)),
		xlsxpath: filepath.Clean(filepath.Join(cwd, xlsx)),
		srcpath:  filepath.Clean(filepath.Join(cwd, src)),
		dstpath:  filepath.Clean(filepath.Join(cwd, dst)),
	}
}

func (d *Args) Handle() {
	switch d.action {
	case "pb":
		d.handlePb()
	case "bytes":
		d.handleBytes()
	default:
		d.handleGo()
	}
}

func (d *Args) handleBytes() {
	if !manager.IsAction(d.action) {
		panic(fmt.Errorf("%s不支持", d.action))
	}
	// 加载所有go文件
	files, err := basic.Glob(d.srcpath, ".*\\.go", "", true)
	if err != nil {
		panic(err)
	}
	// 解析go文件
	if err := util.ParseFiles(&parse.Parser{}, files...); err != nil {
		panic(err)
	}
	manager.InitEvals()
	// 解析配置
}

func (d *Args) handlePb() {
	if !manager.IsAction(d.action) {
		panic(fmt.Errorf("%s不支持", d.action))
	}
	// 解析枚举类型
	par := parse.PbParser{}
	if err := par.ParseEnum(filepath.Join(d.xlsxpath, "enum.xlsx")); err != nil {
		panic(err)
	}
	manager.InitEvals()
	// 解析配置结构
	files, err := basic.Glob(d.xlsxpath, ".*\\.xlsx", "enum.xlsx", true)
	if err != nil {
		panic(err)
	}
	for _, filename := range files {
		if err := par.ParseConfig(filename); err != nil {
			panic(err)
		}
	}
	// 生成文件
	if err := manager.Generator(d.action, d.dstpath, nil); err != nil {
		panic(err)
	}
}

func (d *Args) handleGo() {
	if !manager.IsAction(d.action) {
		panic(fmt.Errorf("%s不支持", d.action))
	}
	// 加载模板文件
	tplFile, err := util.OpenTemplate(d.tplpath, ".*\\.tpl", "", true)
	if err != nil {
		util.Panic(err)
	}
	// 加载所有go文件
	files, err := basic.Glob(d.srcpath, ".*\\.go", "", true)
	if err != nil {
		panic(err)
	}
	// 解析文件
	if err := util.ParseFiles(&parse.Parser{}, files...); err != nil {
		panic(err)
	}
	// 生成文件
	if err := manager.Generator(d.action, d.dstpath, tplFile, d.param); err != nil {
		panic(err)
	}
}
