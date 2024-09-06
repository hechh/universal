package util

import "time"

func GetNowTime() time.Time {
	return time.Now()
}

func GetNowUnixSecond() int64 {
	return time.Now().Unix()
}

func GetNowUnixMilli() int64 {
	return time.Now().UnixMilli()
}

func GetNowUnixNano() int64 {
	return time.Now().UnixNano()
}
