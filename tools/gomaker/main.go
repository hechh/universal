package main

import (
	"flag"
	"fmt"
	"go/token"
	"os"
	"strings"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/generator"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/util"
)

func main() {
	var action, tpl, src, dst, param string
	flag.StringVar(&action, "action", "", "操作模式")
	flag.StringVar(&param, "param", "", "参数")
	flag.StringVar(&tpl, "tpl", "", "模板文件目录")
	flag.StringVar(&src, "src", "", "原文件目录")
	flag.StringVar(&dst, "dst", "", "生成文件目录")
	flag.Parse()

	// 获取生成器
	if !manager.IsAction(action) {
		panic(fmt.Sprintf("%s不支持", action))
	}

	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// 获取绝对地址
	tpl = util.GetAbsPath(cwd, tpl)
	src = util.GetAbsPath(cwd, src)
	dst = util.GetAbsPath(cwd, dst)

	// 加载模板文件
	tplFile, err := util.OpenTemplate(tpl)
	if err != nil {
		panic(err)
	}

	// 解析文件
	fset := token.NewFileSet()
	if strings.HasSuffix(src, ".go") {
		err = util.ParseFile(&manager.TypeParser{}, fset, src)
	} else {
		err = util.ParseDir(&manager.TypeParser{}, fset, src)
	}
	if err != nil {
		panic(err)
	}

	// 生成文件
	if err := manager.Generator(action, dst, param, tplFile); err != nil {
		panic(err)
	}
}

func init() {
	flag.Usage = func() {
		flag.PrintDefaults()
		manager.Help()
	}

	manager.Register(domain.HTTPKIT, "生成client代码", generator.HttpKitGenerator)
	manager.Register(domain.PBCLASS, "生成client代码", generator.OmitEmptyGenerator)
	manager.Register(domain.PROTO, "生成client代码", generator.ProtoGenerator)
}
