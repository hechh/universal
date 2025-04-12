package test

import (
	"bytes"
	"hego/Library/basic"
	"hego/tools/cfgtool/domain"
	"hego/tools/cfgtool/internal/manager"
	"hego/tools/cfgtool/internal/parser"
	"hego/tools/cfgtool/service"
	"path/filepath"
	"testing"
)

func TestCfg(t *testing.T) {
	domain.XlsxPath = "../../../configure/table"
	domain.DataPath = "../../../configure/json"
	domain.ProtoPath = "../../../configure/proto"
	domain.PbPath = "../../../common/pb"
	domain.PkgName = filepath.Base(domain.PbPath)
	domain.ProtoPkgName = filepath.Base(domain.ProtoPath)

	// 加载所有配置
	files, err := basic.Glob(domain.XlsxPath, ".*\\.xlsx", "", true)
	if err != nil {
		panic(err)
	}

	// 解析所有文件
	if err := parser.ParseFiles(files...); err != nil {
		panic(err)
	}
	// 生成proto文件数据
	buf := bytes.NewBuffer(nil)
	if err := service.GenProto(domain.ProtoPath, buf); err != nil {
		panic(err)
	}
	if err := service.SaveProto(domain.ProtoPath); err != nil {
		panic(err)
	}

	// 解析proto文件
	if err := manager.ParseProto(); err != nil {
		panic(err)
	}

	if err := service.GenData(domain.DataPath, buf); err != nil {
		panic(err)
	}

}
