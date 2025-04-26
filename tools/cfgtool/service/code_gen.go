package service

import (
	"bytes"
	"path/filepath"
	"sort"
	"strings"
	"universal/library/baselib/uerror"
	"universal/library/file"
	"universal/tools/cfgtool/domain"
	"universal/tools/cfgtool/internal/base"
	"universal/tools/cfgtool/internal/manager"
	"universal/tools/cfgtool/internal/templ"

	"github.com/iancoleman/strcase"
)

type ConfigInfo struct {
	PbImport string
	PbPkg    string
	Pkg      string
	*base.Config
}

type IndexInfo struct {
	Pkg       string
	IndexList []int
}

func GenCode(buf *bytes.Buffer) error {
	if len(domain.PbPath) <= 0 || len(domain.CodePath) <= 0 || len(domain.Module) <= 0 {
		return nil
	}

	// 生成索引
	if err := genIndex(buf); err != nil {
		return err
	}

	pbImport := filepath.Join(domain.Module, domain.PbPath)
	pkg := filepath.Base(domain.PbPath)

	// 对文件分类
	for _, st := range manager.GetConfigMap() {
		buf.Reset()
		dataName := strings.TrimSuffix(st.Name, "Config")
		name := strcase.ToSnake(dataName)
		item := &ConfigInfo{
			PbImport: pbImport,
			PbPkg:    pkg,
			Pkg:      strcase.ToSnake(name),
			Config:   st,
		}
		if err := templ.CodeTpl.Execute(buf, item); err != nil {
			return err
		}
		// 保存代码
		if err := file.SaveGo(filepath.Join(domain.CodePath, name), dataName+"Data.gen.go", buf.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

func genIndex(buf *bytes.Buffer) error {
	indexs := &IndexInfo{
		Pkg:       filepath.Base(domain.PbPath),
		IndexList: manager.GetIndexMap(),
	}

	if len(indexs.IndexList) > 0 {
		sort.Slice(indexs.IndexList, func(i, j int) bool {
			return indexs.IndexList[i] < indexs.IndexList[j]
		})

		buf.Reset()
		if err := templ.IndexTpl.Execute(buf, indexs); err != nil {
			return uerror.New(1, -1, "gen index file error: %s", err.Error())
		}
		return file.SaveGo(domain.PbPath, "index.gen.go", buf.Bytes())
	}
	return nil
}
