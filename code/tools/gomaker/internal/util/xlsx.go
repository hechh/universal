package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"universal/framework/basic/uerror"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/typespec"

	"github.com/xuri/excelize/v2"
)

// 需要规定sheet解析顺序
func sortSheet(list []string) []string {
	j := -1
	for i, name := range list {
		if name == "代对表" {
			j++
			list[i] = list[j]
			list[j] = name
			break
		}
	}
	return list
}

// 需要规定xlsx文件解析顺序
func sortXlsx(files []string) []string {
	j := -1
	for i, filename := range files {
		if filepath.Base(filename) == "define.xlsx" {
			j++
			files[i] = files[j]
			files[j] = filename
			break
		}
	}
	return files
}

func parseXlsx(v domain.IParser, filename string) error {
	// 加载文件
	fb, err := excelize.OpenFile(filename)
	if err != nil {
		return uerror.NewUError(1, -1, "filename: %v, error: %v", filename, err)
	}
	defer fb.Close()
	// 解析内容
	isDefine := filepath.Base(filename) == "define.xlsx"
	for _, sheet := range sortSheet(fb.GetSheetList()) {
		values, err := fb.GetRows(sheet)
		if err != nil {
			return uerror.NewUError(1, -1, "filename: %s, error: %v", filename, err)
		}
		// 解析
		if strings.HasPrefix(sheet, "@") {
			if v.Visit(typespec.NewTableNode(sheet, values[0], values[1])) == nil {
				continue
			}
			v.Visit(&typespec.ValueNode{Sheet: sheet, Values: values[2:]})
		} else if isDefine || sheet == "代对表" {
			for _, vals := range values {
				for _, val := range vals {
					ss := strings.Split(val, ":")
					switch len(ss) {
					case 4:
						v.Visit(typespec.NewEnumNode(sheet, ss))
					case 3:
						v.Visit(typespec.NewProxyNode(sheet, ss))
					}
				}
			}
		}
	}
	return nil
}

func ParseXlsxs(v domain.IParser, files ...string) error {
	for _, filename := range sortXlsx(files) {
		v.SetFile(filename)

		if err := parseXlsx(v, filename); err != nil {
			return err
		}
	}
	return nil
}

// 保存文件
func SaveJson(filename string, data []map[string]interface{}) error {
	// 创建目录
	if err := os.MkdirAll(filepath.Dir(filename), os.FileMode(0777)); err != nil {
		return uerror.NewUError(1, -1, "filename: %s, error: %v", filename, err)
	}

	// 写入文件
	buf, _ := json.MarshalIndent(&data, "", "	")
	if err := ioutil.WriteFile(filename, buf, os.FileMode(0666)); err != nil {
		return uerror.NewUError(1, -1, "filename: %s, error: %v", filename, err)
	}
	return nil
}
