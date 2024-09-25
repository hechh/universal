package main

import (
	"flag"
	"os"
	"path/filepath"
	"universal/framework/basic"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/parser"
	"universal/tools/gomaker/internal/util"
	"universal/tools/gomaker/repository"
)

func init() {
	flag.Usage = func() {
		flag.PrintDefaults()
		manager.Help()
	}
	repository.Init()
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		util.Panic(err)
	}
	GetArgs(cwd).Handle()
}

type Args struct {
	action string // 模式
	tpl    string // tpl模板文件目录
	xlsx   string // xlsx文件目录
	src    string // 解析原文件目录
	dst    string // 生成go文件目录
	param  string // action模式参数
}

func GetArgs(cwd string) *Args {
	ret := new(Args)
	flag.StringVar(&ret.action, "action", "", "操作模式")
	flag.StringVar(&ret.xlsx, "xlsx", "", "xlsx文件目录")
	flag.StringVar(&ret.param, "param", "", "参数")
	flag.StringVar(&ret.tpl, "tpl", "", "模板文件目录")
	flag.StringVar(&ret.src, "src", "", "原文件目录")
	flag.StringVar(&ret.dst, "dst", "", "生成文件目录")
	flag.Parse()
	ret.tpl = filepath.Clean(filepath.Join(cwd, ret.tpl))
	ret.xlsx = filepath.Clean(filepath.Join(cwd, ret.xlsx))
	ret.src = filepath.Clean(filepath.Join(cwd, ret.src))
	ret.dst = filepath.Clean(filepath.Join(cwd, ret.dst))
	return ret
}

func (d *Args) Handle() {
	switch d.action {
	case "proto", "config":
		d.handleProto()
	case "bytes":
		d.handleBytes()
	default:
		d.handleGo()
	}
}

func (d *Args) handleBytes() {
	// 加载所有go文件
	files, err := basic.Glob(d.src, ".*\\.go", "", true)
	if err != nil {
		util.Panic(err)
	}
	util.Panic(util.ParseFiles(&parser.Parser{}, files...))
	// 初始化枚举数据
	manager.InitEvals()
	// 生成文件
	xlsxs, err := basic.Glob(d.xlsx, ".*\\.xlsx", "enum.xlsx", true)
	if err != nil {
		util.Panic(err)
	}
	util.Panic(manager.Generator(d.action, d.dst, nil, xlsxs...))
}

func (d *Args) handleProto() {
	files, err := basic.Glob(d.xlsx, ".*\\.xlsx", "", true)
	if err != nil {
		util.Panic(err)
	}
	// 解析xlsx文件生成表
	par := &parser.XlsxParser{}
	if err := par.ParseFiles(files...); err != nil {
		util.Panic(err)
	}
	// 生成文件
	util.Panic(manager.Generator(d.action, d.dst, nil))
}

func (d *Args) handleGo() {
	// 加载模板文件
	tpls, err := util.OpenTemplate(d.tpl, ".*\\.tpl", "", false)
	if err != nil {
		util.Panic(err)
	}
	// 加载所有go文件
	files, err := basic.Glob(d.src, ".*\\.go", "", true)
	if err != nil {
		util.Panic(err)
	}
	// 以目的目录设置pkg
	domain.DefaultPkg = filepath.Base(d.dst)
	util.Panic(util.ParseFiles(&parser.Parser{}, files...))
	// 生成go文件
	util.Panic(manager.Generator(d.action, d.dst, tpls, d.param))
}
