package xlsx

/*
import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"text/template"
	"unicode"
	"hego/framework/uerror"
	"hego/tools/gomaker/internal/manager"
	"hego/tools/gomaker/internal/parser"
	"hego/tools/gomaker/internal/typespec"
	"hego/tools/gomaker/internal/util"

	_ "hego/common/pb"

	"github.com/golang/protobuf/proto"
)

func BytesGen(dst string, tpls *template.Template, extra ...string) error {
	for _, filename := range extra {
		// 解析生成表
		cfgs := []*typespec.Sheet{}
		if err := parser.ParseGenTable(filename, nil, &cfgs); err != nil {
			return err
		}
		// 对游戏配置生成bytes文件
		for _, sh := range cfgs {
			if err := parseSheet(sh, dst); err != nil {
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
	pbItem := manager.GetStruct(st.Type.GetPkgType())
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
		jsbuf, bytes := jsonToProto(pbItem.Type, tmps)
		util.SaveFile(filepath.Join(dst, fmt.Sprintf("%s.json", pbItem.Type.Name)), jsbuf)
		return util.SaveFile(filename, bytes)
	}
	return nil
}

func jsonToProto(item *typespec.Type, data interface{}) (jsbuf, bytes []byte) {
	ary := fmt.Sprintf("%sAry", item.GetPkgType())
	jsbuf, _ = json.Marshal(map[string]interface{}{"Ary": data})
	pbData := reflect.New(proto.MessageType(ary).Elem()).Interface()
	json.Unmarshal(jsbuf, pbData)
	bytes, _ = proto.Marshal(pbData.(proto.Message))
	return
}
*/
