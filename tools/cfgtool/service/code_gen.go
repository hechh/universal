package service

import (
	"bytes"
	"path"
	"sort"
	"universal/library/fileutil"
	"universal/library/uerror"
	"universal/tools/cfgtool/domain"
	"universal/tools/cfgtool/internal/base"
	"universal/tools/cfgtool/internal/manager"
	"universal/tools/cfgtool/internal/templ"

	"strings"

	"github.com/iancoleman/strcase"
)

type IndexInfo struct {
	Pkg       string
	IndexList []int
}

type ConfigInfo struct {
	PbPkg string
	Pkg   string
	*base.Config
}

func GenCode(buf *bytes.Buffer) error {
	if len(domain.PbPath) <= 0 || len(domain.CodePath) <= 0 {
		return nil
	}

	// 生成索引
	if err := genIndex(buf); err != nil {
		return err
	}
	// 对文件分类
	for _, st := range manager.GetConfigMap() {
		buf.Reset()
		dataName := strings.TrimSuffix(st.Name, "Config")
		name := strcase.ToSnake(st.Name)
		item := &ConfigInfo{
			PbPkg:  domain.PkgName,
			Pkg:    name,
			Config: st,
		}
		if err := templ.CodeTpl.Execute(buf, item); err != nil {
			return err
		}
		// 保存代码
		if err := fileutil.SaveGo(path.Join(domain.CodePath, name), dataName+"Data.gen.go", buf.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

func genIndex(buf *bytes.Buffer) error {
	indexs := &IndexInfo{
		Pkg:       "pb",
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
		return fileutil.SaveGo(domain.PbPath, "index.gen.go", buf.Bytes())
	}
	return nil
}
