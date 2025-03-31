package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"universal/framework/uerror"
	"universal/tools/xlsx/domain"
	"universal/tools/xlsx/internal/base"
	"universal/tools/xlsx/internal/manager"

	"github.com/xuri/excelize/v2"
)

// 解析xlsx文件
func ParseXlsx(files ...string) error {
	for _, fileName := range files {
		fp, err := excelize.OpenFile(fileName)
		if err != nil {
			return uerror.NewUError(1, -1, "打开文件%s失败: %v", fileName, err)
		}

		cols, err := fp.GetCols(domain.GENERATE_TABLE)
		if err != nil {
			return uerror.NewUError(1, -1, "获取列失败: %v", err)
		}

		for _, vals := range cols {
			table := parseTable(fp, vals[0])
			manager.AddTable(table)
			for _, val := range vals[1:] {
				manager.AddEnum(parseValue(table, val))
			}
		}
	}
	return nil
}

// 解析结构类型
func ParseType() error {
	// 解析所有枚举
	for _, table := range manager.GetTables(domain.TYPE_OF_ENUM) {
		cols, err := table.GetCols()
		if err != nil {
			return uerror.NewUError(1, -1, "获取列失败: %v", err)
		}
		for _, vals := range cols {
			for _, val := range vals {
				if len(val) <= 0 {
					continue
				}
				manager.AddEnum(parseValue(table, val))
			}
		}
	}
	// 解析所有结构
	for _, table := range manager.GetTables(domain.TYPE_OF_STRUCT) {
		rows, err := table.GetRows()
		if err != nil {
			return uerror.NewUError(1, -1, "获取行失败: %v", err)
		}
		st := parseStruct(table, rows[:3])
		manager.AddStruct(st)
		for _, vals := range rows[3:] {
			parseStructConvert(st, vals...)
		}
	}
	// 解析所有配置结构
	for _, table := range manager.GetTables(domain.TYPE_OF_CONFIG) {
		rows, err := table.ScanRows(3)
		if err != nil {
			return uerror.NewUError(1, -1, "获取行失败: %v", err)
		}
		manager.AddConfig(parseConfig(table, rows))
	}
	// 文件分类
	manager.UpdateFile()
	return nil
}

func SaveJson(jspath string, buf *bytes.Buffer) error {
	for _, table := range manager.GetTables(domain.TYPE_OF_CONFIG) {
		rows, err := table.GetRows()
		if err != nil {
			return uerror.NewUError(1, -1, "获取行失败: %v", err)
		}

		cfg := manager.GetConfig(table.FileName)
		rets := []map[string]interface{}{}
		for _, row := range rows[3:] {
			rets = append(rets, ConvertConfig(cfg, row...))
		}

		jsData, err := json.MarshalIndent(rets, "", "  ")
		if err != nil {
			return uerror.NewUError(1, -1, "json.Marshal: %v", err)
		}
		buf.Reset()
		buf.Write(jsData)

		if err := base.SaveFile(path.Join(jspath, fmt.Sprintf("%s.json", table.FileName)), buf.Bytes()); err != nil {
			return uerror.NewUError(1, -1, "保存文件失败: %v", err)
		}
	}
	return nil
}
