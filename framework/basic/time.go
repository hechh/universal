package basic

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

func String2Time(str string) time.Time {
	tt, _ := time.Parse("2006-01-02 15:04:05", str)
	return tt
}
