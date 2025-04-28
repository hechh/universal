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

type TWO[T1, T2 any] struct {
	T1 T1
	T2 T2
}

type THREE[T1, T2, T3 any] struct {
	T1 T1
	T2 T2
	T3 T3
}

type FOUR[T1, T2, T3, T4 any] struct {
	T1 T1
	T2 T2
	T3 T3
	T4 T4
}

type FIVE[T1, T2, T3, T4, T5 any] struct {
	T1 T1
	T2 T2
	T3 T3
	T4 T4
	T5 T5
}
