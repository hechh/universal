package util

import "sync"

func Ifelse[T any](flag bool, a, b T) T {
	if flag {
		return a
	}
	return b
}

func Pool[T any]() *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			return new(T)
		},
	}
}

func PoolSlice[T any](size int) *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			return make([]T, size)
		},
	}
}
