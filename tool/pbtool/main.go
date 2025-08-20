package main

import (
	"bytes"
	"flag"
	"path/filepath"
	"universal/library/fileutil"
	"universal/tool/pbtool/domain"
	"universal/tool/pbtool/internal/manager"
	"universal/tool/pbtool/internal/parse"
	"universal/tool/pbtool/internal/tpl"
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

	buff := bytes.NewBuffer(nil)
	if err := tpl.PackageTpl.Execute(buff, filepath.Base(domain.PbPath)); err != nil {
		panic(err)
	}
	for _, cls := range manager.GetAll() {
		if err := tpl.MemberTpl.Execute(buff, cls); err != nil {
			panic(err)
		}
	}
	if err := tpl.FactoryTpl.Execute(buff, manager.GetAll()); err != nil {
		panic(err)
	}
	if err := fileutil.SaveGo(domain.PbPath, "pb.gen.go", buff.Bytes()); err != nil {
		panic(err)
	}
}
