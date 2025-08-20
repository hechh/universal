package main

import (
	"flag"
	"universal/library/fileutil"
	"universal/tool/pbtool/domain"
	"universal/tool/pbtool/internal/parse"
)

func main() {
	flag.StringVar(&domain.PbPath, "pb", "", ".pb.go文件目录")
	flag.Parse()

	if len(domain.PbPath) <= 0 {
		panic(".pb.go文件目录为空")
	}

	files, err := fileutil.Glob(domain.PbPath, ".*\\.pb\\.go", true)
	if err != nil {
		panic(err)
	}

	fac := parse.NewTypeParser()
	if err := fileutil.ParseFiles(fac, files...); err != nil {
		panic(err)
	}

	par := parse.NewGoParser(fac, "state", "sizeCache", "unknownFields")
	if err := fileutil.ParseFiles(par, files...); err != nil {
		panic(err)
	}
}
