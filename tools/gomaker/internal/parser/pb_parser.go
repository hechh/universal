package parser

import (
	"fmt"
	"strings"
	"universal/tools/gomaker/domain"
	"universal/tools/gomaker/internal/typespec"
	"universal/tools/gomaker/internal/util"

	"github.com/spf13/cast"
)

type PbParser struct {
	docs   []string // 顶部注释
	buf    []byte   // 文件数据
	size   int      // buf大小
	begin  int      // 上一个位置
	cursor int      // 当前游标
	char   rune     // 当前字符
}

// 移动游标
func (d *PbParser) next(times int) *PbParser {
	d.cursor += times
	if d.size <= d.cursor {
		d.cursor = d.size
		d.char = -1
	} else {
		d.char = rune(d.buf[d.cursor])
	}
	return d
}

// 移动游标
func (d *PbParser) prev(times int) *PbParser {
	d.cursor -= times
	if d.cursor < 0 {
		d.cursor = 0
	}
	d.char = rune(d.buf[d.cursor])
	return d
}

func (d *PbParser) refresh() *PbParser {
	d.begin = d.cursor
	return d
}

func (d *PbParser) read() string {
	return string(d.buf[d.begin:d.cursor])
}

func (d *PbParser) getChar() rune {
	for ; true; d.next(1).refresh() {
		if d.char == ' ' || d.char == '\r' || d.char == '\t' || d.char == '\n' {
			continue
		}
		break
	}
	return d.char
}

// 过滤指定字符，并统计tt出现次数
func (d *PbParser) skip(tt rune, chs ...rune) int {
	tmps := map[rune]int{tt: 0}
	for _, ch := range chs {
		tmps[ch] = 0
	}
	for ; true; d.next(1).refresh() {
		if _, ok := tmps[d.char]; ok {
			tmps[d.char]++
			continue
		}
		break
	}
	return tmps[tt]
}

// 解析字符集
func (d *PbParser) parseWord() string {
	for ; d.char != -1 && (util.IsLetter(d.char) || util.IsNumber(d.char)); d.next(1) {
	}
	val := d.read()
	d.refresh()
	return val
}

// 解析注释
func (d *PbParser) parseDoc() (str string) {
	d.next(1)
	switch d.char {
	case '/':
		for d.next(1).refresh(); d.char != -1 && d.char != '\n'; d.next(1) {
		}
		str = strings.TrimSpace(d.read())
		d.next(1).refresh()
	case '*':
		for d.next(2).refresh(); d.char != -1 && d.buf[d.cursor-1] != '*' && d.char != '/'; d.next(1) {
		}
		str = strings.TrimSpace(d.read())
		d.next(1).refresh()
	default:
		d.prev(1)
	}
	return
}

// 解析字符串
func (d *PbParser) parseString() string {
	for d.next(1).refresh(); d.char != -1 && (d.char != '"' || d.buf[d.cursor-1] == '\\' && d.char == '"'); d.next(1) {
	}
	name := d.read()
	d.next(1).refresh()
	return name
}

func (d *PbParser) set(buf []byte) {
	d.docs = d.docs[:0]
	d.buf = buf
	d.size = len(buf)
	d.begin = 0
	d.cursor = 0
	d.char = rune(buf[0])
}

// 解析import
func (d *PbParser) parseImport(word string) *typespec.Import {
	d.skip(' ', '\t', '\n', '\r') // 过滤空格
	pkg := d.parseString()        // 解析引用文件
	d.skip(';', ' ', '\r', '\t')  // 解析分号
	item := &typespec.Import{
		Docs:    d.docs,
		Type:    word,
		File:    pkg,
		Comment: d.parseDoc(),
	}
	d.docs = d.docs[:0]
	return item
}

// 解析option
func (d *PbParser) parseOption(word string) *typespec.Option {
	d.skip(' ', '\t', '\n', '\r')      // 过滤空格
	opname := d.parseWord()            // 解析键
	d.skip('=', ' ', '\t', '\r', '\n') // 解析等号
	val := d.parseString()             // 解析值
	d.skip(';', ' ', '\t', '\r')       // 解析分号
	item := &typespec.Option{
		Docs:    d.docs,
		Type:    word,
		Key:     opname,
		Value:   val,
		Comment: d.parseDoc(),
	}
	d.docs = d.docs[:0]
	return item
}

// 解析package
func (d *PbParser) parsePackage(word string) *typespec.Package {
	d.skip(' ', '\t', '\n', '\r') // 过滤空格
	pkgname := d.parseWord()      // 解析包名
	d.skip(';', ' ', '\t', '\r')  // 过滤空格
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
func (d *PbParser) parseSyntax(word string) *typespec.Syntax {
	d.skip('=', ' ', '\r', '\t', '\n') // 解析 = 号
	lit := d.parseString()             // 解析proto版本
	d.skip(';', ' ', '\t', '\r')       // 解析 ; 号
	item := &typespec.Syntax{
		Docs:    d.docs,
		Type:    word,
		Name:    lit,
		Comment: d.parseDoc(),
	}
	d.docs = d.docs[:0]
	return item
}

// 解析message信息
func (d *PbParser) parseMessage(word string) *typespec.Message {
	d.skip(' ', '\t', '\n', '\r')      // 过滤空格
	stname := d.parseWord()            // 解析结构名字
	d.skip('{', ' ', '\r', '\t', '\n') // 过滤 { 符号
	// 解析field
	fs := []*typespec.Attribute{}
	for {
		// 解析日志
		item := &typespec.Attribute{}
	inner:
		if d.getChar() == '/' {
			if str := d.parseDoc(); len(str) > 0 {
				item.Docs = append(item.Docs, str)
			}
			goto inner
		}
		d.skip(' ', '\t', '\n', '\r') // 过滤空格
		item.Type = d.parseWord()     // 解析字段类型
		d.skip(' ', '\t', '\n', '\r') // 过滤空格
		item.Name = d.parseWord()     // 解析字段名
		if item.Type == "repeated" {
			item.IsRepeat = true          // 是否为数组
			d.skip(' ', '\t', '\n', '\r') // 过滤空格
			item.Type, item.Name = item.Name, d.parseWord()
		}
		d.skip('=', ' ', '\t', '\n', '\r')     // 过滤 =
		item.Index = cast.ToInt(d.parseWord()) // 解析pb的index
		d.skip(';', ' ', '\r', '\t')           // 过滤 ;
		item.Comment = d.parseDoc()            // 解析注释
		fs = append(fs, item)
		if d.skip('}', ' ', '\r', '\t', '\n') == 1 {
			break
		}
	}
	item := &typespec.Message{
		Docs:       d.docs,
		Type:       word,
		Name:       stname,
		Attributes: fs,
	}
	d.docs = d.docs[:0]
	return item
}

func (d *PbParser) ParseFile(buf []byte) {
	d.set(buf)
loop:
	switch d.getChar() {
	case '/':
		if str := d.parseDoc(); len(str) > 0 {
			d.docs = append(d.docs, str)
		}
		goto loop
	case -1:
		return
	default:
		word := d.parseWord()
		switch word {
		case domain.SYNTAX:
			fmt.Println("----->", d.parseSyntax(word))
			goto loop
		case domain.PACKAGE:
			fmt.Println("----->", d.parsePackage(word))
			goto loop
		case domain.OPTION:
			fmt.Println("----->", d.parseOption(word))
			goto loop
		case domain.IMPORT:
			fmt.Println("----->", d.parseImport(word))
			goto loop
		case domain.MESSAGE:
			fmt.Println("----->", d.parseMessage(word))
			goto loop
		case domain.ENUM:
		}
	}
	return
}
