package repository

import (
	"universal/tools/gomaker/internal/manager"
	"universal/tools/gomaker/repository/convert"
	"universal/tools/gomaker/repository/xlsx"
)

func Init() {
	manager.Register("proto", "xlsx转pb结构和config代码", xlsx.EnumGen, xlsx.TableGen)
	manager.Register("config", "xlsx生成config代码", xlsx.ConfigGen)
	//manager.Register("bytes", "xlsx转bytes文件", xlsx.BytesGen)
}

func init() {
	manager.AddConv("uint32", "uint32", "uint32", convert.ToUint32)
	manager.AddConv("uint64", "uint64", "uint64", convert.ToUint64)
	manager.AddConv("int32", "int32", "int32", convert.ToInt32)
	manager.AddConv("int", "int", "int32", convert.ToInt32)
	manager.AddConv("int64", "int64", "int64", convert.ToInt64)
	manager.AddConv("float32", "float32", "float", convert.ToFloat32)
	manager.AddConv("float64", "float64", "double", convert.ToFloat64)
	manager.AddConv("bool", "bool", "bool", convert.ToBool)
	manager.AddConv("string", "string", "string", convert.ToString)
	manager.AddConv("bytes", "[]byte", "bytes", convert.ToBytes)
	manager.AddConv("second", "int64", "int64", convert.ToUnixSecond)
	manager.AddConv("millisecond", "int64", "int64", convert.ToUnixMilli)
}
