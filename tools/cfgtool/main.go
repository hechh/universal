package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
)

var (
	CwdPath string
)

func init() {
	var err error
	if CwdPath, err = os.Getwd(); err != nil {
		panic(err)
	}
}

func main() {
	var pb, xlsx string
	flag.StringVar(&pb, "pb", "", "go文件或文件目录")
	flag.StringVar(&xlsx, "xlsx", "", "xlsx文件或者文件目录")
	flag.Parse()

	cell, err := excelize.CoordinatesToCellName(1, 1)
	fmt.Println(cell, "----->", err)
	ee, err := excelize.ColumnNumberToName(2)
	fmt.Println(ee, "----->", err)
}
