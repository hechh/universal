package parse

import (
	"bytes"
	"strings"
	"universal/framework/uerror"
	"universal/tools/xlsx/domain"
	"universal/tools/xlsx/internal/base"
	"universal/tools/xlsx/internal/manager"

	"github.com/spf13/cast"
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
			table := parseTable(fp, vals[0])
			manager.AddTable(table)
			for _, val := range vals[1:] {
				manager.AddEnum(table, val)
			}
		}
	}
	return nil
}

func parseTable(fp *excelize.File, str string) *base.Table {
	strs := strings.Split(str, "|")
	index := strings.Index(strs[2], "@")
	item := &base.Table{
		Sheet:    strs[2][:index],
		FileName: strs[2][index+1:],
		Fp:       fp,
	}
	switch strs[0] {
	case "enum":
		item.Priority = 2
		item.TypeOf = domain.TypeOfEnum
	case "struct":
		item.Priority = 1
		item.TypeOf = domain.TypeOfStruct
	case "config":
		item.TypeOf = domain.TypeOfConfig
	}
	return item
}

func Parse() error {
	buf := bytes.NewBuffer(nil)
	for _, table := range manager.GetTableList() {
		switch table.TypeOf {
		case domain.TypeOfEnum:
			cols, err := table.Fp.GetCols(table.Sheet)
			if err != nil {
				return uerror.NewUError(1, -1, "获取列失败: %v", err)
			}
			for _, vals := range cols {
				for _, val := range vals {
					if len(val) <= 0 {
						continue
					}
					manager.AddEnum(table, val)
				}
			}
		case domain.TypeOfStruct:
			rows, err := table.Fp.GetRows(table.Sheet)
			if err != nil {
				return uerror.NewUError(1, -1, "获取行失败: %v", err)
			}
			st := parseStruct(table, rows[0], rows[1], rows[2])
			for _, vals := range rows[3:] {
				parseStructConvert(st, vals...)
			}
		case domain.TypeOfConfig:
			rows, err := table.Fp.GetRows(table.Sheet)
			if err != nil {
				return uerror.NewUError(1, -1, "获取行失败: %v", err)
			}
			cf := parseConfig(table, rows[0], rows[1], rows[2])
			manager.AddConfig(cf)
			manager.ParseRows(cf, rows[3:], buf)
		}
		table.Fp.Close()
	}
	return nil
}

func parseConfig(st *base.Table, val0, val1, val2 []string) *base.Config {
	ret := &base.Config{
		Name:     st.FileName,
		Sheet:    st.Sheet,
		FileName: st.FileName,
	}
	for i, val := range val1 {
		if len(val) <= 0 {
			continue
		}
		valueOf := uint32(domain.ValueOfSingle)
		typeOf := uint32(domain.TypeOfBase)
		if strings.HasPrefix(val, "[]") {
			valueOf = domain.ValueOfArray
			val = strings.TrimPrefix(val, "[]")
		}
		if manager.IsEnum(val) {
			typeOf = domain.TypeOfEnum
		} else if manager.IsStruct(val) {
			typeOf = domain.TypeOfStruct
		}
		ret.List = append(ret.List, &base.Field{
			Name: val0[i],
			Type: &base.Type{
				Name:    val,
				TypeOf:  typeOf,
				ValueOf: valueOf,
			},
			Desc:     val2[i],
			Position: i,
		})
	}
	return ret
}

func parseStructConvert(st *base.Struct, vals ...string) {
	for i, val := range vals {
		if len(val) <= 0 || cast.ToInt(val) <= 0 {
			continue
		}
		st.Converts[vals[0]] = append(st.Converts[vals[0]], st.List[i])
	}
}

func parseStruct(st *base.Table, val0, val1, val2 []string) *base.Struct {
	ret := &base.Struct{
		Name:     st.FileName,
		Sheet:    st.Sheet,
		FileName: st.FileName,
		Converts: map[string][]*base.Field{},
	}
	for i, val := range val1 {
		if len(val) <= 0 {
			continue
		}
		typeOf := uint32(domain.TypeOfBase)
		if manager.IsEnum(val) {
			typeOf = domain.TypeOfEnum
		}
		ret.List = append(ret.List, &base.Field{
			Name: val0[i],
			Type: &base.Type{
				Name:    val,
				TypeOf:  typeOf,
				ValueOf: domain.ValueOfSingle,
			},
			Desc:     val2[i],
			Position: i,
		})
	}
	return ret
}
