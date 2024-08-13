package gomaker

import (
	"flag"
	"os"
	"universal/tool/gomaker/internal/util"
)

func main() {
	var action, tpl, src, dst, param string
	flag.StringVar(&action, "action", "", "操作模式")
	flag.StringVar(&param, "param", "", "参数")
	flag.StringVar(&tpl, "tpl", "", "模板文件目录")
	flag.StringVar(&src, "src", "", "原文件目录")
	flag.StringVar(&dst, "dst", "", "生成文件目录")
	flag.Parse()

	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// 获取绝对地址
	tpl = util.GetAbsPath(cwd, tpl)
	src = util.GetAbsPath(cwd, src)
	dst = util.GetAbsPath(cwd, dst)

}
