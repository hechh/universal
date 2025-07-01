package util

import (
	"sync"
	"time"
)

func Or[T any](flag bool, a, b T) T {
	if flag {
		return a
	}
	return b
}

func Index[T any](arr []T, pos int, def T) T {
	if ll := len(arr); ll <= 0 || pos < 0 || pos >= ll {
		return def
	}
	return arr[pos]
}

func Prefix[T any](arr []T, pos int) []T {
	if pos < 0 || pos >= len(arr) {
		return nil
	}
	return arr[:pos]
}

func Suffix[T any](arr []T, pos int) []T {
	if pos < 0 || pos >= len(arr) {
		return nil
	}
	return arr[pos:]
}

func Retry(attempts int, sleep time.Duration, f func() error) (err error) {
	for i := 0; i < attempts; i++ {
		if err = f(); err == nil {
			return
		}

		time.Sleep(sleep)
		sleep *= 2
	}
	return err
}

type INumber interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

func Add[T INumber](a, b T) T {
	return a + b
}

func Sub[T INumber](a, b T) T {
	return a - b
}

func Pool[T any]() *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			return new(T)
		},
	}
}

func ArrayPool[T any](size int) *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			return make([]T, size)
		},
	}
}
