package typespec

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
)

type Enum struct {
	Type  string // 枚举类型
	Name  string // 枚举名字
	Value int32  // 枚举值
	Doc   string // 枚举注释
}

type Table struct {
	IsCreate    bool   // 是否生成配置代码
	EnglishName string // 中文名称
	ChineseName string // 英文名称
}

type Field struct {
	IsProxy bool   // 是否取代对值
	Name    string // 字段名字
	Doc     string // 注释
}

type Struct struct {
	Fields []*Field                 // 字段数据
	Jsons  []map[string]interface{} // json数据
}

type FileInfo struct {
	Name   string              // 文件名
	fb     *excelize.File      // xlsx文件句柄
	Enums  map[string][]*Enum  // 枚举类型
	Alls   map[string]*Enum    // 所有枚举类型
	Proxys map[string]*Table   // 代对表
	Tables map[string][]*Field // 配置表
}

func NewFileInfo(filename string, alls map[string]*Enum) (*FileInfo, error) {
	// 打开文件
	fb, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	return &FileInfo{
		Name:   filepath.Base(filename),
		fb:     fb,
		Alls:   alls,
		Enums:  make(map[string][]*Enum),
		Proxys: make(map[string]*Table),
		Tables: make(map[string][]*Field),
	}, nil
}

func (d *FileInfo) GetFPointer() *excelize.File {
	return d.fb
}

func (d *FileInfo) Close() {
	d.fb.Close()
}

func (d *Field) GetValue(str string, alls map[string]*Enum) interface{} {
	valuef := func(val string) interface{} {
		if vv, ok := alls[val]; ok {
			return vv.Value
		}
		return val
	}
	str = strings.ReplaceAll(str, "|", ",")
	switch strings.Split(d.Name, "_")[0] {
	case "i":
		return cast.ToUint64(valuef(str))
	case "il":
		rets := []interface{}{}
		for _, ss := range strings.Split(str, ",") {
			rets = append(rets, cast.ToUint64(valuef(ss)))
		}
		return rets
	case "ill":
		rets := []interface{}{}
		for _, substr := range strings.Split(str, "#") {
			subs := []interface{}{}
			for _, ss := range strings.Split(substr, ",") {
				subs = append(subs, cast.ToUint64(valuef(ss)))
			}
			rets = append(rets, subs...)
		}
		return rets
	case "b":
		return cast.ToBool(valuef(str))
	case "f":
		return cast.ToFloat64(valuef(str))
	case "fl":
		rets := []interface{}{}
		for _, ss := range strings.Split(str, ",") {
			rets = append(rets, cast.ToFloat64(valuef(ss)))
		}
		return rets
	}
	return str
}

func (d *FileInfo) GetJsonName(dst, name string) string {
	return filepath.Join(dst, fmt.Sprintf("%s.%s.json", strings.TrimSuffix(d.Name, ".xlsx"), name))
}

func (d *FileInfo) GetProxy(sheet string) *Table {
	return d.Proxys[sheet]
}

func (d *FileInfo) ParseTable(sheetName string) ([]map[string]interface{}, error) {
	values, err := d.fb.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	// 第一行一定是字段, 第二行一定值注释
	fields := make([]*Field, len(values[0]))
	d.Tables[strings.TrimPrefix(sheetName, "@")] = fields
	for i, name := range values[0] {
		doc := ""
		if len(values[1]) > i {
			doc = values[1][i]
		}
		fields[i] = &Field{
			IsProxy: strings.HasPrefix(name, "$"),
			Name:    strings.TrimPrefix(name, "$"),
			Doc:     doc,
		}
	}

	// 解析json数据
	jsons := []map[string]interface{}{}
	for _, vals := range values[2:] {
		js := map[string]interface{}{}
		for i, val := range vals {
			if len(fields) > i {
				js[fields[i].Name] = fields[i].GetValue(strings.TrimSpace(val), d.Alls)
			}
		}
		jsons = append(jsons, js)
	}
	return jsons, nil
}

// 解析代对表
func (d *FileInfo) ParseProxy(sheetName string) error {
	values, err := d.fb.GetRows(sheetName)
	if err != nil {
		return err
	}
	for _, vals := range values {
		for _, val := range vals {
			if !strings.Contains(val, ":") {
				continue
			}
			ss := strings.Split(strings.TrimSpace(val), ":")
			switch ss[0] {
			case "C", "c":
				d.Proxys[ss[1]] = &Table{IsCreate: false, EnglishName: ss[2], ChineseName: ss[1]}
			case "CS", "cs", "Cs", "cS", "SC", "Sc", "sC", "sc":
				d.Proxys[ss[1]] = &Table{IsCreate: true, EnglishName: ss[2], ChineseName: ss[1]}
			case "E", "e":
				item := &Enum{
					Type:  ss[2],
					Name:  ss[3],
					Value: cast.ToInt32(ss[4]),
					Doc:   ss[1],
				}
				// 保存
				d.Alls[item.Doc] = item
				if _, ok := d.Enums[ss[2]]; !ok {
					d.Enums[ss[2]] = []*Enum{}
				}
				d.Enums[ss[2]] = append(d.Enums[ss[2]], item)
			}
		}
	}
	return nil
}
