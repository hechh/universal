package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"hego/Library/basic"
	"hego/Library/file"
	"hego/tools/xlsx/domain"
	"hego/tools/xlsx/internal/base"
	"hego/tools/xlsx/internal/generate"
	"hego/tools/xlsx/internal/manager"
	"hego/tools/xlsx/internal/parser"
	"path/filepath"
)

func main() {
	var jsonPath, cfgPath, codePath, xlsxPath string
	flag.StringVar(&jsonPath, "json", "", "json文件目录")
	flag.StringVar(&cfgPath, "cfg", "", "配置文件目录")
	flag.StringVar(&codePath, "code", "", "代码文件目录")
	flag.StringVar(&xlsxPath, "xlsx", "", "xlsx文件目录")
	flag.Parse()
	domain.PkgName = filepath.Base(cfgPath)

	// 读取文件
	files, err := basic.Glob(xlsxPath, ".*\\.xlsx", "", true)
	if err != nil {
		panic(err)
	}

	// 解析table结构
	for _, fileName := range files {
		if err := parser.ParseXlsx(fileName); err != nil {
			panic(err)
		}
	}

	// 解析结构
	parseType()
	// 保存结构
	saveType(cfgPath)
	saveCode(codePath)
	// 解析数据，并且保存
	saveJson(jsonPath)
}

func parseType() {
	for _, table := range manager.GetTableList() {
		switch table.TypeOf {
		case domain.TypeOfEnum:
			if err := parser.ParseEnum(table); err != nil {
				panic(err)
			}
		case domain.TypeOfStruct:
			if err := parser.ParseStruct(table); err != nil {
				panic(err)
			}
		case domain.TypeOfConfig:
			if err := parser.ParseConfig(table); err != nil {
				panic(err)
			}
		}
	}
}

func saveType(cfgPath string) {
	buf := bytes.NewBuffer(nil)
	for fileName, items := range manager.GetFileInfo() {
		buf.WriteString(fmt.Sprintf("package %s\n", filepath.Base(cfgPath)))
		for _, item := range items {
			switch val := item.(type) {
			case *base.Enum:
				val.Format(buf)
			case *base.Struct:
				val.Format(buf)
			case *base.Config:
				val.Format(buf)
			}
		}
		if err := file.SaveGo(cfgPath, fmt.Sprintf("%s.gen.go", fileName), buf.Bytes()); err != nil {
			panic(err)
		}
		buf.Reset()
	}
}

func saveCode(codePath string) {
	buf := bytes.NewBuffer(nil)
	manager.WalkConfig(func(cfg *base.Config) bool {
		if err := generate.Generate(codePath, cfg, buf); err != nil {
			panic(err)
		}
		buf.Reset()
		return true
	})
}

func saveJson(jsonPath string) {
	for _, table := range manager.GetTables(domain.TypeOfConfig) {
		data, err := parser.ParseData(table)
		if err != nil {
			panic(err)
		}
		jsbuf, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			panic(err)
		}
		if err := file.Save(jsonPath, fmt.Sprintf("%s.json", table.FileName), jsbuf); err != nil {
			panic(err)

		}
	}
}
