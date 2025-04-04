package main

import (
	"bytes"
	"flag"
	"hego/framework/basic"
	"hego/tools/xlsx/internal/parser"
)

var (
	jsonPath string
	cfgPath  string
	xlsx     string
)

func main() {
	flag.StringVar(&jsonPath, "json", "", "json文件目录")
	flag.StringVar(&cfgPath, "cfg", "", "配置文件目录")
	flag.StringVar(&xlsx, "xlsx", "", "xlsx文件目录")
	flag.Parse()

	// 读取文件
	files, err := basic.Glob(xlsx, ".*\\.xlsx", "", true)
	if err != nil {
		panic(err)
	}
	// 解析table结构
	if err := parser.ParseXlsx(files...); err != nil {
		panic(err)
	}

	// 解析结构
	buf := bytes.NewBuffer(nil)
	if err := parser.SaveType(cfgPath, buf); err != nil {
		panic(err)
	}

	// 生成json
	if err := parser.SaveJson(jsonPath, buf); err != nil {
		panic(err)
	}
}
