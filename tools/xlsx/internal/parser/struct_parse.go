package parser

import (
	"hego/Library/uerror"
	"hego/tools/xlsx/domain"
	"hego/tools/xlsx/internal/base"
	"hego/tools/xlsx/internal/manager"
)

func ParseStruct(tab *base.Table) error {
	rows, err := tab.Fp.GetRows(tab.SheetName)
	if err != nil {
		return uerror.New(1, -1, "获取行失败: %v", err)
	}

	item := &base.Struct{
		Name:     tab.TypeName,
		FileName: tab.FileName,
		Converts: map[string][]*base.Field{},
	}
	manager.AddStruct(item)

	for i, val := range rows[1] {
		if len(val) <= 0 {
			continue
		}
		item.List = append(item.List, &base.Field{
			Type: &base.Type{
				Name:    val,
				TypeOf:  manager.GetTypeOf(val),
				ValueOf: domain.ValueOfBase,
			},
			Name:     rows[0][i],
			Desc:     rows[2][i],
			Position: i,
		})
	}

	for _, vals := range rows[3:] {
		for i, val := range vals {
			if len(val) <= 0 || val == "0" {
				continue
			}
			item.Converts[vals[0]] = append(item.Converts[vals[0]], item.List[i])
		}
	}
	return nil
}
