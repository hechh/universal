package manager

import (
	"strings"
	"universal/tools/cfgtool/domain"
	"universal/tools/cfgtool/internal/typespec"
	"universal/tools/cfgtool/internal/util"
)

var (
	alls  = make(map[string]*typespec.Enum)
	files = make(map[string]*typespec.FileInfo)
)

func ParseXlsx(dst string, files ...string) error {
	// 解析define.xlsx文件
	info, err := typespec.NewFileInfo(files[0], alls)
	if err != nil {
		return err
	}
	for _, sheet := range info.GetFPointer().GetSheetList() {
		info.ParseProxy(sheet)
	}

	// 解析其他文件
	for _, filename := range files[1:] {
		info, err := typespec.NewFileInfo(filename, alls)
		if err != nil {
			return err
		}

		// 保证优先解析代对表
		info.ParseProxy(domain.ProxyTable)
		for _, sheet := range info.GetFPointer().GetSheetList() {
			if !strings.HasPrefix(sheet, "@") {
				continue
			}
			datas, err := info.ParseTable(sheet)
			if err != nil {
				return err
			}

			if tab := info.GetProxy(strings.TrimPrefix(sheet, "@")); tab != nil {
				filename := info.GetJsonName(dst, tab.EnglishName)
				if err := util.SaveJson(filename, datas); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
