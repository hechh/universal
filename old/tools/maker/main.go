package main

import (
	"flag"
	"os"
	"path/filepath"
	"universal/framework/basic"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/base"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/parser"
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
		base.Panic(err)
	}
	GetArgs(cwd).Handle()
}

type Args struct {
	action string // 模式
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
	flag.StringVar(&ret.src, "src", "", "原文件目录")
	flag.StringVar(&ret.dst, "dst", "", "生成文件目录")
	flag.Parse()
	ret.xlsx = filepath.Clean(filepath.Join(cwd, ret.xlsx))
	ret.src = filepath.Clean(filepath.Join(cwd, ret.src))
	ret.dst = filepath.Clean(filepath.Join(cwd, ret.dst))
	return ret
}

func (d *Args) Handle() {
	switch d.action {
	case "proto", "config":
		d.handleXlsx()
	default:
		d.handleGo()
	}
}

func (d *Args) handleXlsx() {
	// 解析自定义proto文件
	/*
		protos, err := basic.Glob(d.src, ".*\\.proto", ".*.gen.proto", true)
		if err != nil {
			base.Panic(err)
		}
		ppar := &parser.PbParser{}
		for _, filename := range protos {
			if err := ppar.ParseFile(filename); err != nil {
				base.Panic(err)
			}
		}
	*/
	// 解析xlsx文件生成表
	files, err := basic.Glob(d.xlsx, ".*\\.xlsx", "", true)
	if err != nil {
		base.Panic(err)
	}
	par := &parser.XlsxParser{}
	if err := par.ParseFiles(files...); err != nil {
		base.Panic(err)
	}
	// 生成文件
	base.Panic(manager.Generator(d.action, d.dst))
}

func (d *Args) handleGo() {
	// 加载所有go文件
	files, err := basic.Glob(d.src, ".*\\.pb\\.go", "", true)
	if err != nil {
		base.Panic(err)
	}
	// 以目的目录设置pkg
	domain.DefaultPkg = filepath.Base(d.dst)
	base.Panic(base.ParseFiles(&parser.Parser{}, files...))
	// 生成go文件
	base.Panic(manager.Generator(d.action, d.dst, d.param))
}
