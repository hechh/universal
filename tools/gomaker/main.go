package main

import (
	"flag"
	"fmt"
	"os"
	"universal/tools/gomaker/internal/base"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/service"
	"universal/tools/gomaker/repository/playerFun"
	"universal/tools/gomaker/repository/uerrors"
)

var (
	CwdPath string
)

func main() {
	var action, param, tpl, src, dst string
	flag.StringVar(&action, "action", "", "操作模式")
	flag.StringVar(&param, "param", "", "不同action的param使用方式不同")
	flag.StringVar(&tpl, "tpl", "", "模板文件路径, 默认从TPL_GO环境变量中读取, -tpl={tpl文件目录}")
	flag.StringVar(&src, "src", "", "解析go文件或目录, -param={go文件目录 或 xxx/*.go}")
	flag.StringVar(&dst, "dst", "", "生成文件或目录, -dst={生成文件目录}")
	flag.Parse()

	// 将相对路径转成绝对路径
	if len(action) <= 0 {
		fmt.Println("-action: parameter is empty")
		return
	}
	if tpl = base.GetAbsPath(base.GetPathDefault(tpl, os.Getenv("TPL_GO")), CwdPath); len(tpl) <= 0 {
		fmt.Println("-tpl(TPL_GO): parameter is empty")
		return
	}
	dst = base.GetAbsPath(base.GetPathDefault(dst, CwdPath), CwdPath)
	// 加载模板文件
	manager.InitTpl(tpl)
	// 解析go文件
	if err := service.ParseFiles(src, CwdPath); err != nil {
		fmt.Printf("parseFiles is faield, error: %v", err)
		return
	}
	// 生成文件
	if err := manager.Gen(action, dst, param); err != nil {
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

	playerFun.Init()
	uerrors.Init()
}
