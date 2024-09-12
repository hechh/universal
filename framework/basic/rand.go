package basic

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 在[min, max]范围内随机生成一个值
func RangeInt63n(min, max int64) int64 {
	return rand.Int63n(max-min) + min
}

func RangeInt31n(min, max int32) int32 {
	return rand.Int31n(max-min) + min
}

func RangeIntn(min, max int) int {
	return rand.Intn(max-min) + min
}
