package parser

import (
	"hego/Library/uerror"
	"hego/tools/xlsx/internal/base"
	"hego/tools/xlsx/internal/manager"
)

func ParseData(tab *base.Table) (rets []interface{}, err error) {
	rows, err := tab.Fp.GetRows(tab.SheetName)
	if err != nil {
		return nil, uerror.New(1, -1, "获取行失败: %v", err)
	}

	cfg := manager.GetConfig(tab.TypeName)
	for _, row := range rows[3:] {
		if len(row) <= 0 {
			continue
		}
		rets = append(rets, ConvertConfig(cfg, row...))
	}
	return
}
