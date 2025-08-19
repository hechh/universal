package parser

import (
	"fmt"
	"path"
	"sort"
	"strings"
	"universal/library/uerror"
	"universal/library/util"
	"universal/tools/cfgtool/domain"
	"universal/tools/cfgtool/internal/manager"
	"universal/tools/cfgtool/internal/typespec"

	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
)

func ParseFiles(files ...string) error {
	for _, file := range files {
		fmt.Printf("解析文件: %s\n", path.Base(file))
		if err := parseFile(file); err != nil {
			return err
		}
	}
	for _, item := range manager.GetStructs() {
		parseStruct(item)
	}
	for _, item := range manager.GetConfigs() {
		parseConfig(item)
	}
	parseProto()
	return nil
}

func parseFile(filename string) error {
	fp, err := excelize.OpenFile(filename)
	if err != nil {
		return err
	}
	defer fp.Close()

	rows, err := fp.GetRows("生成表")
	if err != nil {
		if _, ok := err.(excelize.ErrSheetNotExist); ok {
			fmt.Printf("%s没有定义生成表\n", filename)
			return nil
		}
		return uerror.New(1, -1, "获取生成表失败:%s", err.Error())
	}
	file := strings.TrimSuffix(path.Base(filename), path.Ext(filename))
	defaultFile := path.Base(file)
	/*
	   @config[:filename]|sheet:结构名|map:字段名[,字段名]:别名|group:字段名[,字段名]:别名
	   @struct[:filename]|sheet:结构名
	   @enum[:filename]|sheet
	*/
	for _, items := range rows {
		for _, val := range items {
			strs := strings.Split(val, "|")
			if len(strs) <= 0 {
				continue
			}
			rules := strings.Split(strs[0], ":")
			switch strings.ToLower(util.Index[string](rules, 0, strs[0])) {
			case "@enum":
				datas, err := fp.GetRows(strs[1])
				if err != nil {
					return uerror.New(1, -1, "获取枚举数据失败:%s", err.Error())
				}
				parseEnum(datas, util.Index[string](rules, 1, defaultFile), strs[1])
			case "@struct":
				ones := strings.Split(strs[1], ":")
				datas, err := fp.GetRows(ones[0])
				if err != nil {
					return uerror.New(1, -1, "获取结构数据失败:%s", err.Error())
				}
				st := manager.GetOrNewStruct(ones[1])
				st.Set(util.Index[string](rules, 1, defaultFile), ones[0], datas)
			case "@config":
				ones := strings.Split(strs[1], ":")
				datas, err := fp.GetRows(ones[0])
				if err != nil {
					return uerror.New(1, -1, "获取结构数据失败:%s", err.Error())
				}
				cfg := manager.GetOrNewConfig(ones[1])
				cfg.Set(util.Index[string](rules, 1, defaultFile), ones[0], datas, util.Suffix[string](strs, 2))
			}
		}
	}
	return nil
}

func parseEnum(rows [][]string, file, sheet string) {
	for _, vals := range rows {
		for _, val := range vals {
			if !strings.HasPrefix(val, "E|") && !strings.HasPrefix(val, "e|") {
				continue
			}
			strs := strings.Split(val, "|")
			ee := manager.GetOrNewEnum(strs[2])
			ee.Set(file, sheet)
			ee.AddValue(strs[2]+strs[3], cast.ToInt32(strs[4]), strs[1])
		}
	}
}

func parseStruct(st *typespec.Struct) {
	for i, val := range st.Rows[1] {
		if len(val) <= 0 || len(st.Rows[0][i]) <= 0 {
			continue
		}
		isArr := strings.HasPrefix(val, "[]")
		vType := strings.TrimPrefix(val, "[]")
		st.AddField(&typespec.Field{
			Name: st.Rows[0][i],
			Type: &typespec.Type{
				Name:    manager.GetConvType(vType),
				TypeOf:  manager.GetTypeOf(vType),
				ValueOf: util.Or[int](isArr, domain.ValueOfList, domain.ValueOfBase),
			},
			Desc:     st.Rows[2][i],
			Position: i,
			ConvFunc: manager.GetConvFunc(vType),
		})
	}
	for _, vals := range st.Rows[3:] {
		for i, val := range vals {
			if len(val) <= 0 || val == "0" {
				continue
			}
			st.Converts[vals[0]] = append(st.Converts[vals[0]], st.FieldList[i])
		}
	}
}

func parseConfig(cfg *typespec.Config) {
	for i, val := range cfg.Rows[1] {
		if len(val) <= 0 || len(cfg.Rows[0][i]) <= 0 {
			continue
		}
		isArr := strings.HasPrefix(val, "[]")
		vType := strings.TrimPrefix(val, "[]")
		cfg.AddField(&typespec.Field{
			Name: cfg.Rows[0][i],
			Type: &typespec.Type{
				Name:    manager.GetConvType(vType),
				TypeOf:  manager.GetTypeOf(vType),
				ValueOf: util.Or[int](isArr, domain.ValueOfList, domain.ValueOfBase),
			},
			Desc:     cfg.Rows[2][i],
			Position: i,
			ConvFunc: manager.GetConvFunc(vType),
		})
	}

	// 默认索引
	cfg.AddIndex(&typespec.Index{
		Name: "List",
		Type: &typespec.Type{TypeOf: domain.TypeOfBase, ValueOf: domain.ValueOfList},
	})

	// 解析索引   map:字段名[,字段名]:别名|group:字段名[,字段名]:别名
	for _, val := range cfg.Rules {
		strs := strings.Split(val, ":")
		keys := []*typespec.Field{}
		for _, field := range strings.Split(strs[1], ",") {
			if cfg.Fields[field] == nil {
				panic(fmt.Sprintf("索引字段不存在:%s %s", val, field))
			}
			keys = append(keys, cfg.Fields[field])
		}
		cfg.AddIndex(&typespec.Index{
			Name: util.Index[string](strs, 2, strings.ReplaceAll(strs[1], ",", "")),
			Type: &typespec.Type{
				Name:    util.Or[string](len(keys) == 1, keys[0].Type.Name, "string"),
				TypeOf:  util.Or[int](len(keys) == 1, keys[0].Type.TypeOf, domain.TypeOfBase),
				ValueOf: util.Or[int](strings.ToLower(strs[0]) == "map", domain.ValueOfMap, domain.ValueOfGroup),
			},
			FieldList: keys,
		})
	}
}

func parseProto() {
	for _, item := range manager.GetEnums() {
		sort.Slice(item.ValueList, func(i, j int) bool {
			return item.ValueList[i].Value < item.ValueList[j].Value
		})
		pp := manager.GetOrNewProto(item.FileName)
		pp.AddEnum(item)
	}
	for _, item := range manager.GetStructs() {
		pp := manager.GetOrNewProto(item.FileName)
		pp.AddStruct(item)
		for _, field := range item.FieldList {
			switch field.Type.TypeOf {
			case domain.TypeOfEnum:
				enum := manager.GetEnum(field.Type.Name)
				pp.AddReference(enum.FileName)
			case domain.TypeOfStruct:
				st := manager.GetStruct(field.Type.Name)
				pp.AddReference(st.FileName)
			case domain.TypeOfConfig:
				cfg := manager.GetConfig(field.Type.Name)
				pp.AddReference(cfg.FileName)
			}
		}
	}
	for _, item := range manager.GetConfigs() {
		pp := manager.GetOrNewProto(item.FileName)
		pp.AddConfig(item)
		for _, field := range item.FieldList {
			switch field.Type.TypeOf {
			case domain.TypeOfEnum:
				enum := manager.GetEnum(field.Type.Name)
				pp.AddReference(enum.FileName)
			case domain.TypeOfStruct:
				st := manager.GetStruct(field.Type.Name)
				pp.AddReference(st.FileName)
			case domain.TypeOfConfig:
				cfg := manager.GetConfig(field.Type.Name)
				pp.AddReference(cfg.FileName)
			}
		}
	}
}
