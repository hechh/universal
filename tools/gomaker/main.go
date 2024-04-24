package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"universal/tools/gomaker/internal/base"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/parse"
	"universal/tools/gomaker/repository/uerrors"
)

var (
	CwdPath string
)

func main() {
	var src, dst, tpl, action string
	flag.StringVar(&tpl, "tpl", "", "加载.tpl文件路径, 默认从${TPL_GO}环境变量中读取")
	flag.StringVar(&src, "src", "", "解析.go文件路径")
	flag.StringVar(&dst, "dst", "", "生成.gen.go文件路径")
	flag.StringVar(&action, "action", "uerrors", "操作模式")
	flag.Parse()

	// 将相对路径转成绝对路径
	var err error
	if tpl, err = base.GetAbsPath(base.GetPathDefault(tpl, "TPL_GO"), CwdPath); err != nil {
		fmt.Println("-tpl: ", err)
		return
	}
	if src, err = base.GetAbsPath(base.GetPathDefault(src, "SRC_GO"), CwdPath); err != nil {
		fmt.Println("-src: ", err)
		return
	}
	if dst, err = base.GetAbsPath(dst, CwdPath); err != nil {
		fmt.Println("-dst: ", err)
		return
	}

	// 加载配置文件
	manager.InitTpl(tpl)
	// 解析go文件
	fmt.Println("tpl: ", tpl)
	par := parse.NewTypeParser()
	for _, pp := range strings.Split(src, ",") {
		if !filepath.IsAbs(pp) {
			pp = filepath.Join(CwdPath, pp)
		}
		if !strings.HasSuffix(pp, ".go") {
			pp = filepath.Join(pp, "*.go")
		}
		fmt.Println("src: ", pp)
		files, err := filepath.Glob(pp)
		if err != nil {
			fmt.Printf("error: %v", err)
			return
		}
		if err = par.ParseFiles(files...); err != nil {
			fmt.Printf("error: %v", err)
			return
		}
	}
	// 生成文件
	fmt.Println("dst: ", dst)
	if err := manager.Gen(action, dst); err != nil {
		fmt.Printf("error: %v", err)
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
	initRegister()
}

func initRegister() {
	manager.Register("uerrors", uerrors.Gen, "生成errorCode错误码文件")
}
