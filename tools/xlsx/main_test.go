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

	if err := parser.ParseType(); err != nil {
		panic(err)
	}

	// 生成json
	if err := parser.ParseAndSaveJson("./json_tmp", bytes.NewBuffer(nil)); err != nil {
		panic(err)
	}
}
