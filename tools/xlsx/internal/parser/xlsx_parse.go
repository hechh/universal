package parser

import (
	"fmt"
	"hego/Library/uerror"
	"hego/Library/util"
	"hego/tools/xlsx/domain"
	"hego/tools/xlsx/internal/base"
	"hego/tools/xlsx/internal/manager"
	"path/filepath"
	"strings"

	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
)

func ParseXlsx(filename string) error {
	fp, err := excelize.OpenFile(filename)
	if err != nil {
		return uerror.New(1, -1, "打开文件%s失败: %v", filename, err)
	}

	cols, err := fp.GetCols("生成表")
	if err != nil {
		return uerror.New(1, -1, "获取列失败: %v", err)
	}

	fileName := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	for _, vals := range cols {
		for _, val := range vals {
			if len(val) <= 0 {
				continue
			}
			strs := strings.Split(val, "|")
			pos := strings.Index(strs[1], "@")
			switch strings.ToLower(strs[0]) {
			case "enum":
				manager.AddTable(&base.Table{
					TypeOf:    domain.TYPE_OF_ENUM,
					SheetName: util.Ifelse(pos > 0, util.GetPrefix(strs[1], pos), strs[1]),
					FileName:  util.Ifelse(pos > 0, util.GetSuffix(strs[1], pos+1), fileName),
					Fp:        fp,
				})
			case "struct":
				manager.AddTable(&base.Table{
					TypeOf:    domain.TYPE_OF_STRUCT,
					SheetName: strs[1][:pos],
					TypeName:  strs[1][pos+1:],
					FileName:  fileName,
					Fp:        fp,
				})
			case "config":
				manager.AddTable(&base.Table{
					TypeOf:    domain.TYPE_OF_CONFIG,
					SheetName: strs[1][:pos],
					TypeName:  strs[1][pos+1:],
					FileName:  fileName,
					Rules:     util.Suffix(strs, 2),
					Fp:        fp,
				})
			case "e":
				manager.AddEnum(&base.Value{
					TypeOf:   domain.TYPE_OF_ENUM,
					Type:     strs[2],
					Name:     fmt.Sprintf("%s_%s", strs[2], strs[3]),
					Value:    cast.ToUint32(strs[4]),
					Desc:     strs[1],
					FileName: fileName,
				})
			}
		}
	}
	return nil
}
