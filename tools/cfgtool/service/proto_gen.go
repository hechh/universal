package service

import (
	"bytes"
	"hego/Library/file"
	"hego/Library/uerror"
	"hego/tools/cfgtool/domain"
	"hego/tools/cfgtool/internal/base"
	"hego/tools/cfgtool/internal/manager"
	"hego/tools/cfgtool/internal/templ"
	"path/filepath"
	"sort"
)

type ProtoInfo struct {
	Pkg        string
	GoPkg      string
	RefList    []string
	EnumList   []*base.Enum
	StructList []*base.Struct
	ConfigList []*base.Config
}

func GenProto(protoPath string, buf *bytes.Buffer) error {
	pkg := filepath.Base(protoPath)
	// 根据文件分类
	tmps := map[string]*ProtoInfo{}
	for _, val := range manager.GetEnumList() {
		sort.Slice(val.ValueList, func(i, j int) bool {
			return val.ValueList[i].Value < val.ValueList[j].Value
		})
		if _, ok := tmps[val.FileName]; !ok {
			tmps[val.FileName] = &ProtoInfo{Pkg: pkg, GoPkg: domain.PkgName}
		}
		tmps[val.FileName].EnumList = append(tmps[val.FileName].EnumList, val)
	}

	for _, val := range manager.GetStructList() {
		if _, ok := tmps[val.FileName]; !ok {
			tmps[val.FileName] = &ProtoInfo{Pkg: pkg, GoPkg: domain.PkgName}
		}
		tmps[val.FileName].StructList = append(tmps[val.FileName].StructList, val)
	}

	for _, val := range manager.GetConfigList() {
		if _, ok := tmps[val.FileName]; !ok {
			tmps[val.FileName] = &ProtoInfo{Pkg: pkg, GoPkg: domain.PkgName}
		}
		tmps[val.FileName].ConfigList = append(tmps[val.FileName].ConfigList, val)
	}

	// 生成proto文件
	for fileName, data := range tmps {
		buf.Reset()
		data.RefList = manager.GetRefList(fileName)
		if err := templ.ProtoTpl.Execute(buf, data); err != nil {
			return uerror.New(1, -1, "gen proto file error: %s", err.Error())
		}
		manager.AddProto(fileName, buf)
	}
	return nil
}

func SaveProto(protoPath string) error {
	for fileName, data := range manager.GetProtoMap() {
		if err := file.Save(protoPath, fileName, []byte(data)); err != nil {
			return err
		}
	}
	return nil
}
