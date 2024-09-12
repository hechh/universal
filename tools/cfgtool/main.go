package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"universal/framework/basic/uerror"
	"universal/framework/basic/util"
	"universal/tools/cfgtool/internal/typespec"

	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
)

func panicout(err error) {
	if uerr, ok := err.(*uerror.UError); ok {
		panic(uerr.String())
	} else {
		panic(err)
	}
}

func main() {
	var src, dst string
	flag.StringVar(&src, "src", "./", "xlsx配置文件目录")
	flag.StringVar(&dst, "dst", "./", "proto生成目录")
	flag.Parse()
	// 获取当前路径
	cwd, err := os.Getwd()
	if err != nil {
		panic(uerror.NewUError(1, -1, "%v", err).String())
	}
	// 设置绝对路径
	src = filepath.Clean(filepath.Join(cwd, src))
	dst = filepath.Clean(filepath.Join(cwd, dst))
	// 优先解析enum.xlsx文件
	enums := make(map[string]*typespec.Enum)
	types := make(map[string][]*typespec.Enum)
	converts := make(map[string]string)
	filename := filepath.Join(src, "enum.xlsx")
	if err := parseEnum(filename, enums, types, converts); err != nil {
		panicout(err)
	}
	// 生成enum.gen.proto
	if err := genEnum(dst, types); err != nil {
		panicout(err)
	}
	// 解析table表
	files, err := util.Glob(src, ".*xlsx", "enum.xlsx", true)
	if err != nil {
		panicout(err)
	}
	sts := make(map[string][]*typespec.Field)
	for _, filename := range files {
		if err := parseTable(filename, types, converts, sts); err != nil {
			panicout(err)
		}
	}
	// 生成table.gen.proto
	if err := genTable(dst, converts, sts); err != nil {
		panicout(err)
	}
}

func parseTable(filename string, types map[string][]*typespec.Enum, convs map[string]string, sts map[string][]*typespec.Field) error {
	fb, err := excelize.OpenFile(filename)
	if err != nil {
		return err
	}
	defer fb.Close()
	// 解析table表
	tmps := map[string]string{}
	if tables, err := fb.GetRows("生成表"); err != nil {
		return uerror.NewUError(1, -1, "%v", err)
	} else {
		for _, vals := range tables {
			for _, val := range vals {
				val = strings.TrimSpace(val)
				if len(val) <= 0 {
					continue
				}
				if strings.Contains(val, ":") {
					ss := strings.Split(val, ":")
					tmps[ss[0]] = ss[1]
				}
			}
		}
	}
	// 解析proto结构
	for k, v := range tmps {
		if values, err := fb.GetRows(k); err == nil {
			sts[v] = typespec.ParseField(values[0], values[1])
		}
	}
	return nil
}

func parseEnum(filename string, enums map[string]*typespec.Enum, types map[string][]*typespec.Enum, convs map[string]string) error {
	fb, err := excelize.OpenFile(filename)
	if err != nil {
		return err
	}
	defer fb.Close()
	// 解析枚举类型
	for _, sheet := range fb.GetSheetList() {
		values, err := fb.GetRows(sheet)
		if err != nil {
			continue
		}
		if sheet == "类型配置表" {
			for _, row := range values[2:] {
				convs[strings.TrimSpace(row[0])] = strings.TrimSpace(row[1])
			}
		} else {
			for _, row := range values {
				for _, val := range row {
					ss := strings.Split(val, ":")
					eval := &typespec.Enum{
						ID:    ss[1],
						Type:  ss[2],
						Name:  fmt.Sprintf("%s_%s", ss[2], ss[3]),
						Value: cast.ToInt32(ss[4]),
					}
					enums[eval.ID] = eval
					if _, ok := types[eval.Type]; !ok {
						types[eval.Type] = []*typespec.Enum{}
						convs[eval.Type] = eval.Type
					}
					types[eval.Type] = append(types[eval.Type], eval)
				}
			}
		}
	}
	return nil
}

const (
	head = `syntax = "proto3";
package pb;
option go_package = "../../common/pb";`
)

func genEnum(dst string, types map[string][]*typespec.Enum) error {
	// 排序
	list := []string{}
	for key, vals := range types {
		sort.Slice(vals, func(i, j int) bool {
			return vals[i].Value < vals[j].Value
		})
		list = append(list, key)
	}
	sort.Slice(list, func(i, j int) bool {
		return strings.Compare(list[i], list[j]) <= 0
	})
	// 生成
	buf := bytes.NewBufferString(head)
	for _, key := range list {
		strs := []string{}
		for _, val := range types[key] {
			strs = append(strs, fmt.Sprintf("\t%s  = %d; // %s", val.Name, val.Value, val.ID))
		}
		buf.WriteString(fmt.Sprintf("\n\nenum %s {\n %s \n};", key, strings.Join(strs, "\n")))
	}
	// 创建目录
	if err := os.MkdirAll(dst, os.FileMode(0777)); err != nil {
		return uerror.NewUError(1, -1, "%v", err)
	}
	// 写入文件
	if err := ioutil.WriteFile(filepath.Join(dst, "enum.gen.proto"), buf.Bytes(), os.FileMode(0666)); err != nil {
		return uerror.NewUError(1, -1, "%v", err)
	}
	return nil
}

func genTable(dst string, convs map[string]string, sts map[string][]*typespec.Field) error {
	// 排序
	list := []string{}
	for key := range sts {
		list = append(list, key)
	}
	sort.Slice(list, func(i, j int) bool {
		return strings.Compare(list[i], list[j]) <= 0
	})
	// 生成
	buf := bytes.NewBufferString(`syntax = "proto3";
package pb;
import "enum.gen.proto";
option go_package = "../../common/pb";`)
	for _, key := range list {
		strs := []string{}
		for i, val := range sts[key] {
			strs = append(strs, fmt.Sprintf("\t%s %s  = %d; // %s", convs[val.Type], val.Name, i+1, val.Doc))
		}
		buf.WriteString(fmt.Sprintf("\n\nmessage %s {\n %s \n};", key, strings.Join(strs, "\n")))
	}
	// 创建目录
	if err := os.MkdirAll(dst, os.FileMode(0777)); err != nil {
		return uerror.NewUError(1, -1, "%v", err)
	}
	// 写入文件
	if err := ioutil.WriteFile(filepath.Join(dst, "table.gen.proto"), buf.Bytes(), os.FileMode(0666)); err != nil {
		return uerror.NewUError(1, -1, "%v", err)
	}

	return nil
}
