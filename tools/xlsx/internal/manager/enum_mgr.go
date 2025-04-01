package manager

import (
	"sort"
	"universal/tools/xlsx/domain"
	"universal/tools/xlsx/internal/base"
)

var (
	enumMgr   = make(map[string]*base.Enum)
	structMgr = make(map[string]*base.Struct)
	configMgr = make(map[string]*base.Config)
	tableMgr  = make(map[string]*base.Table)
	tableList = []*base.Table{}
)

func GetFileInfo() (rets map[string][]interface{}) {
	rets = make(map[string][]interface{})
	for _, item := range enumMgr {
		rets[item.FileName] = append(rets[item.FileName], item)
	}
	for _, item := range structMgr {
		rets[item.FileName] = append(rets[item.FileName], item)
	}
	for _, item := range configMgr {
		rets[item.FileName] = append(rets[item.FileName], item)
	}
	return rets
}

func GetEnum(name string) *base.Enum {
	return enumMgr[name]
}

func GetStruct(name string) *base.Struct {
	return structMgr[name]
}

func GetConfig(name string) *base.Config {
	return configMgr[name]
}

func GetTable(name string) *base.Table {
	return tableMgr[name]
}

func GetTableList() (rets []*base.Table) {
	rets = append(rets, tableList...)
	sort.Slice(rets, func(i, j int) bool {
		return rets[i].TypeOf < rets[j].TypeOf
	})
	return
}

func GetTables(typeOf uint32) (rets []*base.Table) {
	for _, item := range tableMgr {
		if item.TypeOf == typeOf {
			rets = append(rets, item)
		}
	}
	return
}

func AddConfig(item *base.Config) {
	configMgr[item.Name] = item
}

func AddStruct(item *base.Struct) {
	structMgr[item.Name] = item
}

func AddTable(item *base.Table) {
	tableMgr[item.FileName] = item
	tableList = append(tableList, item)
}

func AddEnum(item *base.Value) {
	enum, ok := enumMgr[item.Type]
	if !ok {
		enum = &base.Enum{
			Name:     item.Type,
			FileName: item.FileName,
			Values:   make(map[string]*base.EValue),
		}
		enumMgr[item.Type] = enum
	}
	enum.Values[item.Desc] = &base.EValue{
		Name:  item.Name,
		Value: item.Value,
		Desc:  item.Desc,
	}
}

func GetTypeOf(name string) uint32 {
	if _, ok := enumMgr[name]; ok {
		return domain.TYPE_OF_ENUM
	}
	if _, ok := structMgr[name]; ok {
		return domain.TYPE_OF_STRUCT
	}
	if _, ok := configMgr[name]; ok {
		return domain.TYPE_OF_CONFIG
	}
	return domain.TYPE_OF_BASE
}
