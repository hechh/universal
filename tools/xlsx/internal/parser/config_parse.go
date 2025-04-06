package parser

import (
	"hego/tools/xlsx/domain"
	"hego/tools/xlsx/internal/base"
	"hego/tools/xlsx/internal/manager"
	"strings"
)

func ParseConfig(tab *base.Table) error {
	rows, err := tab.ScanRows(3)
	if err != nil {
		return err
	}
	item := &base.Config{
		Name:     tab.TypeName,
		FileName: tab.FileName,
		IndexList: []*base.Index{
			{Name: "list",
				Type: &base.Type{
					TypeOf:  domain.TypeOfBase,
					ValueOf: domain.ValueOfList,
				},
			},
		},
	}
	manager.AddConfig(item)

	tmps := map[string]*base.Field{}
	for i, val := range rows[1] {
		if len(val) <= 0 {
			continue
		}
		valueOf := uint32(domain.ValueOfBase)
		if strings.HasPrefix(val, "[]") {
			valueOf = domain.ValueOfList
			val = strings.TrimPrefix(val, "[]")
		}
		tmps[rows[0][i]] = &base.Field{
			Type: &base.Type{
				Name:    val,
				TypeOf:  manager.GetTypeOf(val),
				ValueOf: valueOf,
			},
			Name:     rows[0][i],
			Desc:     rows[2][i],
			Position: i,
		}
		item.List = append(item.List, tmps[rows[0][i]])
	}

	parseIndex(tab, item, tmps)
	return nil
}

func parseIndex(tab *base.Table, item *base.Config, tmps map[string]*base.Field) {
	for _, val := range tab.Rules {
		strs := strings.Split(val, ":")

		keys := []*base.Field{}
		fields := []string{}
		for _, field := range strings.Split(strs[1], ",") {
			keys = append(keys, tmps[field])
			fields = append(fields, field)
		}

		name := strings.Join(fields, "")
		if len(strs) > 2 {
			name = strs[2]
		}

		typeOf := uint32(domain.TypeOfStruct)
		tname := strings.Join(fields, "")
		if len(keys) == 1 {
			tname = keys[0].Type.GetType(tab.SheetName)
			typeOf = domain.TypeOfBase
		}
		valueOf := uint32(domain.ValueOfGroup)
		if strings.ToLower(strs[0]) == "map" {
			valueOf = domain.ValueOfMap
		}

		item.IndexList = append(item.IndexList, &base.Index{
			Name: name,
			Type: &base.Type{
				Name:    tname,
				TypeOf:  typeOf,
				ValueOf: valueOf,
			},
			List: keys,
		})
	}
}
