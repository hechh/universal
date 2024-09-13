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
					manager.AddConvType(strings.TrimSpace(vals[0]), strings.TrimSpace(vals[1]))
				}
			}
		default:
			for _, vals := range values {
				for _, val := range vals {
					if !strings.HasPrefix(val, "E:") {
						continue
					}
					ss := strings.Split(val, ":")
					tt := manager.GetTypeReference(&typespec.Type{
						Kind:    domain.KIND_ENUM,
						PkgName: "pb",
						Name:    ss[2],
						Doc:     ss[1],
					})
					manager.AddConvType(tt.GetPkgType(), ss[2])
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
	// 解析生成表
	tmps := map[string]string{}
	values, _ := fb.GetRows("生成表")
	for _, vals := range values {
		for _, val := range vals {
			if pos := strings.Index(val, ":"); pos > 0 {
				tmps[strings.TrimSpace(val[:pos])] = strings.TrimSpace(val[pos+1:])
			}
		}
	}
	// 根据生成表解析config结构
	for k, v := range tmps {
		item := &typespec.Struct{
			Type: manager.GetTypeReference(&typespec.Type{
				Kind:    domain.KIND_STRUCT,
				PkgName: "pb",
				Name:    v,
			}),
			Fields: make(map[string]*typespec.Field),
		}
		values, _ := fb.GetRows(k)
		for i, val := range values[0] {
			if len(values[1]) <= i {
				values[1] = append(values[1], "")
			}
			val = strings.TrimSpace(val)
			if len(val) <= 0 {
				continue
			}
			field := d.parseField(val, values[1][i])
			item.Add(field)
		}
		if len(item.List) > 0 {
			util.Panic(manager.AddStruct(item))
		}
	}
	return nil
}

func (d *PbParser) parseField(str, doc string) *typespec.Field {
	pos := strings.Index(str, "/")
	if pos <= 0 {
		return &typespec.Field{
			Type: manager.GetTypeReference(&typespec.Type{
				Kind: domain.KIND_IDENT,
				Name: "string",
			}),
			Name: str,
			Doc:  doc,
		}
	}
	dot := strings.Index(str, ".")
	if dot <= 0 {
		return &typespec.Field{
			Type: manager.GetTypeReference(&typespec.Type{
				Kind: domain.KIND_IDENT,
				Name: str[pos+1:],
			}),
			Name: str[:pos],
			Doc:  doc,
		}
	}
	return &typespec.Field{
		Type: manager.GetTypeReference(&typespec.Type{
			Kind:    domain.KIND_IDENT,
			Name:    str[dot+1:],
			PkgName: str[pos+1 : dot],
		}),
		Name: str[:pos],
		Doc:  doc,
	}
}
