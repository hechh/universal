package util

import "time"

func OR[T any](flag bool, a, b T) T {
	if flag {
		return a
	}
	return b
}

func INDEX[T any](arr []T, pos int, def T) T {
	if ll := len(arr); ll <= 0 || pos < 0 || pos >= ll {
		return def
	}
	return arr[pos]
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
