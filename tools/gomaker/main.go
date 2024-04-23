package main

import (
	"path/filepath"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/parse"
	"universal/tools/gomaker/repository/uerrors"
)

func main() {
	par := parse.NewTypeParser(manager.GetTypeMgr())

	// 构建指定目录下所有文件的匹配模式
	files, err := filepath.Glob(filepath.Join("../../common/pb/", "*.pb.go"))
	if err != nil {
		panic(err)
	}
	if err = par.ParseFiles(files...); err != nil {
		panic(err)
	}

	// 生成文件
	if err := uerrors.Gen(); err != nil {
		panic(err)
	}
}
