package base

import (
	"bytes"
	"fmt"
	"hego/tools/xlsx/domain"
	"sort"
	"strings"
)

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
		strs = append(strs, fmt.Sprintf("%s %s // %s", val.Name, val.Type.GetType(""), val.Desc))
	}
	buf.WriteString(fmt.Sprintf("type %s struct {\n%s\n}\n", d.Name, strings.Join(strs, "\n")))
}

func (d *Config) Format(buf *bytes.Buffer) {
	strs := []string{}
	for _, val := range d.List {
		strs = append(strs, fmt.Sprintf("%s %s // %s", val.Name, val.Type.GetType(""), val.Desc))
	}
	buf.WriteString(fmt.Sprintf("type %sConfig struct {\n%s\n}\n", d.Name, strings.Join(strs, "\n")))
}

// 获取类型字符串
func (d *Type) GetType(pkgName string) string {
	switch d.TypeOf {
	case domain.TYPE_OF_BASE:
		switch d.ValueOf {
		case domain.VALUE_OF_IDENT:
			return d.Name
		case domain.VALUE_OF_ARRAY:
			return fmt.Sprintf("[]%s", d.Name)
		}
	case domain.TYPE_OF_ENUM:
		switch d.ValueOf {
		case domain.VALUE_OF_IDENT:
			if len(pkgName) > 0 {
				return fmt.Sprintf("%s.%s", pkgName, d.Name)
			}
			return d.Name
		case domain.VALUE_OF_ARRAY:
			if len(pkgName) > 0 {
				return fmt.Sprintf("[]%s.%s", pkgName, d.Name)
			}
			return fmt.Sprintf("[]%s", d.Name)
		}
	case domain.TYPE_OF_STRUCT, domain.TYPE_OF_CONFIG:
		switch d.ValueOf {
		case domain.VALUE_OF_IDENT:
			if len(pkgName) > 0 {
				return fmt.Sprintf("*%s.%s", pkgName, d.Name)
			}
			return fmt.Sprintf("*%s", d.Name)
		case domain.VALUE_OF_ARRAY:
			if len(pkgName) > 0 {
				return fmt.Sprintf("[]*%s.%s", pkgName, d.Name)
			}
			return fmt.Sprintf("[]*%s", d.Name)
		}
	}
	return ""
}
