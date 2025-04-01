package base

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"universal/framework/uerror"
	"universal/tools/xlsx/domain"
)

func SaveGo(filename string, buf *bytes.Buffer) error {
	result, err := format.Source(buf.Bytes())
	if err != nil {
		return uerror.NewUError(1, -1, "格式化失败: %v", err)
	}
	return SaveFile(filename, result)
}

func SaveFile(filename string, buf []byte) error {
	// 创建目录
	if err := os.MkdirAll(filepath.Dir(filename), os.FileMode(0777)); err != nil {
		return uerror.NewUError(1, -1, "filename: %s, error: %v", filename, err)
	}

	// 写入文件
	if err := ioutil.WriteFile(filename, buf, os.FileMode(0666)); err != nil {
		return uerror.NewUError(1, -1, "filename: %s, error: %v", filename, err)
	}
	return nil
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
