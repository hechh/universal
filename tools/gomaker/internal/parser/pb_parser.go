package parser

import (
	"io/ioutil"
	"universal/framework/uerror"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/internal/typespec"
	"universal/tools/gomaker/internal/util"

	"github.com/spf13/cast"
)

type PbParser struct {
	docs []string // 顶部注释
}

// 解析.proto文件
func (d *PbParser) ParseFile(filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return uerror.NewUError(1, -1, "%v %v", filename, err)
	}
	sc := util.NewScanner(buf)
loop:
	d.docs = d.docs[:0]
	switch sc.GetChar() {
	case '/':
		if str := sc.ParseDoc(); len(str) > 0 {
			d.docs = append(d.docs, str)
		}
		goto loop
	case -1:
		return nil
	default:
		word := sc.ParseWord()
		switch word {
		case domain.SYNTAX:
			parseSyntax(sc, word, d.docs...)
			goto loop
		case domain.PACKAGE:
			parsePackage(sc, word, d.docs...)
			goto loop
		case domain.OPTION:
			parseOption(sc, word, d.docs...)
			goto loop
		case domain.IMPORT:
			parseImport(sc, word, d.docs...)
			goto loop
		case domain.MESSAGE:
			if item := parseMessage(sc, word, d.docs...); item != nil {
				manager.AddStruct(messageToStruct(item))
			}
			goto loop
		case domain.ENUM:
			if item := parseEnum(sc, word, d.docs...); item != nil {
				manager.AddEnum(menumToEnum(item))
			}
			goto loop
		}
	}
	return nil
}

// 解析import
func parseImport(sc *util.Scanner, word string, docs ...string) *typespec.Import {
	sc.Skip(' ', '\t', '\n', '\r') // 过滤空格
	pkg := sc.ParseString()        // 解析引用文件
	sc.Skip(';', ' ', '\r', '\t')  // 解析分号
	return &typespec.Import{
		Docs:    docs,
		Type:    word,
		File:    pkg,
		Comment: sc.ParseDoc(),
	}
}

// 解析option
func parseOption(sc *util.Scanner, word string, docs ...string) *typespec.Option {
	sc.Skip(' ', '\t', '\n', '\r')      // 过滤空格
	opname := sc.ParseWord()            // 解析键
	sc.Skip('=', ' ', '\t', '\r', '\n') // 解析等号
	val := sc.ParseString()             // 解析值
	sc.Skip(';', ' ', '\t', '\r')       // 解析分号
	return &typespec.Option{
		Docs:    docs,
		Type:    word,
		Key:     opname,
		Value:   val,
		Comment: sc.ParseDoc(),
	}
}

// 解析package
func parsePackage(sc *util.Scanner, word string, docs ...string) *typespec.Package {
	sc.Skip(' ', '\t', '\n', '\r') // 过滤空格
	pkgname := sc.ParseWord()      // 解析包名
	sc.Skip(';', ' ', '\t', '\r')  // 过滤空格
	// 解析注释
	return &typespec.Package{
		Docs:    docs,
		Type:    word,
		Name:    pkgname,
		Comment: sc.ParseDoc(),
	}
}

// 解析syntax
func parseSyntax(sc *util.Scanner, word string, docs ...string) *typespec.Syntax {
	sc.Skip('=', ' ', '\r', '\t', '\n') // 解析 = 号
	lit := sc.ParseString()             // 解析proto版本
	sc.Skip(';', ' ', '\t', '\r')       // 解析 ; 号
	return &typespec.Syntax{
		Docs:    docs,
		Type:    word,
		Name:    lit,
		Comment: sc.ParseDoc(),
	}
}

// 解析message信息
func parseMessage(sc *util.Scanner, word string, docs ...string) *typespec.Message {
	sc.Skip(' ', '\t', '\n', '\r')      // 过滤空格
	stname := sc.ParseWord()            // 解析结构名字
	sc.Skip('{', ' ', '\r', '\t', '\n') // 过滤 { 符号
	// 解析field
	fs := []*typespec.Attribute{}
	for {
		item := &typespec.Attribute{}
	inner:
		if sc.GetChar() == '/' {
			if str := sc.ParseDoc(); len(str) > 0 {
				item.Docs = append(item.Docs, str)
			}
			goto inner
		}
		sc.Skip(' ', '\t', '\n', '\r') // 过滤空格
		item.Type = sc.ParseWord()     // 解析字段类型
		sc.Skip(' ', '\t', '\n', '\r') // 过滤空格
		item.Name = sc.ParseWord()     // 解析字段名
		if item.Type == "repeated" {
			item.IsRepeat = true           // 是否为数组
			sc.Skip(' ', '\t', '\n', '\r') // 过滤空格
			item.Type, item.Name = item.Name, sc.ParseWord()
		}
		sc.Skip('=', ' ', '\t', '\n', '\r')     // 过滤 =
		item.Index = cast.ToInt(sc.ParseWord()) // 解析pb的index
		sc.Skip(';', ' ', '\r', '\t')           // 过滤 ;
		item.Comment = sc.ParseDoc()            // 解析注释
		fs = append(fs, item)
		if sc.Skip('}', ' ', '\r', '\t', '\n') == 1 {
			break
		}
	}
	return &typespec.Message{
		Docs:       docs,
		Type:       word,
		Name:       stname,
		Attributes: fs,
	}
}

func parseEnum(sc *util.Scanner, word string, docs ...string) *typespec.MEnum {
	sc.Skip(' ', '\t', '\n', '\r')      // 过滤空格
	ttName := sc.ParseWord()            // 枚举类型
	sc.Skip('{', ' ', '\r', '\t', '\n') // 过滤 { 符号
	// 解析field
	fs := []*typespec.MValue{}
	for {
		item := &typespec.MValue{}
	inner:
		if sc.GetChar() == '/' {
			if str := sc.ParseDoc(); len(str) > 0 {
				item.Docs = append(item.Docs, str)
			}
			goto inner
		}
		sc.Skip(' ', '\t', '\n', '\r')          // 过滤空格
		item.Name = sc.ParseWord()              // 解析字段类型
		sc.Skip('=', ' ', '\t', '\n', '\r')     // 过滤 =
		item.Value = cast.ToInt(sc.ParseWord()) // 解析字段名
		sc.Skip(';', ' ', '\r', '\t')           // 过滤 ;
		item.Comment = sc.ParseDoc()            // 解析注释
		fs = append(fs, item)
		if sc.Skip('}', ' ', '\r', '\t', '\n') == 1 {
			break
		}
	}
	return &typespec.MEnum{
		Docs:   docs,
		Type:   word,
		Name:   ttName,
		Values: fs,
	}
}

func menumToEnum(val *typespec.MEnum) *typespec.Enum {
	vals := []*typespec.Value{}
	for _, vv := range val.Values {
		ttt := manager.GetType(domain.KindTypeIdent, "", vv.Name, "")
		vals = append(vals, typespec.VALUE(ttt, vv.Name, int32(vv.Value), vv.Comment))
	}
	return typespec.ENUM(manager.GetType(domain.KindTypeEnum, domain.DefaultPkg, val.Name, ""), "", vals...)
}

func messageToStruct(val *typespec.Message) *typespec.Struct {
	fs := []*typespec.Field{}
	for _, vv := range val.Attributes {
		if _, ok := domain.BasicTypes[vv.Type]; ok {
			tt := manager.GetType(domain.KindTypeIdent, "", vv.Name, "")
			fs = append(fs, typespec.FIELD(tt, vv.Name, vv.Index, "", ""))
		} else {
			tt := manager.GetType(domain.KindTypeStruct, "", vv.Name, "")
			fs = append(fs, typespec.FIELD(tt, vv.Name, vv.Index, "", ""))
		}
	}
	return typespec.STRUCT(manager.GetType(domain.KindTypeStruct, domain.DefaultPkg, val.Name, ""), "", fs...)
}
