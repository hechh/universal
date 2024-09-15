package util

import (
	"strings"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/typespec"

	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
)

// xlsx枚举规则解析
// E:中文注释:枚举类型:枚举成员:枚举值
func ParseXlsxEnum(val string) *typespec.Value {
	if !strings.HasPrefix(val, domain.RuleTypeEnum) {
		return nil
	}
	ss := strings.Split(val, ":")
	tt := manager.GetType(domain.KindTypeEnum, domain.DefaultPkg, ss[2])
	return typespec.VALUE(tt, ss[3], cast.ToInt32(ss[4]), "")
}

// 生成表解析
// @gomaker:类型｜生成文件名｜需要生成的配置名:配置的pb结构名称
func ParseXlsxSheet(val string, fp *excelize.File) *typespec.Sheet {
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
func ParseXlsxField(pos int, ff string, doc string) *typespec.Field {
	i := strings.Index(ff, "/")
	if i <= 0 || len(ff[i+1:]) <= 0 {
		return nil
	}
	fname := ff[:i]
	cfgType := strings.ReplaceAll(ff[i+1:], "[]", "")
	goType := manager.GetGoType(cfgType)
	pkg, kind := "", int32(domain.KindTypeIdent)
	if !domain.BasicTypes[goType] {
		pkg, kind = domain.DefaultPkg, domain.KindTypeProto
	}
	// 解析token
	ts := []int32{}
	for i := 0; i < strings.Count(ff[i+1:], "[]"); i++ {
		ts = append(ts, domain.TokenTypeArray)
	}
	return typespec.FIELD(manager.GetType(kind, pkg, goType), fname, pos, cfgType, doc, ts...)
}
