package main

import (
	"bytes"
	"flag"
	"universal/library/util"
	"universal/tools/cfgtool/domain"
	"universal/tools/cfgtool/internal/parser"
	"universal/tools/cfgtool/service"
)

func main() {
	flag.StringVar(&domain.XlsxPath, "xlsx", "./xlsx", "cfg文件目录")
	flag.StringVar(&domain.ProtoPath, "proto", "", "proto文件目录")
	flag.StringVar(&domain.TextPath, "text", "", "数据文件目录")
	flag.StringVar(&domain.PbPath, "pb", "", "proto生成路径")
	flag.Parse()

	if len(domain.XlsxPath) <= 0 {
		panic("配置文件目录不能为空")
	}

	// 加载所有配置
	files, err := util.Glob(domain.XlsxPath, ".*\\.xlsx", true)
	if err != nil {
		panic(err)
	}
	// 解析所有文件
	if err := parser.ParseFiles(files...); err != nil {
		panic(err)
	}
	// 生成proto文件数据
	buf := bytes.NewBuffer(nil)
	if err := service.GenProto(buf); err != nil {
		panic(err)
	}
	if err := service.GenData(); err != nil {
		panic(err)
	}
}
