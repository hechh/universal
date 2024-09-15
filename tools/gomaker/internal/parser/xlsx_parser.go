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

type XlsxParser struct {
	enums []*typespec.Sheet
	cfgs  []*typespec.Sheet
}

func NewXlsxParser(files ...string) (*XlsxParser, error) {
	enums := []*typespec.Sheet{}
	cfgs := []*typespec.Sheet{}
	for _, filename := range files {
		// 打开文件
		fp, err := excelize.OpenFile(filename)
		if err != nil {
			return nil, uerror.NewUError(1, -1, "开打%s失败：%v", filename, err)
		}
		// 读取生成表
		values, err := fp.GetRows(domain.GenTable)
		if _, ok := err.(excelize.ErrSheetNotExist); ok || len(values) <= 0 {
			err = nil
			continue
		}
		if err != nil {
			return nil, uerror.NewUError(1, -1, "读取%s配置表%s失败: %v", filename, domain.GenTable, err)
		}
		// 解析生成表
		for _, vals := range values {
			for _, val := range vals {
				// 解析E:
				if value := parseXlsxEnum(domain.DefaultEnumClass, val); value != nil {
					manager.LoadEnum(value.Type).AddValue(value)
					continue
				}
				// 解析@gomaker
				if item := parseXlsxSheet(val, fp); item != nil {
					switch item.Rule {
					case domain.RuleTypeBytes:
						cfgs = append(cfgs, item)
					case domain.RuleTypeProto:
						enums = append(enums, item)
					}
				}
			}
		}
	}
	return &XlsxParser{enums: enums, cfgs: cfgs}, nil
}

func (d *XlsxParser) Parse() error {
	// 先解析枚举
	if err := d.parseEnum(); err != nil {
		return err
	}
	manager.InitEvals()
	// 后解析结构
	return d.parseCfg()
}

func (d *XlsxParser) parseEnum() error {
	for _, sh := range d.enums {
		values, err := sh.GetRows()
		if err != nil {
			return uerror.NewUError(1, -1, "读取配置表%s失败: %v", sh.Sheet, err)
		}
		for _, vals := range values {
			for _, val := range vals {
				if value := parseXlsxEnum(sh.Class, val); value != nil {
					manager.LoadEnum(value.Type).Set(sh.Sheet).AddValue(value)
				}
			}
		}
	}
	return nil
}

func (d *XlsxParser) parseCfg() error {
	for _, sh := range d.cfgs {
		values, err := sh.GetRows()
		if err != nil {
			return uerror.NewUError(1, -1, "读取配置表%s失败: %v", sh.Sheet, err)
		}
		// 没有注释的字段设置为空
		for j := len(values[1]); j < len(values[0]); j++ {
			values[1] = append(values[1], "")
		}
		// 解析结构
		fs := []*typespec.Field{}
		for i, val := range values[0] {
			// 过滤空字段
			if val = strings.TrimSpace(val); len(val) <= 0 {
				continue
			}
			// 解析字段
			if ff := parseXlsxField(i, val, values[1][i]); ff != nil {
				fs = append(fs, ff)
			}
		}
		if len(fs) > 0 {
			tt := manager.GetType(domain.KindTypeStruct, domain.DefaultPkg, sh.Config, sh.Class)
			manager.AddStruct(typespec.STRUCT(tt, sh.Sheet, fs...))
		}
	}
	return nil
}

// xlsx枚举规则解析
// E:中文注释:枚举类型:枚举成员:枚举值
func parseXlsxEnum(class, val string) *typespec.Value {
	if !strings.HasPrefix(val, domain.RuleTypeEnum) {
		return nil
	}
	ss := strings.Split(val, ":")
	tt := manager.GetType(domain.KindTypeEnum, domain.DefaultPkg, ss[2], class)
	return typespec.VALUE(tt, ss[3], cast.ToInt32(ss[4]), ss[1])
}

// 生成表解析
// @gomaker:类型｜生成文件名｜需要生成的配置名:配置的pb结构名称
func parseXlsxSheet(val string, fp *excelize.File) *typespec.Sheet {
	if !strings.HasPrefix(val, domain.RuleTypeGomaker) {
		return nil
	}
	ss := strings.Split(val, "|")
	switch ss[0] {
	case domain.RuleTypeProto:
		return typespec.SHEET(ss[0], ss[1], ss[2], "", fp)
	case domain.RuleTypeBytes:
		if pos := strings.Index(ss[2], ":"); pos > 0 {
			return typespec.SHEET(ss[0], ss[1], ss[2][:pos], ss[2][pos+1:], fp)
		} else {
			return typespec.SHEET(ss[0], ss[1], ss[2], ss[2], fp)
		}
	}
	return nil
}

// 解析字段
// 字段名/[][]配置类型
func parseXlsxField(pos int, ff string, doc string) *typespec.Field {
	i := strings.Index(ff, "/")
	if i <= 0 || len(ff[i+1:]) <= 0 {
		return nil
	}
	fname := ff[:i]
	cfgType := strings.ReplaceAll(ff[i+1:], "[]", "")
	goType := manager.GetGoType(cfgType)
	pkg, kind := "", int32(domain.KindTypeIdent)
	if !domain.BasicTypes[goType] {
		if ttt := manager.GetType(0, domain.DefaultPkg, goType, ""); ttt != nil {
			pkg, kind = ttt.Pkg, ttt.Kind
		} else {
			pkg, kind = domain.DefaultPkg, domain.KindTypeStruct
		}
	}
	// 解析token
	ts := []int32{}
	for i := 0; i < strings.Count(ff[i+1:], "[]"); i++ {
		ts = append(ts, domain.TokenTypeArray)
	}
	return typespec.FIELD(manager.GetType(kind, pkg, goType, ""), fname, pos, cfgType, doc, ts...)
}
