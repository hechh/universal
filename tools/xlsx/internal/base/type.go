package base

import (
	"bytes"
	"fmt"
	"hego/tools/xlsx/domain"
	"sort"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Type struct {
	Name    string // 类型名称
	TypeOf  uint32 // 1表示内置类型，2表示枚举，3表示结构体
	ValueOf uint32 // 1 表示单值，2:表示数组，3:表示map，4 表示group
}

type Field struct {
	Name     string // 字段名
	Type     *Type  // 字段类型
	Desc     string // 字段标签
	Position int    // 字段索引
}

type EValue struct {
	Name  string // 枚举值名称
	Value uint32 // 枚举值
	Desc  string // 枚举值描述
}

// 枚举类型定义
type Enum struct {
	Name     string             // 枚举名称
	Values   map[string]*EValue // 枚举值
	FileName string             // 文件名
}

type Struct struct {
	Name     string              // 结构体名称
	List     []*Field            // 字段类型
	Converts map[string][]*Field // 转换表
	FileName string              // 文件名
}

// 生成表
type Config struct {
	Name     string   // 表名称
	List     []*Field // 表列表
	FileName string   // 文件名
	Map      []*Field
	Group    []*Field
}

type Table struct {
	TypeOf    uint32
	SheetName string
	TypeName  string
	FileName  string
	Rules     []string
	Fp        *excelize.File
}

type Value struct {
	TypeOf   uint32
	Type     string // 类型名称
	Name     string // 枚举值名称
	Value    uint32 // 枚举值
	Desc     string // 枚举值描述
	FileName string // 文件名
}

func (d *Table) ScanRows(count int) (rets [][]string, err error) {
	rows, err := d.Fp.Rows(d.SheetName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for i := 0; i < count && rows.Next(); i++ {
		row, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		rets = append(rets, row)
	}
	return
}

func (d *Enum) Format(buf *bytes.Buffer) {
	vals := []*EValue{}
	for _, v := range d.Values {
		vals = append(vals, v)
	}
	sort.Slice(vals, func(i, j int) bool {
		return vals[i].Value < vals[j].Value
	})
	strs := []string{}
	for _, val := range vals {
		strs = append(strs, fmt.Sprintf("%s %s = %d // %s", val.Name, d.Name, val.Value, val.Desc))
	}
	buf.WriteString(fmt.Sprintf("type %s uint32\n const(%s\n)\n", d.Name, strings.Join(strs, "\n")))
}

func (d *Struct) Format(buf *bytes.Buffer) {
	strs := []string{}
	for _, val := range d.List {
		strs = append(strs, fmt.Sprintf("%s %s // %s", val.Name, val.Type.String(), val.Desc))
	}
	buf.WriteString(fmt.Sprintf("type %s struct {\n%s\n}\n", d.Name, strings.Join(strs, "\n")))
}

func (d *Config) Format(buf *bytes.Buffer) {
	strs := []string{}
	for _, val := range d.List {
		strs = append(strs, fmt.Sprintf("%s %s // %s", val.Name, val.Type.String(), val.Desc))
	}
	buf.WriteString(fmt.Sprintf("type %sConfig struct {\n%s\n}\n", d.Name, strings.Join(strs, "\n")))
}

func (d *Type) String() string {
	switch d.TypeOf {
	case domain.TYPE_OF_BASE, domain.TYPE_OF_ENUM:
		switch d.ValueOf {
		case domain.VALUE_OF_IDENT:
			return d.Name
		case domain.VALUE_OF_ARRAY:
			return fmt.Sprintf("[]%s", d.Name)
		}
	case domain.TYPE_OF_STRUCT, domain.TYPE_OF_CONFIG:
		switch d.ValueOf {
		case domain.VALUE_OF_IDENT:
			return fmt.Sprintf("*%s", d.Name)
		case domain.VALUE_OF_ARRAY:
			return fmt.Sprintf("[]*%s", d.Name)
		case domain.VALUE_OF_MAP:
		case domain.VALUE_OF_GROUP:
		}
	}
	return ""
}
