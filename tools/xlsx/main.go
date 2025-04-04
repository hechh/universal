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
	"hego/tools/xlsx/internal/manager"
	"hego/tools/xlsx/internal/parser"
)

var (
	jsonPath string
	cfgPath  string
	codePath string
	xlsx     string
)

func main() {
	flag.StringVar(&jsonPath, "json", "", "json文件目录")
	flag.StringVar(&cfgPath, "cfg", "", "配置文件目录")
	flag.StringVar(&codePath, "code", "", "代码文件目录")
	flag.StringVar(&xlsx, "xlsx", "", "xlsx文件目录")
	flag.Parse()

	// 读取文件
	files, err := basic.Glob(xlsx, ".*\\.xlsx", "", true)
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
	for _, table := range manager.GetTableList() {
		switch table.TypeOf {
		case domain.TYPE_OF_ENUM:
			if err := parser.ParseEnum(table); err != nil {
				panic(err)
			}
		case domain.TYPE_OF_STRUCT:
			if err := parser.ParseStruct(table); err != nil {
				panic(err)
			}
		case domain.TYPE_OF_CONFIG:
			if err := parser.ParseConfig(table); err != nil {
				panic(err)
			}
		}
	}

	// 保存结构
	buf := bytes.NewBuffer(nil)
	for fileName, items := range manager.GetFileInfo() {
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

	// 解析数据，并且保存
	for _, table := range manager.GetTables(domain.TYPE_OF_CONFIG) {
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
