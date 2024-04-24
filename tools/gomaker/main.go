package main

import (
	"flag"
	"fmt"
	"os"
	"universal/tools/gomaker/domain"
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
	var src, dst, tpl, action, param string
	flag.StringVar(&action, "action", "", "操作模式")
	flag.StringVar(&param, "param", "", "设置参数，不同的action的params含义不同")
	flag.StringVar(&tpl, "tpl", "", "加载.tpl文件路径, 默认从${TPL_GO}环境变量中读取")
	flag.StringVar(&src, "src", "", "解析.go文件路径")
	flag.StringVar(&dst, "dst", "", "生成.gen.go文件路径")
	flag.Parse()

	var err error
	// 将相对路径转成绝对路径
	if len(action) <= 0 {
		fmt.Println("-action: parametar is empty")
		return
	}
	if tpl, err = base.GetAbsPath(base.GetPathDefault(tpl, "TPL_GO"), CwdPath); err != nil {
		fmt.Println("-tpl: ", err)
		return
	}
	if len(src) > 0 {
		src, _ = base.GetAbsPath(src, CwdPath)
	}
	if dst, err = base.GetAbsPath(dst, CwdPath); err != nil {
		fmt.Println("-dst: ", err)
		return
	}
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
		manager.HelpAction()
	}

	manager.Register(domain.UERRORS, uerrors.Gen, "生成errorCode错误码文件")
	manager.Register(domain.PLAYER_FUN, playerFun.Gen, "生成playerFun模板, 减少手写")
}
