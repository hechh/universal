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

type XlsxParser struct{}

func (d *XlsxParser) ParseFiles(files ...string) error {
	// 解析所有配置的生成表
	for _, filename := range files {
		if err := ParseGenTable(filename, manager.GetMEnumsPointer(), manager.GetMessagePointer()); err != nil {
			return err
		}
	}
	// 解析枚举类型
	for _, sh := range manager.GetMEnumList() {
		values, err := sh.GetRows()
		if err != nil {
			return uerror.NewUError(1, -1, "读取配置表%s失败: %v", sh.Sheet, err)
		}
		for _, vals := range values {
			for _, val := range vals {
				if value := ParseXlsxEnum(sh.Class, val); value != nil {
					manager.LoadEnum(value.Type).Set(sh.Sheet).AddValue(value)
				}
			}
		}
	}
	manager.InitEvals()
	// 解析message信息
	for _, sh := range manager.GetMessageList() {
		values, err := sh.GetRows()
		if err != nil {
			return uerror.NewUError(1, -1, "读取配置表%s失败: %v", sh.Sheet, err)
		}
		// 解析结构
		if st := ParseXlsxStruct(sh, values[0], values[1]); st != nil {
			manager.AddStruct(st)
		}
	}
	manager.InitSheet()
	return nil
}

// 解析生成表
func ParseGenTable(filename string, enums, cfgs *[]*typespec.Sheet) (err error) {
	// 打开文件
	fp, err := excelize.OpenFile(filename)
	if err != nil {
		return uerror.NewUError(1, -1, "开打%s失败：%v", filename, err)
	}
	// 读取生成表
	values, err := fp.GetRows(domain.GenTable)
	if _, ok := err.(excelize.ErrSheetNotExist); ok || len(values) <= 0 {
		return nil
	}
	if err != nil {
		return uerror.NewUError(1, -1, "读取%s配置表%s失败: %v", filename, domain.GenTable, err)
	}
	// 解析规则
	for _, vals := range values {
		for _, val := range vals {
			ss := strings.Split(val, "|")
			switch ss[0] {
			case domain.GomakerTypeMessage:
				if cfgs != nil {
					*cfgs = append(*cfgs, ParseXlsxSheet(ss, fp))
				}
			case domain.GomakerTypeEnum:
				if enums != nil {
					*enums = append(*enums, ParseXlsxSheet(ss, fp))
				}
			}
		}
	}
	return
}

// 生成表解析
// @gomaker:类型｜out:生成文件名｜sheet:需要生成的配置名:配置的pb结构名称
func ParseXlsxSheet(ss []string, fp *excelize.File) *typespec.Sheet {
	result := typespec.NewSheet(ss[0], fp)
	for _, str := range ss[1:] {
		vals := strings.Split(str, ":")
		switch vals[0] {
		case "out":
			result.Class = vals[1]
		case "sheet":
			result.Sheet = vals[1]
			if len(vals) > 2 {
				result.Config = vals[2]
			}
		case "struct":
			result.IsStruct = true
		case "list":
			result.IsList = true
		case "map":
			tmps := []*typespec.Field{}
			for _, param := range strings.Split(vals[1], ",") {
				tmps = append(tmps, &typespec.Field{Name: param[:strings.Index(param, "@")]})
			}
			result.Map = append(result.Map, tmps)
		case "group":
			tmps := []*typespec.Field{}
			for _, param := range strings.Split(vals[1], ",") {
				tmps = append(tmps, &typespec.Field{Name: param[:strings.Index(param, "@")]})
			}
			result.Group = append(result.Group, tmps)
		}
	}
	return result
}

// xlsx枚举规则解析
// E:中文注释:枚举类型:枚举成员:枚举值
func ParseXlsxEnum(class, val string) *typespec.Value {
	if !strings.HasPrefix(val, domain.RuleTypeEnum) {
		return nil
	}
	ss := strings.Split(val, ":")
	tt := manager.GetType(domain.KindTypeEnum, domain.SourceTypeXlsx, domain.DefaultPkg, ss[2], class)
	return typespec.VALUE(tt, ss[3], cast.ToInt32(ss[4]), ss[1])
}

// 解析配置表中的message结构信息
func ParseXlsxStruct(sh *typespec.Sheet, val01, val02 []string) *typespec.Struct {
	// 没有注释的字段设置为空
	for j := len(val01); j < len(val01); j++ {
		val02 = append(val02, "")
	}
	// 解析结构
	fs := []*typespec.Field{}
	for i, val := range val01 {
		// 过滤空字段
		if val = strings.TrimSpace(val); len(val) <= 0 {
			continue
		}
		// 解析字段
		if ff := ParseXlsxField(i, val, val02[i]); ff != nil {
			fs = append(fs, ff)
		}
	}
	if len(fs) > 0 {
		return typespec.STRUCT(manager.GetType(domain.KindTypeStruct, domain.SourceTypeXlsx, domain.DefaultPkg, sh.Config, sh.Class), sh.Sheet, fs...)
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
	source := int32(domain.SourceTypeXlsx)
	fname := ff[:i]
	cfgType := strings.ReplaceAll(ff[i+1:], "[]", "")
	goType := manager.GetGoType(cfgType)
	pkg, kind := "", int32(domain.KindTypeIdent)
	if !domain.BasicTypes[goType] {
		// 枚举类型一定先于struct结构
		if ttt := manager.GetType(0, source, domain.DefaultPkg, goType, ""); ttt != nil {
			pkg, kind = ttt.Pkg, ttt.Kind
		} else {
			pkg, kind = domain.DefaultPkg, domain.KindTypeStruct
		}
	}
	// 解析token
	ts := []rune{}
	for i := 0; i < strings.Count(ff[i+1:], "[]"); i++ {
		ts = append(ts, domain.TokenTypeArray)
	}
	return typespec.FIELD(manager.GetType(kind, int32(source), pkg, goType, ""), fname, pos, cfgType, doc, ts...)
}
