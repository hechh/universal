package util

import "sync"

func Ifelse[T any](flag bool, a, b T) T {
	if flag {
		return a
	}
	return b
}

func Prefix[T any](vals []T, pos int) []T {
	if pos < 0 || pos >= len(vals) {
		return nil
	}
	return vals[:pos]
}

func Suffix[T any](vals []T, pos int) []T {
	if pos < 0 || pos >= len(vals) {
		return nil
	}
	return vals[pos:]
}

func NewSyncPool[T any]() *sync.Pool {
	return &sync.Pool{New: func() interface{} { return new(T) }}
}
