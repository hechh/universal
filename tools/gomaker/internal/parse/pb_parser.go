package parse

import (
	"strings"
	"universal/framework/uerror"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/typespec"
	"universal/tools/gomaker/internal/util"

	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
)

type PbParser struct{}

func (d *PbParser) ParseEnum(filename string) error {
	fb, err := excelize.OpenFile(filename)
	if err != nil {
		return uerror.NewUError(1, -1, "%v", err)
	}
	defer fb.Close()
	for _, sheet := range fb.GetSheetList() {
		values, _ := fb.GetRows(sheet)
		switch sheet {
		case "类型配置表":
			for _, vals := range values[2:] {
				if len(vals) > 1 {
					manager.AddConv(strings.TrimSpace(vals[0]), strings.TrimSpace(vals[1]), nil)
				}
			}
		default:
			for _, vals := range values {
				for _, val := range vals {
					if !strings.HasPrefix(val, "E:") && !strings.HasPrefix(val, "e:") {
						continue
					}
					ss := strings.Split(val, ":")
					tt := manager.GetTypeReference(&typespec.Type{
						Kind:    domain.KIND_ENUM,
						PkgName: "pb",
						Name:    ss[2],
						Doc:     ss[1],
					})
					manager.AddConv(tt.GetPkgType(), ss[2], manager.DefaultEnumConv)
					manager.GetOrNewEnum(tt).Add(tt, ss[3], cast.ToInt32(ss[4]), ss[1])
				}
			}
		}
	}
	return nil
}

func (d *PbParser) ParseConfig(filename string) error {
	fb, err := excelize.OpenFile(filename)
	if err != nil {
		return uerror.NewUError(1, -1, "%v", err)
	}
	defer fb.Close()
	// 根据生成表解析config结构
	for k, v := range GetTables(fb) {
		values, _ := fb.GetRows(k)
		item := ParseXlsxStruct(v, values[0], values[1])
		if len(item.List) <= 0 {
			continue
		}
		util.Panic(manager.AddStruct(item))
	}
	return nil
}

func ParseXlsxStruct(tableName string, val01, val02 []string) *typespec.Struct {
	item := &typespec.Struct{
		Type: manager.GetTypeReference(&typespec.Type{
			Kind:    domain.KIND_STRUCT,
			PkgName: "pb",
			Name:    tableName,
		}),
		Fields: make(map[string]*typespec.Field),
	}
	// 解析结构类型
	for i, val := range val01 {
		if len(val02) <= i {
			val02 = append(val02, "")
		}
		val = strings.TrimSpace(val)
		if len(val) <= 0 {
			continue
		}
		field := parseXlsxField(i, val, val02[i])
		item.Add(field)
	}
	return item
}

func GetTables(fb *excelize.File) map[string]string {
	tmps := map[string]string{}
	values, _ := fb.GetRows("生成表")
	for _, vals := range values {
		for _, val := range vals {
			if pos := strings.Index(val, ":"); pos > 0 {
				tmps[strings.TrimSpace(val[:pos])] = strings.TrimSpace(val[pos+1:])
			}
		}
	}
	return tmps
}

func parseXlsxField(i int, str, doc string) *typespec.Field {
	pos := strings.Index(str, "/")
	if pos <= 0 {
		return &typespec.Field{
			Type: manager.GetTypeReference(&typespec.Type{
				Kind: domain.KIND_IDENT,
				Name: "string",
			}),
			Name:  str,
			Doc:   doc,
			Index: i,
		}
	}
	dot := strings.Index(str, ".")
	if dot <= 0 {
		return &typespec.Field{
			Type: manager.GetTypeReference(&typespec.Type{
				Kind: domain.KIND_IDENT,
				Name: str[pos+1:],
			}),
			Name:  str[:pos],
			Doc:   doc,
			Index: i,
		}
	}
	return &typespec.Field{
		Type: manager.GetTypeReference(&typespec.Type{
			Kind:    domain.KIND_IDENT,
			Name:    str[dot+1:],
			PkgName: str[pos+1 : dot],
		}),
		Name:  str[:pos],
		Doc:   doc,
		Index: i,
	}
}
