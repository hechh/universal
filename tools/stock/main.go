package main

import (
	"flag"
	"os"
	"path"
	"path/filepath"
	"stock/internal/manager"
	"strings"

	"github.com/xuri/excelize/v2"
)

func main() {
	// 解析参数
	var fpath, dpath, output string
	flag.StringVar(&fpath, "f", "", "题材配置表")
	flag.StringVar(&dpath, "d", "", "涨停数据表")
	flag.StringVar(&output, "o", "./output.xlsx", "热点输出")
	flag.Parse()
	// 获取工作目录
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// 获取绝对地址
	fpath = filepath.Clean(filepath.Join(cwd, fpath))
	dpath = filepath.Clean(filepath.Join(cwd, dpath))
	output = filepath.Clean(filepath.Join(cwd, output))
	// 判断
	if !strings.HasSuffix(output, ".xlsx") {
		output = path.Join(output, "output.xlsx")
	}
	// 解析数据
	if err := manager.ParseFilters(fpath); err != nil {
		panic(err)
	}
	if err := manager.ParseStocks(dpath); err != nil {
		panic(err)
	}
	// 分析数据
	manager.Analyse()
	fp := excelize.NewFile()
	defer fp.SaveAs(output)
	manager.Write(fp)
}
