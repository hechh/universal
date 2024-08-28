package main

import (
	"flag"
	"path/filepath"
	"universal/tools/cfgtool/internal/manager"
	"universal/tools/cfgtool/internal/util"
)

func main() {
	var src, dst string
	flag.StringVar(&src, "src", "", "xlsx文件目录")
	flag.StringVar(&dst, "dst", "", "生成json目录")
	flag.Parse()

	// 读取所有xlsx文件
	files, err := util.Glob(src, "*.xlsx", true)
	if err != nil {
		panic(err)
	}

	// 优先解析define.xlsx文件
	for _, filename := range files {
		if filepath.Base(filename) == "define.xlsx" {
			manager.ParseDefine(filename)
			break
		}
	}
	for _, filename := range files {
		if filepath.Base(filename) == "define.xlsx" {
			continue
		}
		// 解析代对表
		if fb, err := manager.ParseProxy(filename); err != nil {
			//panic(err)
			continue
			// 解析配置表
		} else if err := fb.ParseTable(dst); err != nil {
			panic(err)
		}
	}
}
