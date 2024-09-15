package repository

import (
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/repository/convert"
	"universal/tools/gomaker/repository/generator"
)

func Init() {
	manager.Register("proto", "xlsx转pb结构", generator.EnumGen, generator.TableGen)
}

func init() {
	manager.AddConv("uint32", convert.ToUint32)
	manager.AddConv("uint64", convert.ToUint64)
	manager.AddConv("int32", convert.ToInt32)
	manager.AddConv("int64", convert.ToInt64)
	manager.AddConv("float32", convert.ToFloat32)
	manager.AddConv("float64", convert.ToFloat64)
	manager.AddConv("bool", convert.ToBool)
	manager.AddConv("string", convert.ToString)
	manager.AddConv("[]byte", convert.ToBytes)
	manager.AddConv("second", convert.ToUnixSecond)
	manager.AddConv("millisecond", convert.ToUnixMilli)
}
