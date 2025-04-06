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
	}
	manager.AddConfig(item)

	tmps := map[string]*base.Field{}
	for i, val := range rows[1] {
		if len(val) <= 0 {
			continue
		}
		valueOf := uint32(domain.VALUE_OF_IDENT)
		if strings.HasPrefix(val, "[]") {
			valueOf = domain.VALUE_OF_ARRAY
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

	for _, val := range tab.Rules {
		strs := strings.Split(val, ":")
		keys := []*base.Field{}
		for _, field := range strings.Split(strs[1], ",") {
			keys = append(keys, tmps[field])
		}

		if len(keys) > 0 {
			switch strings.ToLower(strs[0]) {
			case "map":
				item.MapList = append(item.MapList, keys)
			case "group":
				item.GroupList = append(item.GroupList, keys)
			}
		}
	}
	return nil
}
