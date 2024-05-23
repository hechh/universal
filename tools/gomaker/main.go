package main

import (
	"flag"
	"fmt"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/common/base"
	"universal/tools/gomaker/internal/common/types"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/repository/uerrors"
)

func main() {
	cmdLine := &domain.CmdLine{}
	if err := base.InitCmdLine(cmdLine); err != nil {
		panic(err)
	}
	// 获取解析器
	par := manager.GetMaker(cmdLine.Action)
	if par == nil {
		panic(fmt.Sprintf("-action=%s not suppoerted", cmdLine.Action))
	}
	// 加载模板文件
	if err := par.OpenTpl(cmdLine); err != nil {
		panic(err)
	}
	// 解析go文件
	if err := par.ParseFile(cmdLine, types.NewTypeParser(manager.AddType)); err != nil {
		fmt.Printf("parseFiles is faield, error: %v", err)
		return
	}
	manager.Finish()
	// 生成文件
	if err := par.Gen(cmdLine); err != nil {
		fmt.Printf("error: %v", err.Error())
	}
}

func init() {
	uerrors.Init()
	// 设置默认的help函数
	flag.Usage = func() {
		flag.PrintDefaults()
		manager.Help()
	}
}
