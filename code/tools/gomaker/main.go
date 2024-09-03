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
)

func main() {
	var action, tpl, src, dst, param, jspath string
	flag.StringVar(&action, "action", "", "操作模式")
	flag.StringVar(&param, "param", "", "参数")
	flag.StringVar(&tpl, "tpl", "", "模板文件目录")
	flag.StringVar(&jspath, "json", "./", "json文件生成目录")
	flag.StringVar(&src, "src", "", "原文件目录")
	flag.StringVar(&dst, "dst", "", "生成文件目录")
	flag.Parse()

	// 获取生成器
	if !manager.IsAction(action) {
		panic(fmt.Sprintf("%s不支持", action))
	}

	// 获取当前工作目录 + 获取绝对地址
	if cwd, err := os.Getwd(); err != nil {
		panic(err)
	} else {
		tpl = util.GetAbsPath(cwd, tpl)
		src = util.GetAbsPath(cwd, src)
		dst = util.GetAbsPath(cwd, dst)
	}

	// 加载模板文件
	tplFile, err := util.OpenTemplate(tpl)
	if err != nil {
		panic(err)
	}

	// 解析文件
	switch action {
	case domain.XLSX:
		if strings.HasSuffix(src, ".xlsx") {
			if err := util.ParseXlsx(parse.NewXlsxParser(), src); err != nil {
				panic(err)
			}
		} else {
			if err := util.ParseDirXlsx(parse.NewXlsxParser(), src); err != nil {
				panic(err)
			}
		}
		// 生成json文件
		if err := util.SaveJson(jspath, manager.GetJsons()); err != nil {
			panic(err)
		}
	default:
		fset := token.NewFileSet()
		if strings.HasSuffix(src, ".go") {
			if err := util.ParseFile(&parse.GoParser{}, fset, src); err != nil {
				panic(err)
			}
		} else {
			if err := util.ParseDir(&parse.GoParser{}, fset, src); err != nil {
				panic(err)
			}
		}
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

	manager.Register(domain.HTTPKIT, client.HttpKitGenerator, "生成client代码")
	manager.Register(domain.PBCLASS, client.OmitEmptyGenerator, "生成client代码")
	manager.Register(domain.PROTO, client.ProtoGenerator, "生成client代码")
}
