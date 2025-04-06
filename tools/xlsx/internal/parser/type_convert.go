package parser

import (
	"hego/framework/basic"
	"hego/tools/xlsx/domain"
	"hego/tools/xlsx/internal/base"
	"hego/tools/xlsx/internal/manager"
	"strings"

	"github.com/spf13/cast"
)

func convert(name string, strs ...string) interface{} {
	rets := []interface{}{}
	for _, str := range strs {
		switch name {
		case "int32":
			rets = append(rets, cast.ToInt32(str))
		case "int64":
			rets = append(rets, cast.ToInt64(str))
		case "uint32":
			rets = append(rets, cast.ToUint32(str))
		case "uint64":
			rets = append(rets, cast.ToUint64(str))
		case "float32":
			rets = append(rets, cast.ToFloat32(str))
		case "float64":
			rets = append(rets, cast.ToFloat64(str))
		case "bool":
			rets = append(rets, cast.ToBool(str))
		case "bytes":
			rets = append(rets, basic.StringToBytes(str))
		case "string":
			rets = append(rets, str)
		}
	}
	if len(strs) == 1 {
		return rets[0]
	}
	return rets
}

func convertEnum(en *base.Enum, strs ...string) interface{} {
	rets := []uint32{}
	for _, str := range strs {
		rets = append(rets, en.Values[str].Value)
	}
	if len(strs) == 1 {
		return rets[0]
	}
	return rets
}

func convertStruct(st *base.Struct, strs ...string) interface{} {
	ret := map[string]interface{}{}
	for i, field := range st.Converts[strs[0]] {
		switch field.Type.TypeOf {
		case domain.TypeOfBase:
			switch field.Type.ValueOf {
			case domain.ValueOfBase:
				ret[field.Name] = convert(field.Type.Name, strs[i])
			case domain.ValueOfList:
				ret[field.Name] = convert(field.Type.Name, strings.Split(strs[i], ",")...)
			case domain.ValueOfMap:
				// 暂时不支持
			case domain.ValueOfGroup:
				// 暂时不支持
			}
		case domain.TypeOfEnum:
			switch field.Type.ValueOf {
			case domain.ValueOfBase:
				ret[field.Name] = convertEnum(manager.GetEnum(field.Type.Name), strs[i])
			case domain.ValueOfList:
				ret[field.Name] = convertEnum(manager.GetEnum(field.Type.Name), strings.Split(strs[i], ",")...)
			case domain.ValueOfMap:
				// 暂时不支持
			case domain.ValueOfGroup:
				// 暂时不支持
			}
		}
	}
	return ret
}

func ConvertConfig(st *base.Config, vals ...string) map[string]interface{} {
	ret := map[string]interface{}{}
	for _, field := range st.List {
		switch field.Type.TypeOf {
		case domain.TypeOfBase:
			switch field.Type.ValueOf {
			case domain.ValueOfBase:
				ret[field.Name] = convert(field.Type.Name, vals[field.Position])
			case domain.ValueOfList:
				ret[field.Name] = convert(field.Type.Name, strings.Split(vals[field.Position], ",")...)
			case domain.ValueOfMap:
				// todo
			case domain.ValueOfGroup:
				// todo
			}
		case domain.TypeOfEnum:
			switch field.Type.ValueOf {
			case domain.ValueOfBase:
				ret[field.Name] = convertEnum(manager.GetEnum(field.Type.Name), vals[field.Position])
			case domain.ValueOfList:
				ret[field.Name] = convertEnum(manager.GetEnum(field.Type.Name), strings.Split(vals[field.Position], ",")...)
			case domain.ValueOfMap:
				// todo
			case domain.ValueOfGroup:
				// todo
			}
		case domain.TypeOfStruct:
			switch field.Type.ValueOf {
			case domain.ValueOfBase:
				ret[field.Name] = convertStruct(manager.GetStruct(field.Type.Name), strings.Split(vals[field.Position], ",")...)
			case domain.ValueOfList:
				tmps := []interface{}{}
				for _, str := range strings.Split(vals[field.Position], "|") {
					tmps = append(tmps, convertStruct(manager.GetStruct(field.Type.Name), strings.Split(str, ",")...))
				}
				ret[field.Name] = tmps
			case domain.ValueOfMap:
				// todo
			case domain.ValueOfGroup:
				// todo
			}
		}
	}
	return ret
}
