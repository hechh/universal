package parser

import (
	"fmt"
	"strings"
	"universal/tools/xlsx/domain"
	"universal/tools/xlsx/internal/base"
	"universal/tools/xlsx/internal/manager"

	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
)

// enum|CST|全局表@Global
func parseTable(fp *excelize.File, str string) *base.Table {
	strs := strings.Split(str, "|")
	index := strings.Index(strs[2], "@")
	item := &base.Table{
		SheetName: strs[2][:index],
		FileName:  strs[2][index+1:],
	}
	item.SetFp(fp)
	switch strings.ToLower(strs[0]) {
	case domain.TYPE_OF_ENUM_NAME:
		item.TypeOf = domain.TYPE_OF_ENUM
	case domain.TYPE_OF_STRUCT_NAME:
		item.TypeOf = domain.TYPE_OF_STRUCT
	case domain.TYPE_OF_CONFIG_NAME:
		item.TypeOf = domain.TYPE_OF_CONFIG
	}
	return item
}

// E:服务类型-game:ServerType:Game:1
func parseValue(table *base.Table, str string) *base.Value {
	strs := strings.Split(str, ":")
	switch strings.ToLower(strs[0]) {
	case "e":
		return &base.Value{
			TypeOf:    domain.TYPE_OF_ENUM,
			Type:      strs[2],
			Name:      fmt.Sprintf("%s_%s", strs[2], strs[3]),
			Value:     cast.ToUint32(strs[4]),
			Desc:      strs[1],
			SheetName: table.SheetName,
			FileName:  table.FileName,
		}
	}
	return nil
}

func parseStruct(table *base.Table, vals [][]string) *base.Struct {
	ret := &base.Struct{
		Name:      table.FileName,
		SheetName: table.SheetName,
		FileName:  table.FileName,
		Converts:  map[string][]*base.Field{},
	}
	for i, val := range vals[1] {
		if len(val) <= 0 {
			continue
		}
		ret.List = append(ret.List, &base.Field{
			Type: &base.Type{
				Name:    val,
				TypeOf:  manager.GetTypeOf(val),
				ValueOf: domain.VALUE_OF_IDENT,
			},
			Name:     vals[0][i],
			Desc:     vals[2][i],
			Position: i,
		})
	}
	return ret
}

func parseStructConvert(st *base.Struct, vals ...string) {
	for i, val := range vals {
		if len(val) <= 0 || val == "0" {
			continue
		}
		st.Converts[vals[0]] = append(st.Converts[vals[0]], st.List[i])
	}
}

func parseConfig(table *base.Table, vals [][]string) *base.Config {
	ret := &base.Config{
		Name:      table.FileName,
		SheetName: table.SheetName,
		FileName:  table.FileName,
	}
	for i, val := range vals[1] {
		if len(val) <= 0 {
			continue
		}
		valueOf := uint32(domain.VALUE_OF_IDENT)
		if strings.HasPrefix(val, "[]") {
			valueOf = domain.VALUE_OF_ARRAY
			val = strings.TrimPrefix(val, "[]")
		}
		ret.List = append(ret.List, &base.Field{
			Type: &base.Type{
				Name:    val,
				TypeOf:  manager.GetTypeOf(val),
				ValueOf: valueOf,
			},
			Name:     vals[0][i],
			Desc:     vals[2][i],
			Position: i,
		})
	}
	return ret
}
