package parser

import (
	"fmt"
	"strings"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/typespec"
	"universal/tools/gomaker/internal/util"
)

type PbParser struct {
	docs   []string // 顶部注释
	buf    []byte   // 文件数据
	size   int      // buf大小
	begin  int      // 上一个位置
	cursor int      // 当前游标
	char   rune     // 当前字符
}

func NewPbParser(buf []byte) *PbParser {
	return &PbParser{buf: buf, size: len(buf), cursor: 0, char: rune(buf[0])}
}

// 移动游标
func (d *PbParser) next() {
	d.cursor++
	if d.size <= d.cursor {
		d.cursor = d.size
		d.char = -1
	} else {
		d.char = rune(d.buf[d.cursor])
	}
}

// 移动游标
func (d *PbParser) prev() {
	d.cursor--
	if d.cursor < 0 {
		d.cursor = 0
	}
	d.char = rune(d.buf[d.cursor])
}

func (d *PbParser) forward(vals ...int) {
	times := 1
	if len(vals) > 0 && vals[0] > 0 {
		times = vals[0]
	}
	for i := 0; i < times && d.cursor < d.size; i++ {
		d.next()
		d.begin = d.cursor
	}
}

func (d *PbParser) back(vals ...int) {
	times := 1
	if len(vals) > 0 && vals[0] > 0 {
		times = vals[0]
	}
	for i := 0; i < times && d.cursor < d.size; i++ {
		d.prev()
		d.begin = d.cursor
	}
}

func (d *PbParser) eof() bool {
	return d.char == -1
}

func (d *PbParser) read() string {
	return string(d.buf[d.begin:d.cursor])
}

func (d *PbParser) reset(vals ...int) {
	if d.cursor >= d.size {
		d.begin = d.cursor
	} else {
		d.forward(vals...)
	}
}

// 过滤指定字符，并统计tt出现次数
func (d *PbParser) skip(tt rune, chs ...rune) int {
	tmps := map[rune]int{tt: 0}
	for _, ch := range chs {
		tmps[ch] = 0
	}
	for ; true; d.forward() {
		if _, ok := tmps[d.char]; ok {
			tmps[d.char]++
			continue
		}
		break
	}
	return tmps[tt]
}

// 解析字符串
func (d *PbParser) parseString() string {
	d.skip(' ', '\r', '\t')
	if d.char != '"' {
		return ""
	}
	for d.forward(); d.char != '"' || d.buf[d.cursor-1] == '\\' && d.char == '"'; d.next() {
	}
	doc := d.read()
	d.reset()
	return strings.Trim(doc, "\"")
}

// 解析注释
func (d *PbParser) parseDoc() (str string) {
	switch d.char {
	case '/':
		d.next()
		switch d.char {
		case '/':
			for d.forward(); !d.eof() && d.char != '\n'; d.next() {
			}
			str = strings.TrimSpace(d.read())
			d.reset()
		case '*':
			for d.forward(); !d.eof() && d.buf[d.cursor+1] != '/' && d.char != '*'; d.next() {
			}
			str = strings.TrimSpace(d.read())
			d.reset()
		}
	}
	return
}

// 解析特殊关键字
func (d *PbParser) parseWord() string {
	for ; !d.eof() && (util.IsLetter(d.char) || util.IsNumber(d.char)); d.next() {
	}
	name := d.read()
	d.reset()
	return name
}

func (d *PbParser) Parse() {
loop:
	d.skip(' ', '\r', '\t', '\n')
	switch d.char {
	case '/':
		if str := d.parseDoc(); len(str) > 0 {
			d.docs = append(d.docs, str)
		}
		goto loop
	case -1: // 文件终止
		return
	default: // 注释
		name := d.parseWord()
		fmt.Println("--->", name)
		switch name {
		case domain.SYNTAX:
			switch vv := d.parseSyntax(name).(type) {
			case error:
				fmt.Println("=====err=======>", vv)
			case *typespec.Syntax:
				fmt.Println("============>", vv)
			}
			goto loop
		case domain.PACKAGE:
			switch vv := d.parsePackage(name).(type) {
			case error:
				fmt.Println("=====err=======>", vv)
			case *typespec.Package:
				fmt.Println("============>", vv)
			}
			goto loop
		case domain.IMPORT:
			switch vv := d.parseImport(name).(type) {
			case error:
				fmt.Println("=====err=======>", vv)
			case *typespec.Import:
				fmt.Println("============>", vv)
			}
			goto loop
		case domain.OPTION:
			switch vv := d.parseOption(name).(type) {
			case error:
				fmt.Println("=====err=======>", vv)
			case *typespec.Option:
				fmt.Println("============>", vv)
			}
			goto loop
		case domain.MESSAGE:
			/*
				switch vv := d.parseMessage(name).(type) {
				case error:
					fmt.Println("=====err=======>", vv)
				case *typespec.Message:
					fmt.Println("============>", vv)
				}
				goto loop
			*/
			return
		case domain.ENUM:
			return
		}
		return
	}
}

func (d *PbParser) parseMessage(word string) interface{} {
	// 解析结构名字
	stname := d.parseWord()
	if len(stname) <= 0 {
		return fmt.Errorf("message语法错误")
	}

	return nil
}

// 解析import
func (d *PbParser) parseImport(word string) interface{} {
	d.skip(' ', '\r', '\t', '\n')
	pkg := d.parseString()
	// 解析分号
	if times := d.skip(';', ' ', '\r', '\t', '\n'); times != 1 {
		return fmt.Errorf("import语法错误")
	}
	item := &typespec.Import{
		Docs:    d.docs,
		File:    pkg,
		Comment: d.parseDoc(),
	}
	d.docs = d.docs[:0]
	return item
}

// 解析option
func (d *PbParser) parseOption(word string) interface{} {
	d.skip(' ', '\r', '\t', '\n')
	// 解析键
	opname := d.parseWord()
	if len(opname) <= 0 {
		return fmt.Errorf("option语法错误")
	}
	// 解析等号
	if times := d.skip('=', ' ', '\t', '\r', '\n'); times != 1 {
		return fmt.Errorf("option语法错误")
	}
	// 解析值
	val := d.parseString()
	// 解析分号
	if times := d.skip(';', ' ', '\t', '\r', '\n'); times != 1 {
		return fmt.Errorf("option语法错误")
	}
	item := &typespec.Option{
		Docs:    d.docs,
		Key:     opname,
		Value:   val,
		Comment: d.parseDoc(),
	}
	d.docs = d.docs[:0]
	return item
}

// 解析package
func (d *PbParser) parsePackage(word string) interface{} {
	d.skip(' ', '\r', '\t', '\n')
	if !util.IsLetter(d.char) {
		return fmt.Errorf("package语法错误")
	}
	// 解析包名
	for ; !d.eof() && d.char != ';'; d.next() {
	}
	pkgname := d.read()
	if len(pkgname) <= 0 {
		return fmt.Errorf("package语法错误, 包含特殊字符")
	}
	if times := d.skip(';', ' ', '\r', '\t', '\n'); times != 1 {
		return fmt.Errorf("package语法错误")
	}
	// 解析注释
	item := &typespec.Package{
		Docs:    d.docs,
		Type:    word,
		Name:    pkgname,
		Comment: d.parseDoc(),
	}
	d.docs = d.docs[:0]
	return item
}

// 解析syntax
func (d *PbParser) parseSyntax(word string) interface{} {
	// 解析 = 号
	if times := d.skip('=', ' ', '\r', '\t', '\n'); times != 1 {
		return fmt.Errorf("syntax语法错误, 没有=")
	}
	// 解析proto版本
	lit := d.parseString()
	if len(lit) <= 0 {
		return fmt.Errorf("syntax语法错误, proto版本为空")
	}
	// 解析 ; 号
	if times := d.skip(';', ' ', '\t', '\r'); times != 1 {
		return fmt.Errorf("syntax语法错误，没有;")
	}
	item := &typespec.Syntax{
		Docs:    d.docs,
		Type:    word,
		Name:    lit,
		Comment: d.parseDoc(),
	}
	d.docs = d.docs[:0]
	return item
}
