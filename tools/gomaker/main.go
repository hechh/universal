package main

import (
	"flag"
	"path/filepath"
	"universal/tools/gomaker/internal/parse"
	"universal/tools/gomaker/repository/uerrors"
)

func main() {
	var src, dst, tpl string
	flag.StringVar(&tpl, "tpl", "", ".tpl模版文件")
	flag.StringVar(&src, "src", "", ".go文件路径")
	flag.StringVar(&dst, "dst", "", ".gen.go生成文件")
	flag.Parse()

	par := parse.NewTypeParser()

	// 构建指定目录下所有文件的匹配模式
	files, err := filepath.Glob(filepath.Join("../../common/pb/", "*.pb.go"))
	//files, err := filepath.Glob(src)
	if err != nil {
		panic(err)
	}
	if err = par.ParseFiles(files...); err != nil {
		panic(err)
	}

	// 生成文件
	if err := uerrors.Gen(); err != nil {
		panic(err)
	}
}
