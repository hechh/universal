package fbasic

import "time"

func GetNow() int64 {
	return time.Now().Unix()
}
