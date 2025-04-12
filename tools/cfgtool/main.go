package main

import (
	"bytes"
	"flag"
	"hego/Library/basic"
	"hego/tools/cfgtool/domain"
	"hego/tools/cfgtool/internal/manager"
	"hego/tools/cfgtool/internal/parser"
	"hego/tools/cfgtool/service"
	"path/filepath"
)

func main() {
	flag.StringVar(&domain.XlsxPath, "xlsx", ".", "cfg文件目录")
	flag.StringVar(&domain.DataPath, "data", "", "数据文件目录")
	flag.StringVar(&domain.ProtoPath, "proto", "./cfg_protocol", "proto文件目录")
	flag.StringVar(&domain.CodePath, "code", "", "go代码文件目录")
	flag.StringVar(&domain.Module, "module", "", "项目目录")
	flag.StringVar(&domain.PbPath, "pb", "", "proto生成路径")
	flag.Parse()

	if len(domain.XlsxPath) <= 0 {
		panic("配置文件目录不能为空")
	}
	domain.ProtoPkgName = filepath.Base(domain.ProtoPath)
	if len(domain.PbPath) > 0 {
		domain.PkgName = filepath.Base(domain.PbPath)
	}

	// 加载所有配置
	files, err := basic.Glob(domain.XlsxPath, ".*\\.xlsx", "", true)
	if err != nil {
		panic(err)
	}

	// 解析所有文件
	if err := parser.ParseFiles(files...); err != nil {
		panic(err)
	}

	// 生成proto文件数据
	buf := bytes.NewBuffer(nil)
	if err := service.GenProto(domain.ProtoPath, buf); err != nil {
		panic(err)
	}
	if err := service.SaveProto(domain.ProtoPath); err != nil {
		panic(err)
	}

	// 解析proto文件
	if err := manager.ParseProto(); err != nil {
		panic(err)
	}

	if err := service.GenData(domain.DataPath, buf); err != nil {
		panic(err)
	}
}
