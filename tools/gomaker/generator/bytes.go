package generator

import (
	"text/template"
	"unicode"
	"universal/framework/uerror"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/parse"
	"universal/tools/gomaker/internal/typespec"

	"github.com/xuri/excelize/v2"
)

// 生成enum.gen.proto
func bytesGenerator(dst string, tpls *template.Template, extra ...string) error {
	for _, filename := range extra {
		if err := parseXlsx(filename); err != nil {
			return err
		}
	}
	return nil
}

func parseXlsx(filename string) error {
	fb, err := excelize.OpenFile(filename)
	if err != nil {
		return uerror.NewUError(1, -1, "%v", err)
	}
	defer fb.Close()
	// 根据生成表解析config结构
	for k, v := range parse.GetTables(fb) {
		values, _ := fb.GetRows(k)
		item := parse.ParseXlsxStruct(v, values[0], values[1])
		if len(item.List) <= 0 {
			continue
		}
		// 加载pb.go文件中的结构信息
		pbItem := manager.GetStruct(item.Type.GetPkgType())
		if pbItem == nil {
			return uerror.NewUError(1, -1, "配置表中未定义%s结构", v)
		}
		fields := map[int]*typespec.Field{}
		for _, ff := range pbItem.List {
			if unicode.IsLower(rune(ff.Name[0])) {
				continue
			}
			ff.Index = item.Fields[ff.Name].Index
			fields[ff.Index] = ff
		}
		// 解析配置数据
		/*
			for _, vals := range values[2:] {

			}
		*/
	}
	return nil
}
