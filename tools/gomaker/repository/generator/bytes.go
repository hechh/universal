package generator

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
	"unicode"
	"universal/framework/uerror"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/parser"
	"universal/tools/gomaker/internal/typespec"
	"universal/tools/gomaker/internal/util"

	_ "universal/common/pb"

	"github.com/golang/protobuf/proto"
	"github.com/xuri/excelize/v2"
)

func BytesGen(dst string, tpls *template.Template, extra ...string) error {
	for _, filename := range extra {
		if err := parseXlsx(filename, dst); err != nil {
			return err
		}
	}
	return nil
}

func parseXlsx(filename string, dst string) error {
	// 打开文件
	fp, err := excelize.OpenFile(filename)
	if err != nil {
		return uerror.NewUError(1, -1, "开打%s失败：%v", filename, err)
	}
	defer fp.Close()
	// 读取生成表
	values, err := fp.GetRows(domain.GenTable)
	if _, ok := err.(excelize.ErrSheetNotExist); ok || len(values) <= 0 {
		return nil
	}
	if err != nil {
		return uerror.NewUError(1, -1, "读取%s配置表%s失败: %v", filename, domain.GenTable, err)
	}
	// 解析生成表
	for _, vals := range values {
		for _, val := range vals {
			if !strings.HasPrefix(val, domain.RuleTypeBytes) {
				continue
			}
			// 解析@gomaker
			if err := parseSheet(parser.ParseXlsxSheet(val, fp), dst); err != nil {
				return err
			}
		}
	}
	return nil
}

func parseSheet(sh *typespec.Sheet, dst string) error {
	values, err := sh.GetRows()
	if err != nil {
		return uerror.NewUError(1, -1, "读取配置表%s失败: %v", sh.Sheet, err)
	}
	// 解析结构
	st := parser.ParseXlsxStruct(sh, values[0], values[1])
	if st == nil {
		return nil
	}
	// 加载从pb.go文件解析的结构信息
	pbItem := manager.GetStruct(typespec.GetPkgType(st.Type))
	if pbItem == nil {
		return uerror.NewUError(1, -1, "%s配置表生成的结构被删除", sh.Sheet)
	}
	fields := map[int]*typespec.Field{}
	for _, ff := range pbItem.List {
		if unicode.IsLower(rune(ff.Name[0])) {
			continue
		}
		ff.Index = st.Fields[ff.Name].Index
		fields[ff.Index] = ff
	}
	// 解析配置数据
	if err := toBytes(pbItem, fields, values[2:], dst); err != nil {
		return err
	}
	return nil
}

func toBytes(pbItem *typespec.Struct, fields map[int]*typespec.Field, values [][]string, dst string) error {
	tmps := []interface{}{}
	for _, vals := range values {
		tmp := map[string]interface{}{}
		for pos, field := range fields {
			if len(vals) > pos && len(vals[pos]) > 0 {
				tmp[field.Name] = manager.Cast(field, vals[pos])
			}
		}
		if len(tmp) > 0 {
			tmps = append(tmps, tmp)
		}
	}
	if len(tmps) > 0 {
		filename := filepath.Join(dst, fmt.Sprintf("%s.bytes", pbItem.Type.Name))
		return util.SaveFile(filename, jsonToProto(pbItem.Type, tmps))
	}
	return nil
}

func jsonToProto(item *typespec.Type, data interface{}) []byte {
	ary := fmt.Sprintf("%sAry", typespec.GetPkgType(item))
	buf, _ := json.Marshal(map[string]interface{}{"Ary": data})
	pbData := reflect.New(proto.MessageType(ary).Elem()).Interface()
	json.Unmarshal(buf, pbData)
	result, _ := proto.Marshal(pbData.(proto.Message))
	return result
}
