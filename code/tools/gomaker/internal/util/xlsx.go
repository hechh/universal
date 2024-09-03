package util

import (
	"path/filepath"
	"strings"
	"universal/framework/basic/uerror"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/typespec"

	"github.com/xuri/excelize/v2"
)

// 解析xlsx
func ParseXlsx(v domain.Visitor, filename string) error {
	// 加载文件
	fb, err := excelize.OpenFile(filename)
	if err != nil {
		return err
	}
	defer fb.Close()
	v.SetFile(filename)

	// 解析内容
	for _, sheet := range fb.GetSheetList() {
		values, _ := fb.GetRows(sheet)
		if strings.HasPrefix(sheet, "@") {
			if ret := v.Visit(typespec.NewTableNode(sheet, values[0], values[1])); ret == nil {
				continue
			}
			for _, vals := range values[2:] {
				v.Visit(&typespec.ValueNode{SheetName: sheet, Values: vals})
			}
		} else if filepath.Base(filename) == "define.xlsx" || sheet == "代对表" {
			for _, vals := range values {
				for _, val := range vals {
					switch val[:strings.Index(val, ":")] {
					case "E", "e":
						v.Visit(&typespec.EnumNode{SheetName: sheet, Value: val})
					case "CS", "cs", "Cs", "cS", "SC", "Sc", "sC", "sc", "s", "S":
						v.Visit(&typespec.ProxyNode{SheetName: sheet, Value: val, IsCreator: true})
					case "c", "C":
						v.Visit(&typespec.ProxyNode{SheetName: sheet, Value: val})
					}
				}
			}
		}
	}
	return nil
}

func ParseDirXlsx(v domain.Visitor, dir string) error {
	files, err := Glob(dir, "*.xlsx", true)
	if err != nil {
		return err
	}
	for _, filename := range files {
		if err := ParseXlsx(v, filename); err != nil {
			return uerror.NewUError(1, -1, "%v", err)
		}
	}
	return nil
}
