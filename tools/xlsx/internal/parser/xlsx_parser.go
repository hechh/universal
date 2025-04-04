package parser

import (
	"fmt"
	"hego/Library/basic"
	"hego/Library/uerror"
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
					SheetName: basic.Ifelse(pos > 0, basic.GetPrefix(strs[1], pos), strs[1]),
					FileName:  basic.Ifelse(pos > 0, basic.GetSuffix(strs[1], pos+1), fileName),
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
					Rules:     basic.Suffix(strs, 2),
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
				ValueOf: domain.VALUE_OF_IDENT,
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
		switch strings.ToLower(strs[0]) {
		case "map":
			for _, field := range strings.Split(strs[1], ",") {
				item.Map = append(item.Map, tmps[field])
			}
		case "group":
			for _, field := range strings.Split(strs[1], ",") {
				item.Group = append(item.Group, tmps[field])
			}
		}
	}
	return nil
}

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
