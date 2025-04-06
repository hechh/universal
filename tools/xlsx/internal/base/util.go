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
	case domain.TypeOfBase:
		switch d.ValueOf {
		case domain.ValueOfBase:
			return d.Name
		case domain.ValueOfList:
			return fmt.Sprintf("[]%s", d.Name)
		}
	case domain.TypeOfEnum:
		switch d.ValueOf {
		case domain.ValueOfBase:
			if len(pkgName) > 0 {
				return fmt.Sprintf("%s.%s", pkgName, d.Name)
			}
			return d.Name
		case domain.ValueOfList:
			if len(pkgName) > 0 {
				return fmt.Sprintf("[]%s.%s", pkgName, d.Name)
			}
			return fmt.Sprintf("[]%s", d.Name)
		}
	case domain.TypeOfStruct, domain.TypeOfConfig:
		switch d.ValueOf {
		case domain.ValueOfBase:
			if len(pkgName) > 0 {
				return fmt.Sprintf("*%s.%s", pkgName, d.Name)
			}
			return fmt.Sprintf("*%s", d.Name)
		case domain.ValueOfList:
			if len(pkgName) > 0 {
				return fmt.Sprintf("[]*%s.%s", pkgName, d.Name)
			}
			return fmt.Sprintf("[]*%s", d.Name)
		}
	}
	return ""
}

func (d *Index) Value(ref, split string) string {
	strs := []string{}
	for _, field := range d.List {
		if len(ref) > 0 {
			strs = append(strs, fmt.Sprintf("%s.%s", ref, field.Name))
		} else {
			strs = append(strs, field.Name)
		}
	}
	return strings.Join(strs, split)
}

func (d *Index) Arg(pkg, split string) string {
	strs := []string{}
	for _, field := range d.List {
		strs = append(strs, fmt.Sprintf("%s %s", field.Name, field.Type.GetType(pkg)))
	}
	return strings.Join(strs, split)
}

func (d *Index) IndexValue(val string) string {
	switch d.Type.TypeOf {
	case domain.TypeOfBase:
		return val
	case domain.TypeOfStruct:
		return fmt.Sprintf("%s{%s}", d.Name, val)
	}
	return ""
}

func (d *Index) GetType(cfg string) string {
	switch d.Type.ValueOf {
	case domain.ValueOfMap:
		return fmt.Sprintf("map[%s]*%s", d.Type.Name, cfg)
	case domain.ValueOfGroup:
		return fmt.Sprintf("map[%s][]*%s", d.Type.Name, cfg)
	}
	return fmt.Sprintf("[]*%s", cfg)
}
