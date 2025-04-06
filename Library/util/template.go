package util

import (
	"unicode"
)

func Ifelse[T any](flag bool, a, b T) T {
	if flag {
		return a
	}
	return b
}

func GetPrefix(str string, pos int) string {
	if pos < 0 || pos >= len(str) {
		return ""
	}
	return str[:pos]
}

func GetSuffix(str string, pos int) string {
	if pos < 0 || pos >= len(str) {
		return ""
	}
	return str[pos:]
}

func Prefix[T any](str []T, pos int) []T {
	if pos < 0 || pos >= len(str) {
		return nil
	}
	return str[:pos]
}

func Suffix[T any](str []T, pos int) []T {
	if pos < 0 || pos >= len(str) {
		return nil
	}
	return str[pos:]
}

func ToLowerFirst(str string) string {
	rr := []rune(str)
	rr[0] = unicode.SimpleFold(rr[0])
	return string(rr)
}
