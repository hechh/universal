package parser

import (
	"strings"
	"universal/framework/uerror"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/typespec"
	"universal/tools/gomaker/internal/util"

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
				if value := util.ParseXlsxEnum(val); value != nil {
					manager.LoadEnum(value.Type).Set(domain.DefaultEnumClass, "").AddValue(value)
					continue
				}
				// 解析@gomaker
				if item := util.ParseXlsxSheet(val, fp); item != nil {
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
			if ff := util.ParseXlsxField(i, val, values[1][i]); ff != nil {
				fs = append(fs, ff)
			}
		}
		if len(fs) > 0 {
			tt := manager.GetType(domain.KindTypeStruct, domain.DefaultPkg, sh.Config)
			manager.AddStruct(typespec.STRUCT(tt, sh.Class, sh.Sheet, fs...))
		}
	}
	return nil
}

func (d *XlsxParser) parseEnum() error {
	for _, sh := range d.enums {
		values, err := sh.GetRows()
		if err != nil {
			return uerror.NewUError(1, -1, "读取配置表%s失败: %v", sh.Sheet, err)
		}
		for _, vals := range values {
			for _, val := range vals {
				if value := util.ParseXlsxEnum(val); value != nil {
					manager.LoadEnum(value.Type).Set(sh.Class, sh.Sheet).AddValue(value)
				}
			}
		}
	}
	return nil
}
