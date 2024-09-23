package util

import (
	"unicode"
	"unicode/utf8"
)

type Scanner struct {
	buf    []byte // 文件数据
	size   int    // buf大小
	begin  int    // 上一个位置
	cursor int    // 当前游标
	char   rune   // 当前字符
}

func IsSpace(ch rune) bool {
	return ' ' == ch || '\t' == ch || '\n' == ch || '\r' == ch
}

func IsNumber(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func IsDigit(ch rune) bool {
	return IsNumber(ch) || ch >= utf8.RuneSelf && unicode.IsDigit(ch)
}

func IsLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch >= utf8.RuneSelf && unicode.IsLetter(ch)
}

func NewScanner(buf []byte) *Scanner {
	return &Scanner{buf: buf, size: len(buf), char: rune(buf[0])}
}

// 移动游标
func (d *Scanner) Next(times int) *Scanner {
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
func (d *Scanner) Prev(times int) *Scanner {
	d.cursor -= times
	if d.cursor < 0 {
		d.cursor = 0
	}
	d.char = rune(d.buf[d.cursor])
	return d
}

func (d *Scanner) Refresh() *Scanner {
	d.begin = d.cursor
	return d
}

func (d *Scanner) Read() string {
	return string(d.buf[d.begin:d.cursor])
}

func (d *Scanner) GetChar() rune {
	for ; true; d.Next(1).Refresh() {
		if d.char == ' ' || d.char == '\r' || d.char == '\t' || d.char == '\n' {
			continue
		}
		break
	}
	return d.char
}

// 过滤指定字符，并统计tt出现次数
func (d *Scanner) Skip(tt rune, chs ...rune) int {
	tmps := map[rune]int{tt: 0}
	for _, ch := range chs {
		tmps[ch] = 0
	}
	for ; true; d.Next(1).Refresh() {
		if _, ok := tmps[d.char]; ok {
			tmps[d.char]++
			continue
		}
		break
	}
	return tmps[tt]
}
