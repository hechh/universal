package parse

import (
	"strings"
	"universal/framework/uerror"
	"universal/tools/xlsx/internal/manager"

	"github.com/xuri/excelize/v2"
)

const (
	GENERATE_TABLE = "生成表"
)

func OpenFiles(files ...string) error {
	for _, fileName := range files {
		fp, err := excelize.OpenFile(fileName)
		if err != nil {
			return uerror.NewUError(1, -1, "打开文件%s失败: %v", fileName, err)
		}

		cols, err := fp.GetCols(GENERATE_TABLE)
		if err != nil {
			return uerror.NewUError(1, -1, "获取列失败: %v", err)
		}

		for _, vals := range cols {
			for _, val := range vals {
				if strings.HasPrefix(val, "E:") {
					manager.AddEnum(val)
				} else {
					manager.AddTable(fp, val)
				}
			}
		}
	}
	return nil
}
