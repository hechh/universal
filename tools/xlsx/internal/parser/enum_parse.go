package parser

import (
	"fmt"
	"hego/Library/uerror"
	"hego/tools/xlsx/domain"
	"hego/tools/xlsx/internal/base"
	"hego/tools/xlsx/internal/manager"
	"strings"

	"github.com/spf13/cast"
)

func ParseEnum(tab *base.Table) error {
	cols, err := tab.Fp.GetCols(tab.SheetName)
	if err != nil {
		return uerror.New(1, -1, "获取列失败: %v", err)
	}

	for _, vals := range cols {
		for _, val := range vals {
			if len(val) <= 0 {
				continue
			}

			strs := strings.Split(val, "|")
			manager.AddEnum(&base.Value{
				TypeOf:   domain.TYPE_OF_ENUM,
				Type:     strs[2],
				Name:     fmt.Sprintf("%s_%s", strs[2], strs[3]),
				Value:    cast.ToUint32(strs[4]),
				Desc:     strs[1],
				FileName: tab.FileName,
			})
		}
	}
	return nil
}
