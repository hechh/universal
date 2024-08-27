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
	files, err := util.Search(src, "*.xlsx")
	if err != nil {
		panic(err)
	}
	for i := 1; i < len(files); i++ {
		if filepath.Base(files[i]) == "define.xlsx" {
			tmp := files[0]
			files[0] = files[i]
			files[i] = tmp
			break
		}
	}
	// 解析所有xlsx文件
	if err := manager.ParseXlsx(dst, files...); err != nil {
		panic(err)
	}
}
