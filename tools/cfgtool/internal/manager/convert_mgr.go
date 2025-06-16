package manager

import "github.com/spf13/cast"

var (
	convMgr = make(map[string]*ConvertInfo)
)

type ConvertInfo struct {
	Name     string
	ConvFunc func(string) interface{}
}

func GetConvFunc(name string) func(string) interface{} {
	if val, ok := convMgr[name]; ok {
		return val.ConvFunc
	}
	// 默认枚举转换函数
	if item, ok := enums[name]; ok {
		return func(str string) interface{} {
			if vv, ok := item.Values[str]; ok {
				return vv.Value
			}
			return cast.ToInt32(str)
		}
	}
	return nil
}

func GetConvType(name string) string {
	if val, ok := convMgr[name]; ok {
		return val.Name
	}
	return name
}

func init() {
	convMgr["int"] = &ConvertInfo{
		Name: "int32",
		ConvFunc: func(str string) interface{} {
			return cast.ToInt32(str)
		},
	}
	convMgr["int8"] = convMgr["int"]
	convMgr["int16"] = convMgr["int"]
	convMgr["int32"] = convMgr["int"]
	convMgr["int64"] = &ConvertInfo{
		Name: "int64",
		ConvFunc: func(str string) interface{} {
			return cast.ToInt64(str)
		},
	}

	convMgr["uint"] = &ConvertInfo{
		Name: "uint32",
		ConvFunc: func(str string) interface{} {
			return cast.ToUint32(str)
		},
	}
	convMgr["uint8"] = convMgr["uint"]
	convMgr["uint16"] = convMgr["uint"]
	convMgr["uint32"] = convMgr["uint"]
	convMgr["uint64"] = &ConvertInfo{
		Name: "uint64",
		ConvFunc: func(str string) interface{} {
			return cast.ToUint64(str)
		},
	}

	convMgr["float"] = &ConvertInfo{
		Name: "float64",
		ConvFunc: func(str string) interface{} {
			return cast.ToFloat64(str)
		},
	}
	convMgr["bool"] = &ConvertInfo{
		Name: "bool",
		ConvFunc: func(str string) interface{} {
			return cast.ToBool(str)
		},
	}
	convMgr["string"] = &ConvertInfo{
		Name: "string",
		ConvFunc: func(str string) interface{} {
			return str
		},
	}
}
