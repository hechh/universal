package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"universal/framework/basic"
	"universal/framework/uerror"
	"universal/tools/xlsx/domain"
	"universal/tools/xlsx/internal/base"

	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
)

var (
	enumMgr   = make(map[string]*base.Enum)
	structMgr = make(map[string]*base.Struct)
	configMgr = make(map[string]*base.Config)
	tableMgr  = make(map[string]*base.Table)
)

func AddEnum(str string) {
	strs := strings.Split(str, ":")
	enumType := strs[2]
	fieldName := fmt.Sprintf("%s_%s", strs[2], strs[3])
	data, ok := enumMgr[enumType]
	if !ok {
		enumMgr[enumType] = &base.Enum{
			Name:   enumType,
			Values: make(map[string]*base.EnumValue),
		}
		data = enumMgr[enumType]
	}
	data.Values[fieldName] = &base.EnumValue{
		Name:  fieldName,
		Value: cast.ToUint32(strs[4]),
		Desc:  strs[1],
	}
}

func AddTable(fp *excelize.File, str string) error {
	strs := strings.Split(str, "|")
	index := strings.Index(strs[2], "@")
	if _, ok := tableMgr[strs[2][:index]]; ok {
		return uerror.NewUError(1, -1, "表%s已经存在", strs[2])
	}
	switch strs[0] {
	case "enum":
		tableMgr[strs[2][:index]] = &base.Table{
			TypeOf:   domain.TypeOfEnum,
			Sheet:    strs[2][:index],
			FileName: strs[2][index+1:],
			Fp:       fp,
		}
	case "struct":
		tableMgr[strs[2][:index]] = &base.Table{
			TypeOf:   domain.TypeOfStruct,
			Sheet:    strs[2][:index],
			FileName: strs[2][index+1:],
			Fp:       fp,
		}
	case "config":
		tableMgr[strs[2][:index]] = &base.Table{
			TypeOf:   domain.TypeOfConfig,
			Sheet:    strs[2][:index],
			FileName: strs[2][index+1:],
			Fp:       fp,
		}
	}
	return nil
}

func toCast(name string, str string) interface{} {
	switch name {
	case "int32":
		return cast.ToInt32(str)
	case "int64":
		return cast.ToInt64(str)
	case "uint32":
		return cast.ToUint32(str)
	case "uint64":
		return cast.ToUint64(str)
	case "float32":
		return cast.ToFloat32(str)
	case "float64":
		return cast.ToFloat64(str)
	case "bool":
		return cast.ToBool(str)
	case "bytes":
		return basic.StringToBytes(str)
	}
	return str
}

func baseConvert(name string, strs ...string) interface{} {
	if len(strs) == 1 {
		return toCast(name, strs[0])
	}
	rets := []interface{}{}
	for _, str := range strs {
		rets = append(rets, toCast(name, str))
	}
	return rets
}

func enumConvert(en *base.Enum, strs ...string) interface{} {
	if len(strs) == 1 {
		return en.Values[strs[0]].Value
	}
	rets := []uint32{}
	for _, str := range strs {
		rets = append(rets, en.Values[str].Value)
	}
	return rets
}

func structConvert(st *base.Struct, strs ...string) interface{} {
	ret := map[string]interface{}{}
	for i, field := range st.Converts[strs[0]] {
		switch field.Type.TypeOf {
		case domain.TypeOfBase:
			switch field.Type.ValueOf {
			case domain.ValueOfSingle:
				ret[field.Name] = baseConvert(field.Type.Name, strs[i])
			case domain.ValueOfArray:
				ret[field.Name] = baseConvert(field.Type.Name, strings.Split(strs[i], ",")...)
			}
		case domain.TypeOfEnum:
			switch field.Type.ValueOf {
			case domain.ValueOfSingle:
				ret[field.Name] = enumConvert(enumMgr[field.Type.Name], strs[i])
			case domain.ValueOfArray:
				ret[field.Name] = enumConvert(enumMgr[field.Type.Name], strings.Split(strs[i], ",")...)
			}
		}
	}
	return ret
}

func parseRow(st *base.Config, vals ...string) map[string]interface{} {
	ret := map[string]interface{}{}
	for _, field := range st.List {
		switch field.Type.TypeOf {
		case domain.TypeOfBase:
			switch field.Type.ValueOf {
			case domain.ValueOfSingle:
				ret[field.Name] = baseConvert(field.Type.Name, vals[field.Position])
			case domain.ValueOfArray:
				ret[field.Name] = baseConvert(field.Type.Name, strings.Split(vals[field.Position], ",")...)
			case domain.ValueOfMap:
				// todo
			case domain.ValueOfGroup:
				// todo
			}
		case domain.TypeOfEnum:
			switch field.Type.ValueOf {
			case domain.ValueOfSingle:
				ret[field.Name] = enumConvert(enumMgr[field.Type.Name], vals[field.Position])
			case domain.ValueOfArray:
				ret[field.Name] = enumConvert(enumMgr[field.Type.Name], strings.Split(vals[field.Position], ",")...)
			case domain.ValueOfMap:
				// todo
			case domain.ValueOfGroup:
				// todo
			}
		case domain.TypeOfStruct:
			switch field.Type.ValueOf {
			case domain.ValueOfSingle:
				ret[field.Name] = structConvert(structMgr[field.Type.Name], strings.Split(vals[field.Position], ",")...)
			case domain.ValueOfArray:
				tmps := []interface{}{}
				for _, str := range strings.Split(vals[field.Position], "|") {
					tmps = append(tmps, structConvert(structMgr[field.Type.Name], strings.Split(str, ",")...))
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

func ParseRows(name string, rows [][]string, buf *bytes.Buffer) error {
	st := configMgr[name]
	ret := []map[string]interface{}{}
	for _, row := range rows {
		ret = append(ret, parseRow(st, row...))
	}
	jsData, err := json.Marshal(ret)
	if err != nil {
		return err
	}
	buf.Reset()
	buf.Write(jsData)
	return nil
}
