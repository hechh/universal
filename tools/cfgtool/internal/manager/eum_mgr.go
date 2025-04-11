package manager

import (
	"hego/tools/cfgtool/internal/base"
	"sort"
	"strings"
)

var (
	enumMgr = make(map[string]*base.Enum)
)

func GetOrNewEnum(name string) *base.Enum {
	if val, ok := enumMgr[name]; ok {
		return val
	}
	enumMgr[name] = &base.Enum{
		Name:   name,
		Values: make(map[string]*base.EValue),
	}
	return enumMgr[name]
}

func GetEnum(name string) *base.Enum {
	return enumMgr[name]
}

func GetEnumMap() map[string]*base.Enum {
	return enumMgr
}

func GetEnumList() (rets []*base.Enum) {
	for _, val := range enumMgr {
		rets = append(rets, val)
	}
	sort.Slice(rets, func(i, j int) bool {
		return strings.Compare(rets[i].Name, rets[j].Name) <= 0
	})
	return
}
