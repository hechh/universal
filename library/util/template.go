package util

import (
	"sync"
	"time"
)

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
