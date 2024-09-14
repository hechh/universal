package parser

import (
	"strings"
	"universal/framework/uerror"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/typespec"

	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
)

type StructRule struct {
	fb    *excelize.File
	rule  string
	sheet string
	name  string
	class string
}

func (d *StructRule) Parse() error {
	values, err := d.fb.GetRows(d.sheet)
	if err != nil {
		return uerror.NewUError(1, -1, "读取配置表%s失败: %v", d.sheet, err)
	}
	// 解析结构
	fs := []*typespec.Field{}
	for i, val := range values[0] {
		if len(values[1]) <= i {
			values[1] = append(values[1], "")
		}
		val = strings.TrimSpace(val)
		if len(val) <= 0 {
			continue
		}
		// 解析结构
		if pos := strings.Index(val, "/"); pos <= 0 {
			fs = append(fs, &typespec.Field{
				Type:  manager.GetType(domain.KindTypeIdent, "", "string"),
				Name:  val,
				Index: i,
				Doc:   values[1][i],
			})
		} else {
			kind, pkg := manager.GetKindType(val[pos+1:])
			fs = append(fs, &typespec.Field{
				Type:  manager.GetType(kind, pkg, val[:pos]),
				Name:  val[:pos],
				Index: i,
				Doc:   values[1][i],
			})
		}
	}
	if len(fs) > 0 {
		manager.AddStruct(typespec.STRUCT(manager.GetType(domain.KindTypeStruct, domain.DefaultPkg, d.name), d.class, d.sheet, fs...))
	}
	return nil
}

type EnumRule struct {
	fb    *excelize.File
	rule  string
	sheet string
	class string
}

func (d *EnumRule) Parse() error {
	list, err := d.fb.GetRows(d.sheet)
	if err != nil {
		return uerror.NewUError(1, -1, "读取配置表%s失败: %v", d.sheet, err)
	}
	for _, vals := range list {
		for _, val := range vals {
			if !strings.HasPrefix(val, domain.RuleTypeEnum) {
				continue
			}

			addEval(d.class, d.sheet, val)
		}
	}
	return nil
}

func addEval(class, doc, val string) {
	ss := strings.Split(val, ":")
	ev := manager.LoadEnum(manager.GetType(domain.KindTypeEnum, domain.DefaultPkg, ss[2]))

	ev.Class = class
	ev.Doc = doc
	ev.AddValue(ss[3], cast.ToInt32(ss[4]), ss[1])
}
