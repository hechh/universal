package parse

import (
	"strings"
	"universal/framework/basic"
	"universal/tools/gomaker/internal/manager"

	"github.com/spf13/cast"
)

func init() {
	manager.AddConv("uint32", "", ToUint32)
	manager.AddConv("uint64", "", ToUint64)
	manager.AddConv("int32", "", ToInt32)
	manager.AddConv("int64", "", ToInt64)
	manager.AddConv("float32", "", ToFloat32)
	manager.AddConv("float64", "", ToFloat64)
	manager.AddConv("bool", "", ToBool)
	manager.AddConv("string", "", ToString)
	manager.AddConv("time", "", ToTime)
	manager.AddConv("[]byte", "", ToBytes)
	manager.AddConv("[]uint32", "", ToArrayUint32)
	manager.AddConv("[]uint64", "", ToArrayUint64)
	manager.AddConv("[]int32", "", ToArrayInt32)
	manager.AddConv("[]int64", "", ToArrayInt64)
	manager.AddConv("[]float32", "", ToArrayFloat32)
	manager.AddConv("[]float64", "", ToArrayFloat64)
	manager.AddConv("[]bool", "", ToArrayBool)
	manager.AddConv("[]string", "", ToArrayString)
}

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

func ToTime(str string) interface{} {
	return basic.String2Time(str).Unix()
}

func ToBytes(str string) interface{} {
	return basic.StringToBytes(str)
}

func ToArrayUint32(str string) interface{} {
	rets := []uint32{}
	for _, val := range strings.Split(str, ",") {
		rets = append(rets, cast.ToUint32(val))
	}
	return rets
}

func ToArrayUint64(str string) interface{} {
	rets := []uint64{}
	for _, val := range strings.Split(str, ",") {
		rets = append(rets, cast.ToUint64(val))
	}
	return rets
}

func ToArrayInt32(str string) interface{} {
	rets := []int32{}
	for _, val := range strings.Split(str, ",") {
		rets = append(rets, cast.ToInt32(val))
	}
	return rets
}

func ToArrayInt64(str string) interface{} {
	rets := []int64{}
	for _, val := range strings.Split(str, ",") {
		rets = append(rets, cast.ToInt64(val))
	}
	return rets
}

func ToArrayFloat32(str string) interface{} {
	rets := []float32{}
	for _, val := range strings.Split(str, ",") {
		rets = append(rets, cast.ToFloat32(val))
	}
	return rets
}

func ToArrayFloat64(str string) interface{} {
	rets := []float64{}
	for _, val := range strings.Split(str, ",") {
		rets = append(rets, cast.ToFloat64(val))
	}
	return rets
}

func ToArrayBool(str string) interface{} {
	rets := []bool{}
	for _, val := range strings.Split(str, ",") {
		rets = append(rets, cast.ToBool(val))
	}
	return rets
}

func ToArrayString(str string) interface{} {
	return strings.Split(str, ",")
}
