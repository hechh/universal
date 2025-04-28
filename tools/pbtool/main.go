package main

import (
	"flag"
	"universal/library/baselib/util"
	"universal/tools/pbtool/internal/parse"
)

func main() {
	var src string
	flag.StringVar(&src, "src", "", "原文件目录")
	flag.Parse()

	// 加载所有go文件
	files, err := util.Glob(src, ".*\\.pb\\.go", true)
	if err != nil {
		panic(err)
	}
	// 以目的目录设置pkg
	parse.ParseFiles(&parse.Parser{}, files...)

}
