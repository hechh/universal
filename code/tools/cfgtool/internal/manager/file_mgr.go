package manager

import (
	"fmt"
	"path/filepath"
	"strings"
	"universal/tools/cfgtool/domain"
	"universal/tools/cfgtool/internal/util"

	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
)

var (
	alls  = make(map[string]*domain.Enum)     // 中文类型--enum
	files = make(map[string]*domain.FileType) // 文件名 -- fileType
)

// 专门解析define.xlsx文件
func ParseDefine(filename string) error {
	fb, err := excelize.OpenFile(filename)
	if err != nil {
		return err
	}
	defer fb.Close()

	// 创建FileInfo
	info := util.NewFileInfo(filename, alls)
	for _, sheet := range fb.GetSheetList() {
		values, err := fb.GetRows(sheet)
		if err != nil {
			return err
		}
		parseProxy(info, values)
	}
	files[filepath.Base(filename)] = info
	return nil
}

// 解析代对表
func parseProxy(info *domain.FileType, values [][]string) {
	for _, vals := range values {
		for _, val := range vals {
			if !strings.Contains(val, ":") {
				continue
			}

			ss := strings.Split(strings.TrimSpace(val), ":")
			switch ss[0] {
			case "C", "c":
				info.Tables[ss[1]] = &domain.Table{IsServer: false, Name: ss[2]}
			case "CS", "cs", "Cs", "cS", "SC", "Sc", "sC", "sc", "s", "S":
				info.Tables[ss[1]] = &domain.Table{IsServer: true, Name: ss[2]}
			case "E", "e":
				item := &domain.Enum{
					Type:  ss[2],
					Name:  ss[3],
					Value: cast.ToInt32(ss[4]),
					Doc:   ss[1],
				}
				// 保存
				info.Alls[item.Doc] = item
				if _, ok := info.Enums[ss[2]]; !ok {
					info.Enums[ss[2]] = []*domain.Enum{}
				}
				info.Enums[ss[2]] = append(info.Enums[ss[2]], item)
			}
		}
	}
}

type FPointer struct {
	filename string
	fb       *excelize.File
	info     *domain.FileType
}

func ParseProxy(filename string) (*FPointer, error) {
	fb, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}

	// 获取所有数据
	values, err := fb.GetRows("代对表")
	if err != nil {
		return nil, err
	}

	// 解析所有代对
	info := util.NewFileInfo(filename, alls)
	parseProxy(info, values)
	files[filepath.Base(filename)] = info
	return &FPointer{filename: filename, fb: fb, info: info}, nil
}

func (d *FPointer) ParseTable(dst string) error {
	defer d.fb.Close()
	for _, sheet := range d.fb.GetSheetList() {
		if !strings.HasPrefix(sheet, "@") {
			continue
		}
		// 判断是否需要解析数据
		tableName := strings.TrimPrefix(sheet, "@")
		vv, ok := d.info.Tables[tableName]
		if !ok || !vv.IsServer {
			continue
		}
		// 加载所有数据
		values, err := d.fb.GetRows(sheet)
		if err != nil {
			return err
		}
		// 解析字段数据
		vv.Fields = parseField(values[0], values[1])

		// 解析json数据
		data := parseJson(alls, vv.Fields, values[2:])

		// 保存json文件
		jsonName := filepath.Join(dst, fmt.Sprintf("%s.%s.json", d.info.Name, vv.Name))
		if err := util.SaveJson(jsonName, data); err != nil {
			return err
		}
	}
	return nil
}

func parseField(val01, val02 []string) (rets []*domain.Field) {
	docF := func(i int) string {
		if len(val02) > i {
			return val02[i]
		}
		return ""
	}
	for i, name := range val01 {
		if !strings.Contains(name, "_") {
			continue
		}

		index := strings.Index(name, "_")
		rets = append(rets, &domain.Field{
			IsEnum:   strings.HasPrefix(name, "$"),
			Original: strings.TrimPrefix(name, "$"),
			Index:    i,
			Type:     strings.ToLower(strings.TrimPrefix(name[:index], "$")),
			Name:     name[index+1:],
			Doc:      docF(i),
		})
	}
	return
}

func parseJson(alls map[string]*domain.Enum, fields []*domain.Field, values [][]string) (rets []map[string]interface{}) {
	for _, vals := range values {
		js := map[string]interface{}{}
		for _, field := range fields {
			if field.Index < len(vals) {
				value := strings.TrimSpace(vals[field.Index])
				if len(value) > 0 {
					js[field.Name] = parseValue(alls, field, value)
				}
			}
		}
		rets = append(rets, js)
	}
	return
}

// 类型转换
func parseValue(alls map[string]*domain.Enum, ff *domain.Field, value string) interface{} {
	value = strings.ReplaceAll(value, "|", ",")
	proxyF := func(str string) interface{} {
		if vv, ok := alls[str]; ok {
			return vv.Value
		}
		return str
	}

	switch ff.Type {
	case "i":
		return cast.ToUint64(proxyF(value))
	case "il":
		rets := []interface{}{}
		for _, str := range strings.Split(value, ",") {
			rets = append(rets, cast.ToUint64(proxyF(str)))
		}
		return rets
	case "ill":
		rets := []interface{}{}
		for _, ss := range strings.Split(value, "#") {
			subs := []interface{}{}
			for _, str := range strings.Split(ss, ",") {
				subs = append(subs, cast.ToUint64(proxyF(str)))
			}
			rets = append(rets, subs)
		}
		return rets
	case "f":
		return cast.ToFloat64(proxyF(value))
	case "fl":
		rets := []interface{}{}
		for _, str := range strings.Split(value, ",") {
			rets = append(rets, cast.ToFloat64(proxyF(str)))
		}
		return rets
	case "fll":
		rets := []interface{}{}
		for _, ss := range strings.Split(value, "#") {
			subs := []interface{}{}
			for _, str := range strings.Split(ss, ",") {
				subs = append(subs, cast.ToFloat64(proxyF(str)))
			}
			rets = append(rets, subs)
		}
		return rets
	case "b":
		return cast.ToBool(proxyF(value))
	}
	return value
}
