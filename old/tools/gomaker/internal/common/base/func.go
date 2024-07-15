package base

import (
	"strings"
	"unicode"
)

type BaseFunc struct{}

func (d *BaseFunc) TrimPrefix(str, prefix string) string {
	return strings.TrimPrefix(str, prefix)
}

func (d *BaseFunc) TrimSuffix(str, prefix string) string {
	return strings.TrimSuffix(str, prefix)
}

func (d *BaseFunc) Split(str, sp string) []string {
	return strings.Split(str, sp)
}

func (d *BaseFunc) Join(a, b string) string {
	return a + b
}

// 首字符小写
func (d *BaseFunc) FirstCharToLower(str string) string {
	return string(unicode.ToLower(rune(str[0]))) + str[1:]
}

// 首字符大写
func (d *BaseFunc) FirstCharToUpper(str string) string {
	return string(unicode.ToUpper(rune(str[0]))) + str[1:]
}
