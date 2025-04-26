package convert

import (
	"universal/framework/basic"

	"github.com/spf13/cast"
)

func ToUint32(str string) interface{} {
	return cast.ToUint32(str)
}

func ToUint64(str string) interface{} {
	return cast.ToUint64(str)
}

func ToInt32(str string) interface{} {
	return cast.ToInt32(str)
}

func ToInt64(str string) interface{} {
	return cast.ToInt64(str)
}

func ToFloat32(str string) interface{} {
	return cast.ToFloat32(str)
}

func ToFloat64(str string) interface{} {
	return cast.ToFloat64(str)
}

func ToBool(str string) interface{} {
	return cast.ToBool(str)
}

func ToString(str string) interface{} {
	return str
}

func ToBytes(str string) interface{} {
	return basic.StringToBytes(str)
}

func ToUnixSecond(str string) interface{} {
	return basic.String2Time(str).Unix()
}

func ToUnixMilli(str string) interface{} {
	return basic.String2Time(str).UnixMilli()
}
