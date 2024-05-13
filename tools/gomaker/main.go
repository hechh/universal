package main

import (
	"flag"
	"fmt"
	"os"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/base"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/repository/uerrors"
)

var (
	CwdPath string
)

func main() {
	cmdLine := &domain.CmdLine{}
	if err := base.InitCmdLine(CwdPath, cmdLine); err != nil {
		panic(err)
	}
	if len(cmdLine.Tpl) <= 0 {
		fmt.Println("-tpl(TPL_GO): parameter is empty")
		return
	}

	// 获取解析器
	par := manager.GetParser(cmdLine.Action)
	if par == nil {
		panic(fmt.Sprintf("-action=%s not suppoerted", cmdLine.Action))
	}
	// 加载模板文件
	if err := par.OpenTpl(CwdPath, cmdLine); err != nil {
		panic(err)
	}
	// 解析go文件
	if err := par.ParseFile(CwdPath, cmdLine, &manager.TypeParser{}); err != nil {
		fmt.Printf("parseFiles is faield, error: %v", err)
		return
	}
	manager.Finished()
	// 生成文件
	if err := par.Gen(CwdPath, cmdLine); err != nil {
		fmt.Printf("error: %v", err.Error())
	}
}

func init() {
	var err error
	if CwdPath, err = os.Getwd(); err != nil {
		panic(err)
	}
	// 设置默认的help函数
	flag.Usage = func() {
		flag.PrintDefaults()
		manager.Help()
	}

	//playerFun.Init()
	uerrors.Init()
	//entity.Init()
}
