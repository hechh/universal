package util

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

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

// 过滤空字节
func SkipSpace(buf []byte, begin, cursor *int) {
	for len(buf) > *cursor && IsSpace(rune(buf[*cursor])) {
		(*cursor)++
	}
	(*begin) = *cursor
}

func ScanDoc(buf []byte, begin, cursor *int) string {
	// 过滤空格
	SkipSpace(buf, begin, cursor)
	switch string(buf[*cursor : *cursor+1]) {
	case "//":
		beg := *cursor + int(2)
		for buf[*cursor] != '\n' {
			(*cursor)++
		}
		doc := string(buf[beg:*cursor])
		*begin = *cursor
		return strings.TrimSpace(doc)
	case "/*":
		beg := *cursor + 2
		for string(buf[*cursor-1:*cursor]) == "*/" {
			*cursor++
		}
		doc := string(buf[beg : *cursor-1])
		*cursor++
		*begin = *cursor
		return strings.TrimSpace(doc)
	}
	return ""
}
