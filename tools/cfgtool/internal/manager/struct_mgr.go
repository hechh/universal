package manager

import (
	"hego/tools/cfgtool/internal/base"
	"sort"
	"strings"
)

var (
	structMgr = make(map[string]*base.Struct)
)

func GetOrNewStruct(file, sheet, name string) *base.Struct {
	if val, ok := structMgr[name]; ok {
		return val
	}
	structMgr[name] = &base.Struct{
		Name:     name,
		Sheet:    sheet,
		FileName: file,
		Fields:   make(map[string]*base.Field),
		Converts: make(map[string][]*base.Field),
	}
	return structMgr[name]
}

func GetStructList() (rets []*base.Struct) {
	for _, val := range structMgr {
		rets = append(rets, val)
	}
	sort.Slice(rets, func(i, j int) bool {
		return strings.Compare(rets[i].Name, rets[j].Name) <= 0
	})
	return
}

func GetStruct(name string) *base.Struct {
	if val, ok := structMgr[name]; ok {
		return val
	}
	return nil
}

func GetStructMap() map[string]*base.Struct {
	return structMgr
}
