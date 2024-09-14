package parser

import (
	"strings"
	"universal/framework/uerror"
	"universal/tools/gomaker_new/domain"
	"universal/tools/gomaker_new/internal/manager"

	"github.com/xuri/excelize/v2"
)

type XlsxParser struct {
	prev []*EnumRule
	next []*StructRule
}

func (d *XlsxParser) Parse() error {
	// 先解析枚举
	for _, item := range d.prev {
		if err := item.Parse(); err != nil {
			return err
		}
	}
	manager.InitEvals()
	// 后解析结构
	for _, item := range d.next {
		if err := item.Parse(); err != nil {
			return err
		}
	}
	return nil
}

// 解析生成表
func (d *XlsxParser) ParseFiles(files ...string) error {
	for _, filename := range files {
		fb, err := excelize.OpenFile(filename)
		if err != nil {
			return uerror.NewUError(1, -1, "开打%s失败：%v", filename, err)
		}
		// 读取生成表
		values, err := fb.GetRows(domain.GenTable)
		if _, ok := err.(excelize.ErrSheetNotExist); ok || len(values) <= 0 {
			err = nil
			continue
		}
		if err != nil {
			return uerror.NewUError(1, -1, "读取%s配置表%s失败: %v", filename, domain.GenTable, err)
		}
		// 解析生成表
		for _, vals := range values {
			for _, val := range vals {
				// 解析枚举
				if strings.HasPrefix(val, domain.RuleTypeEnum) {
					addEval(domain.DefaultEnumClass, "", val)
					continue
				}
				// 解析生成表
				ss := strings.Split(val, "|")
				switch ss[0] {
				case domain.RuleTypeBytes:
					pos := strings.Index(ss[1], ":")
					d.next = append(d.next, &StructRule{fb: fb, rule: ss[0], sheet: ss[1][:pos], name: ss[1][pos+1:], class: ss[2]})
				case domain.RuleTypeProto:
					d.prev = append(d.prev, &EnumRule{fb: fb, rule: ss[0], sheet: ss[1], class: ss[2]})
				}
			}
		}
	}
	return nil
}
