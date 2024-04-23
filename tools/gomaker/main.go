package main

import (
	"path/filepath"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/visitor"
)

func main() {
	tmgr := manager.NewTypeMgr()
	par := visitor.NewTypeParser(tmgr)

	// 构建指定目录下所有文件的匹配模式
	files, err := filepath.Glob(filepath.Join("./internal/test/", "*.pb.go"))
	if err != nil {
		panic(err)
	}
	if err = par.ParseFiles(files...); err != nil {
		panic(err)
	}

	tmgr.Print()
}
