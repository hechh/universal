package random

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Shuffle(n int, swap func(i, j int)) {
	rand.Shuffle(n, swap)
}

// [0,n)
func Intn(n int) int {
	return rand.Intn(n)
}

// [0,n)
func Uint32n(n uint32) uint32 {
	return uint32(Intn(int(n)))
}

// [0,n)
func Int32n(n int32) int32 {
	return int32(Intn(int(n)))
}

// [0,n)
func Int64n(n int64) int64 {
	return rand.Int63n(n)
}

// 取一个范围的随机数
// [2, 5]
func Uint32Part(min, max uint32) uint32 {
	// [2, 5]
	n := max - min + 1
	return uint32(rand.Int31n(int32(n))) + min
}

func Int64Part(min, max int64) int64 {
	// [2, 5]
	n := max - min + 1
	return rand.Int63n(n) + min
}

func Int32Part(min, max int32) int32 {
	// [2, 5]
	n := max - min + 1
	return rand.Int31n(n) + min
}
