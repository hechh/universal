package main

import (
	"bytes"
	"flag"
	"path/filepath"
	"universal/library/util"
	"universal/tool/pbtool/domain"
	"universal/tool/pbtool/internal/manager"
	"universal/tool/pbtool/internal/parse"
	"universal/tool/pbtool/internal/templ"
)

func main() {
	flag.StringVar(&domain.PbPath, "pb", "", ".pb.go文件目录")
	flag.Parse()

	if len(domain.PbPath) <= 0 {
		panic(".pb.go文件目录为空")
	}

	files, err := util.Glob(domain.PbPath, ".*\\.pb\\.go", true)
	if err != nil {
		panic(err)
	}

	fac := parse.NewTypeParser()
	if err := util.ParseFiles(fac, files...); err != nil {
		panic(err)
	}

	par := parse.NewGoParser(fac, "state", "sizeCache", "unknownFields")
	if err := util.ParseFiles(par, files...); err != nil {
		panic(err)
	}

	buff := bytes.NewBuffer(nil)
	if err := templ.PackageTpl.Execute(buff, filepath.Base(domain.PbPath)); err != nil {
		panic(err)
	}
	for _, cls := range manager.GetAll() {
		if err := templ.MemberTpl.Execute(buff, cls); err != nil {
			panic(err)
		}
	}
	if err := templ.FactoryTpl.Execute(buff, manager.GetAll()); err != nil {
		panic(err)
	}
	if err := util.SaveGo(domain.PbPath, "pb.gen.go", buff.Bytes()); err != nil {
		panic(err)
	}
}
