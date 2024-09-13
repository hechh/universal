package generator

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"text/template"
	"unicode"
	"universal/framework/uerror"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/parse"
	"universal/tools/gomaker/internal/typespec"
	"universal/tools/gomaker/internal/util"

	_ "universal/common/pb"

	"github.com/golang/protobuf/proto"
	"github.com/xuri/excelize/v2"
)

// 生成enum.gen.proto
func bytesGenerator(dst string, tpls *template.Template, extra ...string) error {
	for _, filename := range extra {
		if err := parseXlsx(filename, dst); err != nil {
			return err
		}
	}
	return nil
}

func parseXlsx(filename string, dst string) error {
	fb, err := excelize.OpenFile(filename)
	if err != nil {
		return uerror.NewUError(1, -1, "%v", err)
	}
	defer fb.Close()
	// 根据生成表解析config结构
	for k, v := range parse.GetTables(fb) {
		values, _ := fb.GetRows(k)
		item := parse.ParseXlsxStruct(v, values[0], values[1])
		if len(item.List) <= 0 {
			continue
		}
		// 加载pb.go文件中的结构信息
		pbItem := manager.GetStruct(item.Type.GetPkgType())
		if pbItem == nil {
			return uerror.NewUError(1, -1, "配置表中未定义%s结构", v)
		}
		fields := map[int]*typespec.Field{}
		for _, ff := range pbItem.List {
			if unicode.IsLower(rune(ff.Name[0])) {
				continue
			}
			ff.Index = item.Fields[ff.Name].Index
			fields[ff.Index] = ff
			if ff.Type.Kind == domain.KIND_ENUM {
				manager.AddConv(ff.Type.GetPkgType(), ff.Type.Name, manager.DefaultEnumConv)
			}
		}
		// 解析配置数据
		if err := toBytes(pbItem, fields, values[2:], dst); err != nil {
			return err
		}
	}
	return nil
}

func toBytes(pbItem *typespec.Struct, fields map[int]*typespec.Field, values [][]string, dst string) error {
	tmps := []interface{}{}
	for _, vals := range values {
		tmp := map[string]interface{}{}
		for pos, field := range fields {
			if len(vals) > pos && len(vals[pos]) > 0 {
				tmp[field.Name] = manager.ToConvert(field.Type.GetPkgType(), vals[pos])
			}
		}
		if len(tmp) > 0 {
			tmps = append(tmps, tmp)
		}
	}
	if len(tmps) > 0 {
		data := reflect.New(proto.MessageType(pbItem.Type.GetPkgType() + "Ary").Elem()).Interface()
		buf, _ := json.Marshal(map[string]interface{}{"Ary": tmps})
		//util.SaveJson(filepath.Join(dst, fmt.Sprintf("%s.json", pbItem.Type.Name)), buf)
		json.Unmarshal(buf, data)
		return util.SaveBytes(filepath.Join(dst, fmt.Sprintf("%s.bytes", pbItem.Type.Name)), data.(proto.Message))
	}
	return nil
}
