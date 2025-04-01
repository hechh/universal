package main

import (
	"bytes"
	"testing"
	"universal/framework/basic"
	"universal/tools/xlsx/internal/parser"
)

func TestRun(t *testing.T) {
	//go run main.go -xlsx=. -json=./json_tmp/ -cfg=./cfg_tmp/
	files, err := basic.Glob("./", ".*\\.xlsx", "", true)
	if err != nil {
		panic(err)
	}

	// 解析table结构
	if err := parser.ParseXlsx(files...); err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(nil)
	if err := parser.SaveType("./json_tmp", buf); err != nil {
		panic(err)
	}

	// 生成json
	if err := parser.SaveJson("./json_tmp", buf); err != nil {
		panic(err)
	}
}
